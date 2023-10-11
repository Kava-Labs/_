package e2e_test

import (
	"context"
	"fmt"
	"strconv"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	abci "github.com/tendermint/tendermint/abci/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/kava-labs/kava/tests/util"
	communitytypes "github.com/kava-labs/kava/x/community/types"
	kavadisttypes "github.com/kava-labs/kava/x/kavadist/types"
)

func (suite *IntegrationTestSuite) TestUpgradeInflation_Disable() {
	afterUpgradeCtx := util.CtxAtHeight(suite.UpgradeHeight)

	// Get x/community for switchover time
	params, err := suite.Kava.Community.Params(afterUpgradeCtx, &communitytypes.QueryParamsRequest{})
	suite.Require().NoError(err)

	// Sleep until switchover time + 6 seconds for extra block
	sleepDuration := time.Until(params.Params.UpgradeTimeDisableInflation.Add(6 * time.Second))
	time.Sleep(sleepDuration)

	suite.Require().Eventually(func() bool {
		communityParams, err := suite.Kava.Community.Params(afterUpgradeCtx, &communitytypes.QueryParamsRequest{})
		suite.Require().NoError(err)

		// After params are set in x/community -- non-zero switchover time
		return !communityParams.Params.UpgradeTimeDisableInflation.Equal(time.Time{})
	}, 20*time.Second, 3*time.Second)

	// Fetch exact block when inflation stop event emitted
	// This is run after the switchover, so we don't need to poll
	_, switchoverHeight, err := suite.Kava.GetBeginBlockEventsFromQuery(
		context.Background(),
		fmt.Sprintf(
			"%s.%s EXISTS",
			communitytypes.EventTypeInflationStop,
			communitytypes.AttributeKeyDisableTime,
		),
	)
	suite.Require().NoError(err)
	suite.Require().NotZero(switchoverHeight)

	// 1 block before switchover
	beforeSwitchoverCtx := util.CtxAtHeight(switchoverHeight - 1)
	afterSwitchoverCtx := util.CtxAtHeight(switchoverHeight)

	suite.Run("x/mint inflation before switchover", func() {
		mintParams, err := suite.Kava.Mint.Params(
			beforeSwitchoverCtx,
			&minttypes.QueryParamsRequest{},
		)
		suite.NoError(err)
		kavaDistParams, err := suite.Kava.Kavadist.Params(
			beforeSwitchoverCtx,
			&kavadisttypes.QueryParamsRequest{},
		)
		suite.NoError(err)

		// Use .String() to compare Decs since x/mint uses the deprecated one,
		// mismatch of types but same value.
		suite.Equal(
			sdkmath.LegacyMustNewDecFromStr("0.595000000000000000").String(),
			mintParams.Params.InflationMin.String(),
			"x/mint inflation min should be 59.5%% before switchover",
		)
		suite.Equal(
			sdkmath.LegacyMustNewDecFromStr("0.595000000000000000").String(),
			mintParams.Params.InflationMax.String(),
			"x/mint inflation max should be 59.5%% before switchover",
		)

		suite.True(
			kavaDistParams.Params.Active,
			"x/kavadist should be active before switchover",
		)
	})

	suite.Run("x/mint inflation after switchover", func() {
		mintParams, err := suite.Kava.Mint.Params(
			afterSwitchoverCtx,
			&minttypes.QueryParamsRequest{},
		)
		suite.NoError(err)
		kavaDistParams, err := suite.Kava.Kavadist.Params(
			afterSwitchoverCtx,
			&kavadisttypes.QueryParamsRequest{},
		)
		suite.NoError(err)

		suite.Equal(
			sdkmath.LegacyZeroDec().String(),
			mintParams.Params.InflationMin.String(),
			"x/mint inflation min should be 0% after switchover",
		)
		suite.Equal(
			sdkmath.LegacyZeroDec().String(),
			mintParams.Params.InflationMax.String(),
			"x/mint inflation max should be 0% after switchover",
		)

		suite.False(
			kavaDistParams.Params.Active,
			"x/kavadist should be inactive after switchover",
		)
	})

	// Ensure inflation was still active before switchover
	suite.Run("positive mint events before switchover", func() {
		// 1 block before switchover
		queryHeight := switchoverHeight - 1

		block, err := suite.Kava.TmSignClient.BlockResults(
			context.Background(),
			&queryHeight,
		)
		suite.Require().NoError(err)

		// Mint events should only occur in begin block
		mintEvents := FilterEventsByType(block.BeginBlockEvents, minttypes.EventTypeMint)

		suite.Require().NotEmpty(mintEvents, "mint events should be emitted")

		// Ensure mint amounts are non-zero
		found := false
		for _, event := range mintEvents {
			for _, attribute := range event.Attributes {
				// Bonded ratio and annual provisions unchecked

				if string(attribute.Key) == minttypes.AttributeKeyInflation {
					suite.Equal(
						sdkmath.LegacyMustNewDecFromStr("0.595000000000000000").String(),
						string(attribute.Value),
						"inflation should be 59.5%% before switchover",
					)
				}

				if string(attribute.Key) == sdk.AttributeKeyAmount {
					found = true
					// Parse as native go int, not necessary to use sdk.Int
					value, err := strconv.Atoi(string(attribute.Value))
					suite.Require().NoError(err)

					suite.NotZero(value, "mint amount should be non-zero")
					suite.Positive(value, "mint amount should be positive")
				}
			}
		}

		suite.True(found, "mint amount should be found")
	})

	suite.Run("staking denom supply increases before switchover", func() {
		queryHeight := switchoverHeight - 2

		supply1, err := suite.Kava.Bank.SupplyOf(
			util.CtxAtHeight(queryHeight),
			&types.QuerySupplyOfRequest{
				Denom: suite.Kava.StakingDenom,
			},
		)
		suite.Require().NoError(err)

		suite.NotZero(supply1.Amount, "ukava supply should be non-zero")

		// Next block
		queryHeight += 1
		supply2, err := suite.Kava.Bank.SupplyOf(
			util.CtxAtHeight(queryHeight),
			&types.QuerySupplyOfRequest{
				Denom: suite.Kava.StakingDenom,
			},
		)
		suite.Require().NoError(err)

		suite.NotZero(supply2.Amount, "ukava supply should be non-zero")

		suite.Truef(
			supply2.Amount.Amount.GT(supply1.Amount.Amount),
			"ukava supply before switchover should increase between blocks, %s > %s",
			supply2.Amount.Amount.String(),
		)
	})

	// Check if inflation is ACTUALLY disabled... check if any coins are being
	// minted in the blocks after switchover
	suite.Run("no minting after switchover", func() {
		kavaSupply := sdk.NewCoin(suite.Kava.StakingDenom, sdkmath.ZeroInt())

		// Next 5 blocks after switchover, ensure there's actually no more inflation
		for i := 0; i < 5; i++ {
			queryHeight := switchoverHeight + int64(i)

			suite.Run(
				fmt.Sprintf("x/mint events should with 0 amount @ height=%d", queryHeight),
				func() {
					var block *coretypes.ResultBlockResults
					suite.Require().Eventually(func() bool {
						// Check begin block events
						block, err = suite.Kava.TmSignClient.BlockResults(
							context.Background(),
							&queryHeight,
						)

						return err == nil
					}, 20*time.Second, 3*time.Second)

					var mintEvents []abci.Event

					// Mint events should only occur in begin block, but we just include
					// everything else just in case anything changes in x/mint
					mintEventsBegin := FilterEventsByType(block.BeginBlockEvents, minttypes.EventTypeMint)
					mintEventsEnd := FilterEventsByType(block.EndBlockEvents, minttypes.EventTypeMint)
					mintEventsTx := FilterTxEventsByType(block.TxsResults, minttypes.EventTypeMint)

					mintEvents = append(mintEvents, mintEventsBegin...)
					mintEvents = append(mintEvents, mintEventsEnd...)
					mintEvents = append(mintEvents, mintEventsTx...)

					suite.Require().NotEmpty(mintEvents, "mint events should still be emitted")

					// Ensure mint amounts are 0
					found := false
					for _, event := range mintEvents {
						for _, attribute := range event.Attributes {
							// Bonded ratio and annual provisions unchecked

							if string(attribute.Key) == minttypes.AttributeKeyInflation {
								suite.Equal(sdkmath.LegacyZeroDec().String(), string(attribute.Value))
							}

							if string(attribute.Key) == sdk.AttributeKeyAmount {
								found = true
								suite.Equal(sdkmath.ZeroInt().String(), string(attribute.Value))
							}
						}
					}

					suite.True(found, "mint amount should be found")
				},
			)

			// Run this after the events check, since that one waits for the
			// new block if necessary
			suite.Run(
				fmt.Sprintf("total staking denom supply should not change @ height=%d", queryHeight),
				func() {
					supplyRes, err := suite.Kava.Bank.SupplyOf(
						util.CtxAtHeight(queryHeight),
						&types.QuerySupplyOfRequest{
							Denom: suite.Kava.StakingDenom,
						},
					)
					suite.Require().NoError(err)

					if kavaSupply.IsZero() {
						// First iteration, set supply
						kavaSupply = supplyRes.Amount
					} else {
						suite.Require().Equal(
							kavaSupply,
							supplyRes.Amount,
							"ukava supply should not change",
						)
					}
				},
			)
		}
	})

	suite.Run("no staking rewards from x/community", func() {
		// 1 block before switchover
		queryHeight := switchoverHeight - 1

		block, err := suite.Kava.TmSignClient.BlockResults(
			context.Background(),
			&queryHeight,
		)
		suite.Require().NoError(err)

		// Events are not emitted if amount is 0
		stakingRewardEvents := FilterEventsByType(block.BeginBlockEvents, communitytypes.EventTypeStakingRewardsPaid)
		suite.Require().Empty(stakingRewardEvents, "staking reward events should not be emitted")
	})

	suite.Run("staking rewards pay out from x/community", func() {
		// after switchover
		queryHeight := switchoverHeight

		block, err := suite.Kava.TmSignClient.BlockResults(
			context.Background(),
			&queryHeight,
		)
		suite.Require().NoError(err)

		stakingRewardEvents := FilterEventsByType(block.BeginBlockEvents, communitytypes.EventTypeStakingRewardsPaid)
		suite.Require().NotEmpty(stakingRewardEvents, "staking reward events should be emitted")

		// Ensure amounts are non-zero
		found := false
		for _, attr := range stakingRewardEvents[0].Attributes {
			if string(attr.Key) == communitytypes.AttributeKeyStakingRewardAmount {
				coins, err := sdk.ParseCoinNormalized(string(attr.Value))
				suite.Require().NoError(err, "staking reward amount should be parsable coins")

				suite.True(coins.Amount.IsPositive(), "staking reward amount should be a positive amount")

				found = true
			}
		}

		suite.True(found, "staking reward amount should be found in events")
	})
}

// FilterEventsByType returns a slice of events that match the given type.
func FilterEventsByType(events []abci.Event, eventType string) []abci.Event {
	filteredEvents := []abci.Event{}

	for _, event := range events {
		if event.Type == eventType {
			filteredEvents = append(filteredEvents, event)
		}
	}

	return filteredEvents
}

// FilterTxEventsByType returns a slice of events that match the given type
// from a slice of ResponseDeliverTx.
func FilterTxEventsByType(txs []*abci.ResponseDeliverTx, eventType string) []abci.Event {
	filteredEvents := []abci.Event{}

	for _, tx := range txs {
		events := FilterEventsByType(tx.Events, eventType)
		filteredEvents = append(filteredEvents, events...)
	}

	return filteredEvents
}
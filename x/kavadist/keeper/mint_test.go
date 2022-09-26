package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/kava-labs/kava/x/kavadist/types"
)

func (suite *keeperTestSuite) TestMintExpiredPeriod() {
	initialSupply := suite.BankKeeper.GetSupply(suite.Ctx, types.GovDenom)
	suite.Require().NotPanics(func() { suite.Keeper.SetPreviousBlockTime(suite.Ctx, time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)) })
	ctx := suite.Ctx.WithBlockTime(time.Date(2022, 1, 1, 0, 7, 0, 0, time.UTC))
	err := suite.Keeper.MintPeriodInflation(ctx)
	suite.Require().NoError(err)
	finalSupply := suite.BankKeeper.GetSupply(suite.Ctx, types.GovDenom)
	suite.Require().Equal(initialSupply, finalSupply)
}

func (suite *keeperTestSuite) TestMintPeriodNotStarted() {
	initialSupply := suite.BankKeeper.GetSupply(suite.Ctx, types.GovDenom)
	suite.Require().NotPanics(func() { suite.Keeper.SetPreviousBlockTime(suite.Ctx, time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)) })
	ctx := suite.Ctx.WithBlockTime(time.Date(2019, 1, 1, 0, 7, 0, 0, time.UTC))
	err := suite.Keeper.MintPeriodInflation(ctx)
	suite.Require().NoError(err)
	finalSupply := suite.BankKeeper.GetSupply(suite.Ctx, types.GovDenom)
	suite.Require().Equal(initialSupply, finalSupply)
}

func (suite *keeperTestSuite) TestMintOngoingPeriod() {
	initialSupply := suite.BankKeeper.GetSupply(suite.Ctx, types.GovDenom)
	suite.Require().NotPanics(func() {
		suite.Keeper.SetPreviousBlockTime(suite.Ctx, time.Date(2020, time.March, 1, 1, 0, 1, 0, time.UTC))
	})
	ctx := suite.Ctx.WithBlockTime(time.Date(2021, 2, 28, 23, 59, 59, 0, time.UTC))
	err := suite.Keeper.MintPeriodInflation(ctx)
	suite.Require().NoError(err)
	finalSupply := suite.BankKeeper.GetSupply(suite.Ctx, types.GovDenom)
	suite.Require().True(finalSupply.Amount.GT(initialSupply.Amount))
	mAcc := suite.AccountKeeper.GetModuleAccount(ctx, types.ModuleName)
	mAccSupply := suite.BankKeeper.GetAllBalances(ctx, mAcc.GetAddress()).AmountOf(types.GovDenom)
	suite.Require().True(mAccSupply.Equal(finalSupply.Amount.Sub(initialSupply.Amount)))
	// expect that inflation is ~10%
	expectedSupply := sdk.NewDecFromInt(initialSupply.Amount).Mul(sdk.MustNewDecFromStr("1.1"))
	supplyError := sdk.OneDec().Sub((sdk.NewDecFromInt(finalSupply.Amount).Quo(expectedSupply))).Abs()
	suite.Require().True(supplyError.LTE(sdk.MustNewDecFromStr("0.001")))
}

func (suite *keeperTestSuite) TestMintPeriodTransition() {
	initialSupply := suite.BankKeeper.GetSupply(suite.Ctx, types.GovDenom)
	params := suite.Keeper.GetParams(suite.Ctx)
	periods := []types.Period{
		suite.TestPeriods[0],
		{
			Start:     time.Date(2021, time.March, 1, 1, 0, 0, 0, time.UTC),
			End:       time.Date(2022, time.March, 1, 1, 0, 0, 0, time.UTC),
			Inflation: sdk.MustNewDecFromStr("1.000000003022265980"),
		},
	}
	params.Periods = periods
	suite.Require().NotPanics(func() {
		suite.Keeper.SetParams(suite.Ctx, params)
	})
	suite.Require().NotPanics(func() {
		suite.Keeper.SetPreviousBlockTime(suite.Ctx, time.Date(2020, time.March, 1, 1, 0, 1, 0, time.UTC))
	})
	ctx := suite.Ctx.WithBlockTime(time.Date(2021, 3, 10, 0, 0, 0, 0, time.UTC))
	err := suite.Keeper.MintPeriodInflation(ctx)
	suite.Require().NoError(err)
	finalSupply := suite.BankKeeper.GetSupply(suite.Ctx, types.GovDenom)
	suite.Require().True(finalSupply.Amount.GT(initialSupply.Amount))
}

func (suite *keeperTestSuite) TestMintNotActive() {
	initialSupply := suite.BankKeeper.GetSupply(suite.Ctx, types.GovDenom)
	params := suite.Keeper.GetParams(suite.Ctx)
	params.Active = false
	suite.Require().NotPanics(func() {
		suite.Keeper.SetParams(suite.Ctx, params)
	})
	suite.Require().NotPanics(func() {
		suite.Keeper.SetPreviousBlockTime(suite.Ctx, time.Date(2020, time.March, 1, 1, 0, 1, 0, time.UTC))
	})
	ctx := suite.Ctx.WithBlockTime(time.Date(2021, 2, 28, 23, 59, 59, 0, time.UTC))
	err := suite.Keeper.MintPeriodInflation(ctx)
	suite.Require().NoError(err)
	finalSupply := suite.BankKeeper.GetSupply(suite.Ctx, types.GovDenom)
	suite.Require().Equal(initialSupply, finalSupply)
}

func (suite *keeperTestSuite) TestInfraMinting() {

	type args struct {
		startTime           time.Time
		endTime             time.Time
		infraPeriods        types.Periods
		expectedFinalSupply sdk.Coin
		marginOfError       sdk.Dec
	}

	type errArgs struct {
		expectPass bool
		contains   string
	}

	type test struct {
		name    string
		args    args
		errArgs errArgs
	}

	testCases := []test{
		{
			"5% apy one year",
			args{
				startTime:           time.Date(2022, time.October, 1, 1, 0, 0, 0, time.UTC),
				endTime:             time.Date(2023, time.October, 1, 1, 0, 0, 0, time.UTC),
				infraPeriods:        types.Periods{types.NewPeriod(time.Date(2022, time.October, 1, 1, 0, 0, 0, time.UTC), time.Date(2023, time.October, 1, 1, 0, 0, 0, time.UTC), sdk.MustNewDecFromStr("1.000000001547125958"))},
				expectedFinalSupply: sdk.NewCoin(types.GovDenom, sdk.NewInt(1050000000000)),
				marginOfError:       sdk.MustNewDecFromStr("0.0001"),
			},
			errArgs{
				expectPass: true,
				contains:   "",
			},
		},
		{
			"5% apy 10 seconds",
			args{
				startTime:           time.Date(2022, time.October, 1, 1, 0, 0, 0, time.UTC),
				endTime:             time.Date(2022, time.October, 1, 1, 0, 10, 0, time.UTC),
				infraPeriods:        types.Periods{types.NewPeriod(time.Date(2022, time.October, 1, 1, 0, 0, 0, time.UTC), time.Date(2023, time.October, 1, 1, 0, 0, 0, time.UTC), sdk.MustNewDecFromStr("1.000000001547125958"))},
				expectedFinalSupply: sdk.NewCoin(types.GovDenom, sdk.NewInt(1000000015471)),
				marginOfError:       sdk.MustNewDecFromStr("0.0001"),
			},
			errArgs{
				expectPass: true,
				contains:   "",
			},
		},
	}

	for _, tc := range testCases {
		suite.SetupTest()
		params := types.NewParams(true, types.DefaultPeriods, types.NewInfraParams(tc.args.infraPeriods, types.DefaultInfraParams.PartnerRewards, types.DefaultInfraParams.CoreRewards))
		ctx := suite.Ctx.WithBlockTime(tc.args.startTime)
		suite.Keeper.SetParams(ctx, params)
		suite.Require().NotPanics(func() {
			suite.Keeper.SetPreviousBlockTime(ctx, tc.args.startTime)
		})
		ctx = suite.Ctx.WithBlockTime(tc.args.endTime)
		err := suite.Keeper.MintPeriodInflation(ctx)
		suite.Require().NoError(err)
		finalSupply := suite.BankKeeper.GetSupply(ctx, types.GovDenom)
		marginHigh := tc.args.expectedFinalSupply.Amount.ToDec().Mul(sdk.OneDec().Add(tc.args.marginOfError))
		marginLow := tc.args.expectedFinalSupply.Amount.ToDec().Mul(sdk.OneDec().Sub(tc.args.marginOfError))
		suite.Require().True(finalSupply.Amount.ToDec().LTE(marginHigh))
		suite.Require().True(finalSupply.Amount.ToDec().GTE(marginLow))

	}

}

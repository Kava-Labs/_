package keeper_test

import (
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/harvest/types"
	"github.com/kava-labs/kava/x/pricefeed"
)

func (suite *KeeperTestSuite) TestRepay() {
	type args struct {
		borrower                  sdk.AccAddress
		initialBorrowerCoins      sdk.Coins
		initialModuleCoins        sdk.Coins
		depositCoins              []sdk.Coin
		borrowCoins               sdk.Coins
		repayCoins                sdk.Coins
		expectedAccountBalance    sdk.Coins
		expectedModAccountBalance sdk.Coins
	}

	type errArgs struct {
		expectPass bool
		contains   string
	}

	type borrowTest struct {
		name    string
		args    args
		errArgs errArgs
	}

	testCases := []borrowTest{
		{
			"valid: partial repay",
			args{
				borrower:             sdk.AccAddress(crypto.AddressHash([]byte("test"))),
				initialBorrowerCoins: sdk.NewCoins(sdk.NewCoin("ukava", sdk.NewInt(100*KAVA_CF))),
				initialModuleCoins:   sdk.NewCoins(sdk.NewCoin("ukava", sdk.NewInt(1000*KAVA_CF)), sdk.NewCoin("usdx", sdk.NewInt(1000*USDX_CF))),
				depositCoins:         []sdk.Coin{sdk.NewCoin("ukava", sdk.NewInt(100*KAVA_CF))},
				borrowCoins:          sdk.NewCoins(sdk.NewCoin("ukava", sdk.NewInt(50*KAVA_CF))),
				repayCoins:           sdk.NewCoins(sdk.NewCoin("ukava", sdk.NewInt(10*KAVA_CF))),
			},
			errArgs{
				expectPass: true,
				contains:   "",
			},
		},
		{
			"valid: repay in full",
			args{
				borrower:             sdk.AccAddress(crypto.AddressHash([]byte("test"))),
				initialBorrowerCoins: sdk.NewCoins(sdk.NewCoin("ukava", sdk.NewInt(100*KAVA_CF))),
				initialModuleCoins:   sdk.NewCoins(sdk.NewCoin("ukava", sdk.NewInt(1000*KAVA_CF)), sdk.NewCoin("usdx", sdk.NewInt(1000*USDX_CF))),
				depositCoins:         []sdk.Coin{sdk.NewCoin("ukava", sdk.NewInt(100*KAVA_CF))},
				borrowCoins:          sdk.NewCoins(sdk.NewCoin("ukava", sdk.NewInt(50*KAVA_CF))),
				repayCoins:           sdk.NewCoins(sdk.NewCoin("ukava", sdk.NewInt(50*KAVA_CF))),
			},
			errArgs{
				expectPass: true,
				contains:   "",
			},
		},
		{
			"invalid: insufficent balance for repay",
			args{
				borrower:             sdk.AccAddress(crypto.AddressHash([]byte("test"))),
				initialBorrowerCoins: sdk.NewCoins(sdk.NewCoin("ukava", sdk.NewInt(100*KAVA_CF))),
				initialModuleCoins:   sdk.NewCoins(sdk.NewCoin("ukava", sdk.NewInt(1000*KAVA_CF)), sdk.NewCoin("usdx", sdk.NewInt(1000*USDX_CF))),
				depositCoins:         []sdk.Coin{sdk.NewCoin("ukava", sdk.NewInt(100*KAVA_CF))},
				borrowCoins:          sdk.NewCoins(sdk.NewCoin("ukava", sdk.NewInt(50*KAVA_CF))),
				repayCoins:           sdk.NewCoins(sdk.NewCoin("ukava", sdk.NewInt(51*KAVA_CF))), // Exceeds user's KAVA balance
			},
			errArgs{
				expectPass: false,
				contains:   "account can only repay up to",
			},
		},
		{
			"invalid: overpaid debt",
			args{
				borrower:             sdk.AccAddress(crypto.AddressHash([]byte("test"))),
				initialBorrowerCoins: sdk.NewCoins(sdk.NewCoin("ukava", sdk.NewInt(100*KAVA_CF))),
				initialModuleCoins:   sdk.NewCoins(sdk.NewCoin("ukava", sdk.NewInt(1000*KAVA_CF)), sdk.NewCoin("usdx", sdk.NewInt(1000*USDX_CF))),
				depositCoins:         []sdk.Coin{sdk.NewCoin("ukava", sdk.NewInt(80*KAVA_CF))}, // Deposit less so user still has some KAVA
				borrowCoins:          sdk.NewCoins(sdk.NewCoin("ukava", sdk.NewInt(50*KAVA_CF))),
				repayCoins:           sdk.NewCoins(sdk.NewCoin("ukava", sdk.NewInt(60*KAVA_CF))), // Exceeds borrowed coins but not user's balance
			},
			errArgs{
				expectPass: false,
				contains:   "repayment exceeds loan debt",
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Initialize test app and set context
			tApp := app.NewTestApp()
			ctx := tApp.NewContext(true, abci.Header{Height: 1, Time: tmtime.Now()})

			// Auth module genesis state
			authGS := app.NewAuthGenState(
				[]sdk.AccAddress{tc.args.borrower},
				[]sdk.Coins{tc.args.initialBorrowerCoins})

			// Harvest module genesis state
			harvestGS := types.NewGenesisState(types.NewParams(
				true,
				types.DistributionSchedules{
					types.NewDistributionSchedule(true, "usdx", time.Date(2020, 10, 8, 14, 0, 0, 0, time.UTC), time.Date(2020, 11, 22, 14, 0, 0, 0, time.UTC), sdk.NewCoin("hard", sdk.NewInt(5000)), time.Date(2021, 11, 22, 14, 0, 0, 0, time.UTC), types.Multipliers{types.NewMultiplier(types.Small, 0, sdk.MustNewDecFromStr("0.33")), types.NewMultiplier(types.Medium, 6, sdk.MustNewDecFromStr("0.5")), types.NewMultiplier(types.Medium, 24, sdk.OneDec())}),
					types.NewDistributionSchedule(true, "ukava", time.Date(2020, 10, 8, 14, 0, 0, 0, time.UTC), time.Date(2020, 11, 22, 14, 0, 0, 0, time.UTC), sdk.NewCoin("hard", sdk.NewInt(5000)), time.Date(2021, 11, 22, 14, 0, 0, 0, time.UTC), types.Multipliers{types.NewMultiplier(types.Small, 0, sdk.MustNewDecFromStr("0.33")), types.NewMultiplier(types.Medium, 6, sdk.MustNewDecFromStr("0.5")), types.NewMultiplier(types.Medium, 24, sdk.OneDec())}),
				},
				types.DelegatorDistributionSchedules{types.NewDelegatorDistributionSchedule(
					types.NewDistributionSchedule(true, "usdx", time.Date(2020, 10, 8, 14, 0, 0, 0, time.UTC), time.Date(2025, 10, 8, 14, 0, 0, 0, time.UTC), sdk.NewCoin("hard", sdk.NewInt(500)), time.Date(2026, 10, 8, 14, 0, 0, 0, time.UTC), types.Multipliers{types.NewMultiplier(types.Small, 0, sdk.MustNewDecFromStr("0.33")), types.NewMultiplier(types.Medium, 6, sdk.MustNewDecFromStr("0.5")), types.NewMultiplier(types.Medium, 24, sdk.OneDec())}),
					time.Hour*24,
				),
				},
				types.MoneyMarkets{
					types.NewMoneyMarket("usdx", true, sdk.NewDec(100000000*KAVA_CF), sdk.MustNewDecFromStr("1"), "usdx:usd", sdk.NewInt(USDX_CF), types.NewInterestRateModel(sdk.MustNewDecFromStr("0.05"), sdk.MustNewDecFromStr("2"), sdk.MustNewDecFromStr("0.8"), sdk.MustNewDecFromStr("10"))),
					types.NewMoneyMarket("ukava", false, sdk.NewDec(100000000*KAVA_CF), sdk.MustNewDecFromStr("0.8"), "kava:usd", sdk.NewInt(KAVA_CF), types.NewInterestRateModel(sdk.MustNewDecFromStr("0.05"), sdk.MustNewDecFromStr("2"), sdk.MustNewDecFromStr("0.8"), sdk.MustNewDecFromStr("10"))),
				},
			), types.DefaultPreviousBlockTime, types.DefaultDistributionTimes)

			// Pricefeed module genesis state
			pricefeedGS := pricefeed.GenesisState{
				Params: pricefeed.Params{
					Markets: []pricefeed.Market{
						{MarketID: "usdx:usd", BaseAsset: "usdx", QuoteAsset: "usd", Oracles: []sdk.AccAddress{}, Active: true},
						{MarketID: "kava:usd", BaseAsset: "kava", QuoteAsset: "usd", Oracles: []sdk.AccAddress{}, Active: true},
					},
				},
				PostedPrices: []pricefeed.PostedPrice{
					{
						MarketID:      "usdx:usd",
						OracleAddress: sdk.AccAddress{},
						Price:         sdk.MustNewDecFromStr("1.00"),
						Expiry:        time.Now().Add(1 * time.Hour),
					},
					{
						MarketID:      "kava:usd",
						OracleAddress: sdk.AccAddress{},
						Price:         sdk.MustNewDecFromStr("2.00"),
						Expiry:        time.Now().Add(1 * time.Hour),
					},
				},
			}

			// Initialize test application
			tApp.InitializeFromGenesisStates(authGS,
				app.GenesisState{pricefeed.ModuleName: pricefeed.ModuleCdc.MustMarshalJSON(pricefeedGS)},
				app.GenesisState{types.ModuleName: types.ModuleCdc.MustMarshalJSON(harvestGS)},
			)

			// Mint coins to Harvest module account
			supplyKeeper := tApp.GetSupplyKeeper()
			supplyKeeper.MintCoins(ctx, types.ModuleAccountName, tc.args.initialModuleCoins)

			keeper := tApp.GetHarvestKeeper()
			suite.app = tApp
			suite.ctx = ctx
			suite.keeper = keeper

			var err error

			// Deposit coins to harvest
			depositedCoins := sdk.NewCoins()
			for _, depositCoin := range tc.args.depositCoins {
				err = suite.keeper.Deposit(suite.ctx, tc.args.borrower, depositCoin)
				suite.Require().NoError(err)
				depositedCoins.Add(depositCoin)
			}

			// Borrow coins from harvest
			err = suite.keeper.Borrow(suite.ctx, tc.args.borrower, tc.args.borrowCoins)
			suite.Require().NoError(err)

			err = suite.keeper.Repay(suite.ctx, tc.args.borrower, tc.args.repayCoins)
			if tc.errArgs.expectPass {
				suite.Require().NoError(err)

				// Check borrower balance
				expectedBorrowerCoins := tc.args.initialBorrowerCoins.Sub(tc.args.depositCoins).Add(tc.args.borrowCoins...).Sub(tc.args.repayCoins)
				acc := suite.getAccount(tc.args.borrower)
				suite.Require().Equal(expectedBorrowerCoins, acc.GetCoins())

				// Check module account balance
				expectedModuleCoins := tc.args.initialModuleCoins.Add(tc.args.depositCoins...).Sub(tc.args.borrowCoins).Add(tc.args.repayCoins...)
				mAcc := suite.getModuleAccount(types.ModuleAccountName)
				suite.Require().Equal(expectedModuleCoins, mAcc.GetCoins())

				// Check user's borrow object
				borrow, _ := suite.keeper.GetBorrow(suite.ctx, tc.args.borrower)
				expectedBorrowCoins := tc.args.borrowCoins.Sub(tc.args.repayCoins)
				suite.Require().Equal(expectedBorrowCoins, borrow.Amount)
			} else {
				suite.Require().Error(err)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))

				// Check borrower balance (no repay coins)
				expectedBorrowerCoins := tc.args.initialBorrowerCoins.Sub(tc.args.depositCoins).Add(tc.args.borrowCoins...)
				acc := suite.getAccount(tc.args.borrower)
				suite.Require().Equal(expectedBorrowerCoins, acc.GetCoins())

				// Check module account balance (no repay coins)
				expectedModuleCoins := tc.args.initialModuleCoins.Add(tc.args.depositCoins...).Sub(tc.args.borrowCoins)
				mAcc := suite.getModuleAccount(types.ModuleAccountName)
				suite.Require().Equal(expectedModuleCoins, mAcc.GetCoins())

				// Check user's borrow object (no repay coins)
				borrow, _ := suite.keeper.GetBorrow(suite.ctx, tc.args.borrower)
				suite.Require().Equal(tc.args.borrowCoins, borrow.Amount)
			}
		})
	}
}

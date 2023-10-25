package app

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	communitytypes "github.com/kava-labs/kava/x/community/types"
)

const (
	UpgradeName_Mainnet = "v0.25.0"
	UpgradeName_Testnet = "v0.25.0-alpha.0"
	UpgradeName_E2ETest = "v0.25.0-testing"
)

var (
	// KAVA to ukava - 6 decimals
	kavaConversionFactor = sdk.NewInt(1000_000)
	secondsPerYear       = sdk.NewInt(365 * 24 * 60 * 60)

	// 10 Million KAVA per year in staking rewards, inflation disable time 2024-01-01T00:00:00 UTC
	CommunityParams_Mainnet = communitytypes.NewParams(
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		// before switchover
		sdkmath.LegacyZeroDec(),
		// after switchover - 10M KAVA to ukava per year / seconds per year
		sdkmath.LegacyNewDec(10_000_000).
			MulInt(kavaConversionFactor).
			QuoInt(secondsPerYear),
	)

	// Testnet -- 15 Trillion KAVA per year in staking rewards, inflation disable time 2023-11-16T00:00:00 UTC
	CommunityParams_Testnet = communitytypes.NewParams(
		time.Date(2023, 11, 16, 0, 0, 0, 0, time.UTC),
		// before switchover
		sdkmath.LegacyZeroDec(),
		// after switchover
		sdkmath.LegacyNewDec(15_000_000).
			MulInt64(1_000_000). // 15M * 1M = 15T
			MulInt(kavaConversionFactor).
			QuoInt(secondsPerYear),
	)

	CommunityParams_E2E = communitytypes.NewParams(
		time.Now().Add(10*time.Second).UTC(), // relative time for testing
		sdkmath.LegacyNewDec(0),              // stakingRewardsPerSecond
		sdkmath.LegacyNewDec(1000),           // upgradeTimeSetstakingRewardsPerSecond
	)
)

// RegisterUpgradeHandlers registers the upgrade handlers for the app.
func (app App) RegisterUpgradeHandlers() {
	app.upgradeKeeper.SetUpgradeHandler(
		UpgradeName_Mainnet,
		upgradeHandler(app, UpgradeName_Mainnet, CommunityParams_Mainnet),
	)
	app.upgradeKeeper.SetUpgradeHandler(
		UpgradeName_Testnet,
		upgradeHandler(app, UpgradeName_Testnet, CommunityParams_Testnet),
	)
	app.upgradeKeeper.SetUpgradeHandler(
		UpgradeName_E2ETest,
		upgradeHandler(app, UpgradeName_Testnet, CommunityParams_E2E),
	)

	upgradeInfo, err := app.upgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	doUpgrade := upgradeInfo.Name == UpgradeName_Mainnet ||
		upgradeInfo.Name == UpgradeName_Testnet ||
		upgradeInfo.Name == UpgradeName_E2ETest

	if doUpgrade && !app.upgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{
				// x/community added store
				communitytypes.ModuleName,
			},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

// upgradeHandler returns an UpgradeHandler for the given upgrade parameters.
func upgradeHandler(
	app App,
	name string,
	communityParams communitytypes.Params,
) upgradetypes.UpgradeHandler {
	return func(
		ctx sdk.Context,
		plan upgradetypes.Plan,
		fromVM module.VersionMap,
	) (module.VersionMap, error) {
		app.Logger().Info(fmt.Sprintf("running %s upgrade handler", name))

		toVM, err := app.mm.RunMigrations(ctx, app.configurator, fromVM)
		if err != nil {
			return toVM, err
		}

		app.communityKeeper.SetParams(ctx, communityParams)
		app.Logger().Info(
			"initialized x/community params",
			"UpgradeTimeDisableInflation", communityParams.UpgradeTimeDisableInflation,
			"StakingRewardsPerSecond", communityParams.StakingRewardsPerSecond,
			"UpgradeTimeSetStakingRewardsPerSecond", communityParams.UpgradeTimeSetStakingRewardsPerSecond,
		)

		return toVM, nil
	}
}

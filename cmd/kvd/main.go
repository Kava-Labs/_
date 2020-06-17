package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	dbm "github.com/tendermint/tm-db"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/kava-labs/kava/app"
	kava3 "github.com/kava-labs/kava/contrib/kava-3"
	"github.com/kava-labs/kava/migrate"
)

// kvd custom flags
const flagInvCheckPeriod = "inv-check-period"

var invCheckPeriod uint

func main() {
	appCodec, cdc := app.MakeCodec()

	config := sdk.GetConfig()
	app.SetBech32AddressPrefixes(config)
	app.SetBip44CoinType(config)
	config.Seal()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "kvd",
		Short:             "Kava Daemon (server)",
		PersistentPreRunE: persistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(
		genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(ctx, cdc, banktypes.GenesisBalancetIterator{}, app.DefaultNodeHome),
		migrate.MigrateGenesisCmd(ctx, cdc),
		writeParamsAndConfigCmd(cdc),
		genutilcli.GenTxCmd(
			ctx,
			cdc,
			app.ModuleBasics,
			staking.AppModuleBasic{},
			bank.GenesisBalancesIterator{},
			app.DefaultNodeHome,
			app.DefaultCLIHome),
		genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics),
		AddGenesisAccountCmd(ctx, cdc, appCodec, app.DefaultNodeHome, app.DefaultCLIHome),
		testnetCmd(ctx, cdc, app.ModuleBasics, banktypes.GenesisBalancetIterator{}),
		flags.NewCompletionCmd(rootCmd, true),
		debug.Cmd(cdc),
	)

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "KA", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	var cache sdk.MultiStorePersistentCache

	if viper.GetBool(server.FlagInterBlockCache) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range viper.GetIntSlice(server.FlagUnsafeSkipUpgrades) {
		skipUpgradeHeights[int64(h)] = true
	}

	return app.NewApp(
		logger, db, traceStore, true, skipUpgradeHeights,
		viper.GetString(flags.FlagHome), invCheckPeriod,
		baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetHaltHeight(viper.GetUint64(server.FlagHaltHeight)),
		baseapp.SetHaltTime(viper.GetUint64(server.FlagHaltTime)),
		baseapp.SetInterBlockCache(cache),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, *abci.ConsensusParams, error) {

	if height != -1 {
		tempApp := app.NewApp(logger, db, traceStore, false, map[int64]bool{}, "", uint(1))
		err := tempApp.LoadHeight(height)
		if err != nil {
			return nil, nil, nil, err
		}
		return tempApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}
	tempApp := app.NewApp(logger, db, traceStore, true, map[int64]bool{}, "", uint(1))
	return tempApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

// persistentPreRunEFn wraps the sdk function server.PersistentPreRunEFn to error on invaid pruning config.
func persistentPreRunEFn(ctx *server.Context) func(*cobra.Command, []string) error {

	originalFunc := server.PersistentPreRunEFn(ctx)

	return func(cmd *cobra.Command, args []string) error {

		if err := originalFunc(cmd, args); err != nil {
			return err
		}

		// check pruning config for `kvd start`
		if cmd.Name() == "start" {
			if viper.GetString("pruning") == store.PruningStrategySyncable {
				return fmt.Errorf(
					"invalid app config: pruning == '%s'. Update config (%s) with pruning set to '%s' or '%s'.",
					store.PruningStrategySyncable, viper.ConfigFileUsed(), store.PruningStrategyNothing, store.PruningStrategyEverything,
				)
			}
		}
		return nil
	}
}

// writeParamsAndConfigCmd patches the write-params cmd to additionally update the app pruning config.
func writeParamsAndConfigCmd(cdc *codec.Codec) *cobra.Command {
	cmd := kava3.WriteGenesisParamsCmd(cdc)
	originalFunc := cmd.RunE

	wrappedFunc := func(cmd *cobra.Command, args []string) error {

		if err := originalFunc(cmd, args); err != nil {
			return err
		}

		// fetch the app config from viper
		cfg, err := srvconfig.ParseConfig()
		if err != nil {
			return nil // don't return errors since as failures aren't critical
		}
		// don't prune any state, ie store everything
		cfg.Pruning = store.PruningStrategyNothing
		// write updated config
		if viper.ConfigFileUsed() == "" {
			return nil
		}
		srvconfig.WriteConfigFile(viper.ConfigFileUsed(), cfg)
		return nil
	}

	cmd.RunE = wrappedFunc
	return cmd
}

package simulation

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/kava-labs/kava/x/bep3/types"
)

// Simulation parameter constants
const (
	BnbDeputyAddress = "bnb_deputy_address"
	MinBlockLock     = "min_block_lock"
	MaxBlockLock     = "max_block_lock"
	SupportedAssets  = "supported_assets"
)

var (
	MaxSupplyLimit  = sdk.NewInt(10000000000000000)
	BondedAddresses []sdk.AccAddress
)

// GenBnbDeputyAddress randomized BnbDeputyAddress
func GenBnbDeputyAddress(r *rand.Rand) sdk.AccAddress {
	return BondedAddresses[r.Intn(len(BondedAddresses))]
}

// GenMinBlockLock randomized MinBlockLock
func GenMinBlockLock(r *rand.Rand) int64 {
	min := int(types.AbsoluteMinimumBlockLock)
	max := int(types.AbsoluteMaximumBlockLock)
	return int64(r.Intn(max-min) + min)
}

// GenMaxBlockLock randomized MaxBlockLock
func GenMaxBlockLock(r *rand.Rand, minBlockLock int64) int64 {
	min := int(minBlockLock)
	max := int(types.AbsoluteMaximumBlockLock)
	return int64(r.Intn(max-min) + min)
}

// GenSupportedAssets gets randomized SupportedAssets
func GenSupportedAssets(r *rand.Rand) types.AssetParams {
	var assets types.AssetParams
	for i := 0; i < (r.Intn(10) + 1); i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		asset := genSupportedAsset(r)
		assets = append(assets, asset)
	}
	return assets
}

func genSupportedAsset(r *rand.Rand) types.AssetParam {
	denom := strings.ToLower(simulation.RandStringOfLength(r, (r.Intn(3) + 3)))
	coinID, _ := simulation.RandPositiveInt(r, sdk.NewInt(100000))
	limit, _ := simulation.RandPositiveInt(r, MaxSupplyLimit)
	return types.AssetParam{
		Denom:  denom,
		CoinID: int(coinID.Int64()),
		Limit:  limit,
		Active: true,
	}
}

// RandomizedGenState generates a random GenesisState
func RandomizedGenState(simState *module.SimulationState) {
	BondedAddresses = loadBondedAddresses(simState)

	bep3Genesis := loadRandomBep3GenState(simState)
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, codec.MustMarshalJSONIndent(simState.Cdc, bep3Genesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(bep3Genesis)

	authGenesis, totalCoins := loadAuthGenState(simState, bep3Genesis)
	simState.GenState[auth.ModuleName] = simState.Cdc.MustMarshalJSON(authGenesis)

	// Update supply to match amount of coins in auth
	var supplyGenesis supply.GenesisState
	simState.Cdc.MustUnmarshalJSON(simState.GenState[supply.ModuleName], &supplyGenesis)
	for _, deputyCoin := range totalCoins {
		supplyGenesis.Supply = supplyGenesis.Supply.Add(deputyCoin)
	}
	simState.GenState[supply.ModuleName] = simState.Cdc.MustMarshalJSON(supplyGenesis)
}

func loadRandomBep3GenState(simState *module.SimulationState) types.GenesisState {
	var bnbDeputyAddress sdk.AccAddress
	simState.AppParams.GetOrGenerate(
		simState.Cdc, BnbDeputyAddress, &bnbDeputyAddress, simState.Rand,
		func(r *rand.Rand) { bnbDeputyAddress = GenBnbDeputyAddress(r) },
	)

	fmt.Println("simState:", simState.AppParams)
	// TODO: set minBlockLock/maxBlockLock based off sim.numBlocks
	minBlockLock := int64(types.AbsoluteMinimumBlockLock)
	// var minBlockLock int64
	// simState.AppParams.GetOrGenerate(
	// 	simState.Cdc, MinBlockLock, &minBlockLock, simState.Rand,
	// 	func(r *rand.Rand) { minBlockLock = GenMinBlockLock(r) },
	// )

	maxBlockLock := minBlockLock * 2
	// var maxBlockLock int64
	// simState.AppParams.GetOrGenerate(
	// 	simState.Cdc, MaxBlockLock, &maxBlockLock, simState.Rand,
	// 	func(r *rand.Rand) { maxBlockLock = GenMaxBlockLock(r, minBlockLock) },
	// )

	var supportedAssets types.AssetParams
	simState.AppParams.GetOrGenerate(
		simState.Cdc, SupportedAssets, &supportedAssets, simState.Rand,
		func(r *rand.Rand) { supportedAssets = GenSupportedAssets(r) },
	)

	bep3Genesis := types.GenesisState{
		Params: types.Params{
			BnbDeputyAddress: bnbDeputyAddress,
			MinBlockLock:     minBlockLock,
			MaxBlockLock:     maxBlockLock,
			SupportedAssets:  supportedAssets,
		},
	}

	return bep3Genesis
}

func loadAuthGenState(simState *module.SimulationState, bep3Genesis types.GenesisState) (
	auth.GenesisState, []sdk.Coins) {
	var authGenesis auth.GenesisState
	simState.Cdc.MustUnmarshalJSON(simState.GenState[auth.ModuleName], &authGenesis)

	deputy, found := getAccount(authGenesis.Accounts, bep3Genesis.Params.BnbDeputyAddress)
	if !found {
		panic("deputy address not found in available accounts")
	}

	// Load total limit of each supported asset to deputy's account
	var totalCoins []sdk.Coins
	for _, asset := range bep3Genesis.Params.SupportedAssets {
		assetCoin := sdk.NewCoins(sdk.NewCoin(asset.Denom, asset.Limit))
		if err := deputy.SetCoins(deputy.GetCoins().Add(assetCoin)); err != nil {
			panic(err)
		}
		totalCoins = append(totalCoins, assetCoin)
	}
	authGenesis.Accounts = replaceOrAppendAccount(authGenesis.Accounts, deputy)

	return authGenesis, totalCoins
}

// TODO: This function can be refactored
// loadBondedAddresses loads an array of bonded account addresses
func loadBondedAddresses(simState *module.SimulationState) []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, simState.NumBonded)
	for i := 0; i < int(simState.NumBonded); i++ {
		addr := simState.Accounts[i].Address
		addrs[i] = addr
	}
	return addrs
}

// Return an account from a list of accounts that matches an address.
func getAccount(accounts []authexported.GenesisAccount, addr sdk.AccAddress) (authexported.GenesisAccount, bool) {
	for _, a := range accounts {
		if a.GetAddress().Equals(addr) {
			return a, true
		}
	}
	return nil, false
}

// In a list of accounts, replace the first account found with the same address. If not found, append the account.
func replaceOrAppendAccount(accounts []authexported.GenesisAccount, acc authexported.GenesisAccount) []authexported.GenesisAccount {
	newAccounts := accounts
	for i, a := range accounts {
		if a.GetAddress().Equals(acc.GetAddress()) {
			newAccounts[i] = acc
			return newAccounts
		}
	}
	return append(newAccounts, acc)
}

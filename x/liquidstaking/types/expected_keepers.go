package types // noalias

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error

	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

// AccountKeeper defines the expected keeper interface for interacting with account
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)

	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
	SetModuleAccount(ctx sdk.Context, macc authtypes.ModuleAccountI)
}

// StakingKeeper defines the expected keeper interface for interacting with staking
type StakingKeeper interface {
	Delegate(
		ctx sdk.Context, delAddr sdk.AccAddress, bondAmt sdk.Int, tokenSrc stakingtypes.BondStatus,
		validator stakingtypes.Validator, subtractAccount bool,
	) (newShares sdk.Dec, err error)

	BondDenom(ctx sdk.Context) (res string)

	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)

	GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation stakingtypes.Delegation, found bool)

	// iterate through validators by operator address, execute func for each validator
	IterateValidators(sdk.Context,
		func(index int64, validator stakingtypes.ValidatorI) (stop bool))

	// iterate through bonded validators by operator address, execute func for each validator
	IterateBondedValidatorsByPower(sdk.Context,
		func(index int64, validator stakingtypes.ValidatorI) (stop bool))

	// iterate through the consensus validator set of the last block by operator address, execute func for each validator
	IterateLastValidators(sdk.Context,
		func(index int64, validator stakingtypes.ValidatorI) (stop bool))

	Validator(sdk.Context, sdk.ValAddress) stakingtypes.ValidatorI            // get a particular validator by operator address
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingtypes.ValidatorI // get a particular validator by consensus address

	// slash the validator and delegators of the validator, specifying offence height, offence power, and slash fraction
	Slash(sdk.Context, sdk.ConsAddress, int64, int64, sdk.Dec)
	Jail(sdk.Context, sdk.ConsAddress)   // jail a validator
	Unjail(sdk.Context, sdk.ConsAddress) // unjail a validator

	// Delegation allows for getting a particular delegation for a given validator
	// and delegator outside the scope of the staking module.
	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) stakingtypes.DelegationI

	// MaxValidators returns the maximum amount of bonded validators
	MaxValidators(sdk.Context) uint32

	ValidateUnbondAmount(
		ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, amt sdk.Int,
	) (shares sdk.Dec, err error)

	Unbond(
		ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec,
	) (amount sdk.Int, err error)

	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress,
		fn func(index int64, delegation stakingtypes.DelegationI) (stop bool))
	// TODO: staking methods

}

type DistributionKeeper interface {
	GetDelegatorWithdrawAddr(ctx sdk.Context, delAddr sdk.AccAddress) sdk.AccAddress
	WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error)
}

// LiquidStakingHooks event hooks for other keepers to run code in response to LiquidStaking modifications
type LiquidStakingHooks interface {
}

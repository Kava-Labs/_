package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/kava-labs/kava/x/incentive/types"
	validatorvesting "github.com/kava-labs/kava/x/validator-vesting"
)

// ClaimUSDXMintingReward sends the reward amount to the input receiver address and zero's out the claim in the store
func (k Keeper) ClaimUSDXMintingReward(ctx sdk.Context, owner, receiver sdk.AccAddress, multiplierName types.MultiplierName) error {
	claim, found := k.GetUSDXMintingClaim(ctx, owner)
	if !found {
		return sdkerrors.Wrapf(types.ErrClaimNotFound, "address: %s", owner)
	}

	multiplier, found := k.GetMultiplier(ctx, multiplierName)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidMultiplier, string(multiplierName))
	}

	claimEnd := k.GetClaimEnd(ctx)

	if ctx.BlockTime().After(claimEnd) {
		return sdkerrors.Wrapf(types.ErrClaimExpired, "block time %s > claim end time %s", ctx.BlockTime(), claimEnd)
	}

	claim, err := k.SynchronizeUSDXMintingClaim(ctx, claim)
	if err != nil {
		return err
	}

	rewardAmount := claim.Reward.Amount.ToDec().Mul(multiplier.Factor).RoundInt()
	if rewardAmount.IsZero() {
		return types.ErrZeroClaim
	}
	rewardCoin := sdk.NewCoin(claim.Reward.Denom, rewardAmount)
	length, err := k.GetPeriodLength(ctx, multiplier)
	if err != nil {
		return err
	}

	err = k.SendTimeLockedCoinsToAccount(ctx, types.IncentiveMacc, receiver, sdk.NewCoins(rewardCoin), length)
	if err != nil {
		return err
	}

	k.ZeroUSDXMintingClaim(ctx, claim)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(types.AttributeKeyClaimedBy, claim.GetOwner().String()),
			sdk.NewAttribute(types.AttributeKeyClaimAmount, claim.GetReward().String()),
			sdk.NewAttribute(types.AttributeKeyClaimAmount, claim.GetType()),
		),
	)
	return nil
}

// ClaimHardReward sends the reward amount to the input address and zero's out the claim in the store
func (k Keeper) ClaimHardReward(ctx sdk.Context, owner, receiver sdk.AccAddress, multiplierName types.MultiplierName) error {
	_, found := k.GetHardLiquidityProviderClaim(ctx, owner)
	if !found {
		return sdkerrors.Wrapf(types.ErrClaimNotFound, "address: %s", owner)
	}

	multiplier, found := k.GetMultiplier(ctx, multiplierName)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidMultiplier, string(multiplierName))
	}

	claimEnd := k.GetClaimEnd(ctx)

	if ctx.BlockTime().After(claimEnd) {
		return sdkerrors.Wrapf(types.ErrClaimExpired, "block time %s > claim end time %s", ctx.BlockTime(), claimEnd)
	}

	k.SynchronizeHardLiquidityProviderClaim(ctx, owner)

	claim, found := k.GetHardLiquidityProviderClaim(ctx, owner)
	if !found {
		return sdkerrors.Wrapf(types.ErrClaimNotFound, "address: %s", owner)
	}

	rewardCoins := types.MultiplyCoins(claim.Reward, multiplier.Factor)
	if rewardCoins.IsZero() {
		return types.ErrZeroClaim
	}
	length, err := k.GetPeriodLength(ctx, multiplier)
	if err != nil {
		return err
	}

	err = k.SendTimeLockedCoinsToAccount(ctx, types.IncentiveMacc, receiver, rewardCoins, length)
	if err != nil {
		return err
	}

	k.ZeroHardLiquidityProviderClaim(ctx, claim)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(types.AttributeKeyClaimedBy, claim.GetOwner().String()),
			sdk.NewAttribute(types.AttributeKeyClaimAmount, claim.GetReward().String()),
			sdk.NewAttribute(types.AttributeKeyClaimType, claim.GetType()),
		),
	)
	return nil
}

// ClaimDelegatorReward sends the reward amount to the input address and zero's out the claim in the store
func (k Keeper) ClaimDelegatorReward(ctx sdk.Context, owner, receiver sdk.AccAddress, multiplierName types.MultiplierName) error {
	claim, found := k.GetDelegatorClaim(ctx, owner)
	if !found {
		return sdkerrors.Wrapf(types.ErrClaimNotFound, "address: %s", owner)
	}

	multiplier, found := k.GetMultiplier(ctx, multiplierName)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidMultiplier, string(multiplierName))
	}

	claimEnd := k.GetClaimEnd(ctx)

	if ctx.BlockTime().After(claimEnd) {
		return sdkerrors.Wrapf(types.ErrClaimExpired, "block time %s > claim end time %s", ctx.BlockTime(), claimEnd)
	}

	syncedClaim, err := k.SynchronizeDelegatorClaim(ctx, claim)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrClaimNotFound, "address: %s", owner)
	}

	rewardCoins := types.MultiplyCoins(syncedClaim.Reward, multiplier.Factor)
	if rewardCoins.IsZero() {
		return types.ErrZeroClaim
	}
	length, err := k.GetPeriodLength(ctx, multiplier)
	if err != nil {
		return err
	}

	err = k.SendTimeLockedCoinsToAccount(ctx, types.IncentiveMacc, receiver, rewardCoins, length)
	if err != nil {
		return err
	}

	k.ZeroDelegatorClaim(ctx, syncedClaim)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(types.AttributeKeyClaimedBy, syncedClaim.GetOwner().String()),
			sdk.NewAttribute(types.AttributeKeyClaimAmount, syncedClaim.GetReward().String()),
			sdk.NewAttribute(types.AttributeKeyClaimType, syncedClaim.GetType()),
		),
	)
	return nil
}

func (k Keeper) ValidateIsValidatorVestingAccount(ctx sdk.Context, address sdk.AccAddress) error {
	acc := k.accountKeeper.GetAccount(ctx, address)
	if acc == nil {
		return sdkerrors.Wrapf(types.ErrAccountNotFound, "address not found: %s", address)
	}
	_, ok := acc.(*validatorvesting.ValidatorVestingAccount)
	if !ok {
		return sdkerrors.Wrapf(types.ErrInvalidAccountType, "account is not validator vesting account, %s", address)
	}
	return nil
}


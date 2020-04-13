package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/kava-labs/kava/x/cdp/types"
)

// DepositCollateral adds collateral to a cdp
func (k Keeper) DepositCollateral(ctx sdk.Context, owner sdk.AccAddress, depositor sdk.AccAddress, collateral sdk.Coins) error {
	err := k.ValidateCollateral(ctx, collateral)
	if err != nil {
		return err
	}
	cdp, found := k.GetCdpByOwnerAndDenom(ctx, owner, collateral[0].Denom)
	if !found {
		return sdkerrors.Wrapf(types.ErrCdpNotFound, "owner %s, collateral %s", owner, collateral[0].Denom)
	}

	deposit, found := k.GetDeposit(ctx, cdp.ID, depositor)
	if found {
		deposit.Amount = deposit.Amount.Add(collateral...)
	} else {
		deposit = types.NewDeposit(cdp.ID, depositor, collateral)
	}
	err = k.supplyKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, collateral)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCdpDeposit,
			sdk.NewAttribute(sdk.AttributeKeyAmount, collateral.String()),
			sdk.NewAttribute(types.AttributeKeyCdpID, fmt.Sprintf("%d", cdp.ID)),
		),
	)

	k.SetDeposit(ctx, deposit)

	periods := sdk.NewInt(ctx.BlockTime().Unix()).Sub(sdk.NewInt(cdp.FeesUpdated.Unix()))
	fees := k.CalculateFees(ctx, cdp.Principal.Add(cdp.AccumulatedFees...), periods, cdp.Collateral[0].Denom)
	oldCollateralToDebtRatio := k.CalculateCollateralToDebtRatio(ctx, cdp.Collateral, cdp.Principal.Add(cdp.AccumulatedFees...))
	k.RemoveCdpCollateralRatioIndex(ctx, cdp.Collateral[0].Denom, cdp.ID, oldCollateralToDebtRatio)

	cdp.AccumulatedFees = cdp.AccumulatedFees.Add(fees...)
	cdp.FeesUpdated = ctx.BlockTime()
	cdp.Collateral = cdp.Collateral.Add(collateral...)
	collateralToDebtRatio := k.CalculateCollateralToDebtRatio(ctx, cdp.Collateral, cdp.Principal.Add(cdp.AccumulatedFees...))
	k.SetCdpAndCollateralRatioIndex(ctx, cdp, collateralToDebtRatio)
	return nil
}

// WithdrawCollateral removes collateral from a cdp if it does not put the cdp below the liquidation ratio
func (k Keeper) WithdrawCollateral(ctx sdk.Context, owner sdk.AccAddress, depositor sdk.AccAddress, collateral sdk.Coins) error {
	err := k.ValidateCollateral(ctx, collateral)
	if err != nil {
		return err
	}
	cdp, found := k.GetCdpByOwnerAndDenom(ctx, owner, collateral[0].Denom)
	if !found {
		return sdkerrors.Wrapf(types.ErrCdpNotFound, "owner %s, collateral %s", owner, collateral[0].Denom)
	}
	deposit, found := k.GetDeposit(ctx, cdp.ID, depositor)
	if !found {
		return sdkerrors.Wrapf(types.ErrDepositNotFound, "depositor %s, collateral %s", depositor, collateral[0].Denom)
	}
	if collateral.IsAnyGT(deposit.Amount) {
		return sdkerrors.Wrapf(types.ErrInvalidWithdrawAmount, "owner %s, collateral %s", collateral, deposit.Amount)
	}

	periods := sdk.NewInt(ctx.BlockTime().Unix()).Sub(sdk.NewInt(cdp.FeesUpdated.Unix()))
	fees := k.CalculateFees(ctx, cdp.Principal.Add(cdp.AccumulatedFees...), periods, cdp.Collateral[0].Denom)
	collateralizationRatio, err := k.CalculateCollateralizationRatio(ctx, cdp.Collateral.Sub(collateral), cdp.Principal, cdp.AccumulatedFees.Add(fees...))
	if err != nil {
		return err
	}
	liquidationRatio := k.getLiquidationRatio(ctx, collateral[0].Denom)
	if collateralizationRatio.LT(liquidationRatio) {
		return sdkerrors.Wrapf(types.ErrInvalidCollateralRatio, "owner %s, collateral ratio %s, liquidation ration %s", collateral[0].Denom, collateralizationRatio, liquidationRatio)
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCdpWithdrawal,
			sdk.NewAttribute(sdk.AttributeKeyAmount, collateral.String()),
			sdk.NewAttribute(types.AttributeKeyCdpID, fmt.Sprintf("%d", cdp.ID)),
		),
	)

	err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, collateral)
	if err != nil {
		panic(err)
	}
	oldCollateralToDebtRatio := k.CalculateCollateralToDebtRatio(ctx, cdp.Collateral, cdp.Principal.Add(cdp.AccumulatedFees...))
	k.RemoveCdpCollateralRatioIndex(ctx, cdp.Collateral[0].Denom, cdp.ID, oldCollateralToDebtRatio)

	cdp.AccumulatedFees = cdp.AccumulatedFees.Add(fees...)
	cdp.FeesUpdated = ctx.BlockTime()
	cdp.Collateral = cdp.Collateral.Sub(collateral)
	collateralToDebtRatio := k.CalculateCollateralToDebtRatio(ctx, cdp.Collateral, cdp.Principal.Add(cdp.AccumulatedFees...))
	k.SetCdpAndCollateralRatioIndex(ctx, cdp, collateralToDebtRatio)

	deposit.Amount = deposit.Amount.Sub(collateral)
	if deposit.Amount.IsZero() {
		k.DeleteDeposit(ctx, deposit.CdpID, deposit.Depositor)
	} else {
		k.SetDeposit(ctx, deposit)
	}
	return nil
}

// GetDeposit returns the deposit of a depositor on a particular cdp from the store
func (k Keeper) GetDeposit(ctx sdk.Context, cdpID uint64, depositor sdk.AccAddress) (deposit types.Deposit, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.DepositKeyPrefix)
	bz := store.Get(types.DepositKey(cdpID, depositor))
	if bz == nil {
		return deposit, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &deposit)
	return deposit, true

}

// SetDeposit sets the deposit in the store
func (k Keeper) SetDeposit(ctx sdk.Context, deposit types.Deposit) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.DepositKeyPrefix)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(deposit)
	store.Set(types.DepositKey(deposit.CdpID, deposit.Depositor), bz)

}

// DeleteDeposit deletes a deposit from the store
func (k Keeper) DeleteDeposit(ctx sdk.Context, cdpID uint64, depositor sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.DepositKeyPrefix)
	store.Delete(types.DepositKey(cdpID, depositor))
}

// IterateDeposits iterates over the all the deposits of a cdp and performs a callback function
func (k Keeper) IterateDeposits(ctx sdk.Context, cdpID uint64, cb func(deposit types.Deposit) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.DepositKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, types.GetCdpIDBytes(cdpID))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var deposit types.Deposit
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &deposit)

		if cb(deposit) {
			break
		}
	}
}

// GetDeposits returns all the deposits to a cdp
func (k Keeper) GetDeposits(ctx sdk.Context, cdpID uint64) (deposits types.Deposits) {
	k.IterateDeposits(ctx, cdpID, func(deposit types.Deposit) bool {
		deposits = append(deposits, deposit)
		return false
	})
	return
}

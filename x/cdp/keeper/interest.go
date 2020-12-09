package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/kava-labs/kava/x/cdp/types"
)

var (
	scalingFactor  = 1e18
	secondsPerYear = 31536000
)

// AccumulateInterest calculates the new interest that has accrued for the input collateral type based on the total amount of principal
// that has been created with that collateral type and the amount of time that has passed since interest was last accumulated
func (k Keeper) AccumulateInterest(ctx sdk.Context, ctype string) error {
	previousAccrualTime, found := k.GetPreviousAccrualTime(ctx, ctype)
	if !found {
		k.SetPreviousAccrualTime(ctx, ctype, ctx.BlockTime())
		return nil
	}

	timeElapsed := ctx.BlockTime().Unix() - previousAccrualTime.Unix()
	if timeElapsed == 0 {
		return nil
	}

	totalPrincipalPrior := k.GetTotalPrincipal(ctx, ctype, types.DefaultStableDenom)
	if totalPrincipalPrior.IsZero() || totalPrincipalPrior.IsNegative() {
		k.SetPreviousAccrualTime(ctx, ctype, ctx.BlockTime())
		return nil
	}

	interestFactorPrior, foundInterestFactorPrior := k.GetInterestFactor(ctx, ctype)
	if !foundInterestFactorPrior {
		k.SetInterestFactor(ctx, ctype, sdk.OneDec())
		// set previous accrual time exit early because interest accumulated will be zero
		k.SetPreviousAccrualTime(ctx, ctype, ctx.BlockTime())
		return nil
	}

	borrowRateSpy := k.getFeeRate(ctx, ctype)
	if borrowRateSpy.Equal(sdk.OneDec()) {
		k.SetPreviousAccrualTime(ctx, ctype, ctx.BlockTime())
		return nil
	}
	interestFactor := CalculateInterestFactor(borrowRateSpy, sdk.NewInt(timeElapsed))
	interestAccumulated := (interestFactor.Mul(totalPrincipalPrior.ToDec())).RoundInt().Sub(totalPrincipalPrior)
	if interestAccumulated.IsZero() {
		// in the case accumulated interest rounds to zero, exit early without updating accrual time
		return nil
	}
	err := k.MintDebtCoins(ctx, types.ModuleName, k.GetDebtDenom(ctx), sdk.NewCoin(types.DefaultStableDenom, interestAccumulated))
	if err != nil {
		return err
	}

	interestFactorNew := interestFactorPrior.Mul(interestFactor)
	totalPrincipalNew := totalPrincipalPrior.Add(interestAccumulated)

	k.SetTotalPrincipal(ctx, ctype, types.DefaultStableDenom, totalPrincipalNew)
	k.SetInterestFactor(ctx, ctype, interestFactorNew)
	k.SetPreviousAccrualTime(ctx, ctype, ctx.BlockTime())

	return nil
}

// CalculateInterestFactor calculates the simple interest scaling factor,
// which is equal to: (per-second interest rate ** number of seconds elapsed)
// Will return 1.000x, multiply by principal to get new principal with added interest
func CalculateInterestFactor(perSecondInterestRate sdk.Dec, secondsElapsed sdk.Int) sdk.Dec {
	scalingFactorUint := sdk.NewUint(uint64(scalingFactor))
	scalingFactorInt := sdk.NewInt(int64(scalingFactor))

	// Convert per-second interest rate to a uint scaled by 1e18
	interestMantissa := sdk.NewUint(perSecondInterestRate.MulInt(scalingFactorInt).RoundInt().Uint64())
	// Convert seconds elapsed to uint (*not scaled*)
	secondsElapsedUint := sdk.NewUint(secondsElapsed.Uint64())
	// Calculate the interest factor as a uint scaled by 1e18
	interestFactorMantissa := sdk.RelativePow(interestMantissa, secondsElapsedUint, scalingFactorUint)

	// Convert interest factor to an unscaled sdk.Dec
	return sdk.NewDecFromBigInt(interestFactorMantissa.BigInt()).QuoInt(scalingFactorInt)
}

// SynchronizeInterest updates the input cdp object to reflect the current accumulated interest, updates the cdp state in the store,
// and returns the updated cdp object
func (k Keeper) SynchronizeInterest(ctx sdk.Context, cdp types.CDP) types.CDP {
	globalInterestFactor, found := k.GetInterestFactor(ctx, cdp.Type)
	if !found {
		k.SetInterestFactor(ctx, cdp.Type, sdk.OneDec())
		cdp.InterestFactor = sdk.OneDec()
		cdp.FeesUpdated = ctx.BlockTime()
		k.SetCDP(ctx, cdp)
	}
	cdpInterestFactor := sdk.OneDec().Add((globalInterestFactor.Sub(cdp.InterestFactor)))
	if cdpInterestFactor.Equal(sdk.OneDec()) {
		return cdp
	}
	accumulatedInterest := cdp.GetTotalPrincipal().Amount.ToDec().Mul(cdpInterestFactor).RoundInt().Sub(cdp.GetTotalPrincipal().Amount)

	cdp.AccumulatedFees = cdp.AccumulatedFees.Add(sdk.NewCoin(cdp.AccumulatedFees.Denom, accumulatedInterest))
	cdp.FeesUpdated = ctx.BlockTime()
	cdp.InterestFactor = globalInterestFactor
	collateralToDebtRatio := k.CalculateCollateralToDebtRatio(ctx, cdp.Collateral, cdp.Type, cdp.GetTotalPrincipal())
	k.SetCdpAndCollateralRatioIndex(ctx, cdp, collateralToDebtRatio)
	return cdp
}

// CalculateNewInterest returns the amount of interest that has accrued to the cdp since its interest was last synchronized
func (k Keeper) CalculateNewInterest(ctx sdk.Context, cdp types.CDP) sdk.Coin {
	globalInterestFactor, found := k.GetInterestFactor(ctx, cdp.Type)
	if !found {
		return sdk.NewCoin(cdp.AccumulatedFees.Denom, sdk.ZeroInt())
	}
	cdpInterestFactor := sdk.OneDec().Add((globalInterestFactor.Sub(cdp.InterestFactor)))
	if cdpInterestFactor.Equal(sdk.OneDec()) {
		return sdk.NewCoin(cdp.AccumulatedFees.Denom, sdk.ZeroInt())
	}
	accumulatedInterest := cdp.GetTotalPrincipal().Amount.ToDec().Mul(cdpInterestFactor).RoundInt().Sub(cdp.GetTotalPrincipal().Amount)
	return sdk.NewCoin(cdp.AccumulatedFees.Denom, accumulatedInterest)
}

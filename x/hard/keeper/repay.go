package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/kava-labs/kava/x/hard/types"
)

// Repay borrowed funds
func (k Keeper) Repay(ctx sdk.Context, sender, owner sdk.AccAddress, coins sdk.Coins) error {
	// Check borrow exists here to avoid duplicating store read in ValidateRepay
	borrow, found := k.GetBorrow(ctx, owner)
	if !found {
		return types.ErrBorrowNotFound
	}
	// Call incentive hook
	k.BeforeBorrowModified(ctx, borrow)

	// Sync borrow interest so loan is up-to-date
	k.SyncBorrowInterest(ctx, owner)

	// Validate that sender holds coins for repayment
	err := k.ValidateRepay(ctx, sender, owner, coins)
	if err != nil {
		return err
	}

	payment, err := k.CalculatePaymentAmount(borrow.Amount, coins)
	if err != nil {
		return err
	}

	// Sends coins from user to Hard module account
	err = k.supplyKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleAccountName, payment)
	if err != nil {
		return err
	}

	// If any coin denoms have been completely repaid reset the denom's borrow index factor
	for _, coin := range payment {
		if coin.Amount.Equal(borrow.Amount.AmountOf(coin.Denom)) {
			borrowIndex, removed := borrow.Index.RemoveInterestFactor(coin.Denom)
			if !removed {
				return sdkerrors.Wrapf(types.ErrInvalidIndexFactorDenom, "%s", coin.Denom)
			}
			borrow.Index = borrowIndex
		}
	}

	// Update user's borrow in store
	borrow.Amount = borrow.Amount.Sub(payment)

	if borrow.Amount.Empty() {
		k.DeleteBorrow(ctx, borrow)
	} else {
		k.SetBorrow(ctx, borrow)
	}

	// Update total borrowed amount
	k.DecrementBorrowedCoins(ctx, payment)

	// Call incentive hook
	if !borrow.Amount.Empty() {
		k.AfterBorrowModified(ctx, borrow)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeHardRepay,
			sdk.NewAttribute(types.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, owner.String()),
			sdk.NewAttribute(types.AttributeKeyRepayCoins, payment.String()),
		),
	)

	return nil
}

// ValidateRepay validates a requested loan repay
func (k Keeper) ValidateRepay(ctx sdk.Context, sender, owner sdk.AccAddress, coins sdk.Coins) error {
	moneyMarketCache := map[string]types.MoneyMarket{}
	assetPriceCache := map[string]sdk.Dec{}

	// Get the total USD value of user's existing borrows
	existingBorrowUSDValue := sdk.ZeroDec()
	existingBorrow, found := k.GetBorrow(ctx, owner)
	if found {
		for _, borrowedCoin := range existingBorrow.Amount {
			moneyMarket, ok := moneyMarketCache[borrowedCoin.Denom]
			if !ok { // Fetch money market and store in local cache
				newMoneyMarket, found := k.GetMoneyMarketParam(ctx, borrowedCoin.Denom)
				if !found {
					return sdkerrors.Wrapf(types.ErrMarketNotFound, "no market found for denom %s", borrowedCoin.Denom)
				}
				moneyMarketCache[borrowedCoin.Denom] = newMoneyMarket
				moneyMarket = newMoneyMarket
			}

			assetPrice, ok := assetPriceCache[borrowedCoin.Denom]
			if !ok { // Fetch current asset price and store in local cache
				assetPriceInfo, err := k.pricefeedKeeper.GetCurrentPrice(ctx, moneyMarket.SpotMarketID)
				if err != nil {
					return sdkerrors.Wrapf(types.ErrPriceNotFound, "no price found for market %s", moneyMarket.SpotMarketID)
				}
				assetPriceCache[borrowedCoin.Denom] = assetPriceInfo.Price
				assetPrice = assetPriceInfo.Price
			}

			// Calculate this borrow coin's USD value and add it to the total previous borrowed USD value
			coinUSDValue := sdk.NewDecFromInt(borrowedCoin.Amount).Quo(sdk.NewDecFromInt(moneyMarket.ConversionFactor)).Mul(assetPrice)
			existingBorrowUSDValue = existingBorrowUSDValue.Add(coinUSDValue)
		}
	}

	senderAcc := k.accountKeeper.GetAccount(ctx, sender)
	senderCoins := senderAcc.SpendableCoins(ctx.BlockTime())
	repayTotalUSDValue := sdk.ZeroDec()
	for _, repayCoin := range coins {
		// Check that sender holds enough tokens to make the proposed payment
		if senderCoins.AmountOf(repayCoin.Denom).LT(repayCoin.Amount) {
			return sdkerrors.Wrapf(types.ErrInsufficientBalanceForRepay, "account can only repay up to %s%s", senderCoins.AmountOf(repayCoin.Denom), repayCoin.Denom)
		}

		moneyMarket, ok := moneyMarketCache[repayCoin.Denom]
		if !ok { // Fetch money market and store in local cache
			newMoneyMarket, found := k.GetMoneyMarketParam(ctx, repayCoin.Denom)
			if !found {
				return sdkerrors.Wrapf(types.ErrMarketNotFound, "no market found for denom %s", repayCoin.Denom)
			}
			moneyMarketCache[repayCoin.Denom] = newMoneyMarket
			moneyMarket = newMoneyMarket
		}

		// Calculate this coin's USD value and add it to the repay's total USD value
		assetPrice, ok := assetPriceCache[repayCoin.Denom]
		if !ok { // Fetch current asset price and store in local cache
			assetPriceInfo, err := k.pricefeedKeeper.GetCurrentPrice(ctx, moneyMarket.SpotMarketID)
			if err != nil {
				return sdkerrors.Wrapf(types.ErrPriceNotFound, "no price found for market %s", moneyMarket.SpotMarketID)
			}
			assetPriceCache[repayCoin.Denom] = assetPriceInfo.Price
			assetPrice = assetPriceInfo.Price
		}
		coinUSDValue := sdk.NewDecFromInt(repayCoin.Amount).Quo(sdk.NewDecFromInt(moneyMarket.ConversionFactor)).Mul(assetPrice)
		repayTotalUSDValue = repayTotalUSDValue.Add(coinUSDValue)
	}

	// If the proposed repayment would results in a borrowed USD value below the minimum borrow USD value, reject it.
	// User can overpay their loan to close it out, but underpaying by such a margin that the USD value is in an
	// invalid range is not allowed.
	proposedBorrowNewUSDValue := existingBorrowUSDValue.Sub(repayTotalUSDValue)
	if proposedBorrowNewUSDValue.IsPositive() && proposedBorrowNewUSDValue.LT(k.GetMinimumBorrowUSDValue(ctx)) {
		return sdkerrors.Wrapf(types.ErrBelowMinimumBorrowValue, "the proposed borrow's USD value $%s is below the minimum borrow limit $%s", proposedBorrowNewUSDValue, k.GetMinimumBorrowUSDValue(ctx))
	}

	return nil
}

// CalculatePaymentAmount prevents overpayment when repaying borrowed coins
func (k Keeper) CalculatePaymentAmount(owed sdk.Coins, payment sdk.Coins) (sdk.Coins, error) {
	repayment := sdk.Coins{}

	if !payment.DenomsSubsetOf(owed) {
		return repayment, types.ErrInvalidRepaymentDenom
	}

	for _, coin := range payment {
		if coin.Amount.GT(owed.AmountOf(coin.Denom)) {
			repayment = append(repayment, sdk.NewCoin(coin.Denom, owed.AmountOf(coin.Denom)))
		} else {
			repayment = append(repayment, coin)
		}
	}
	return repayment, nil
}

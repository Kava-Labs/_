package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
)

// GetTotalVestingPeriodLength returns the summed length of all vesting periods
func GetTotalVestingPeriodLength(periods vesting.Periods) int64 {
	length := int64(0)
	for _, period := range periods {
		length += period.Length
	}
	return length
}

// MultiplyCoins multiplies each value in a set of coins by a single decimal value, rounding the result.
func MultiplyCoins(coins sdk.Coins, multiple sdk.Dec) sdk.Coins {
	var result sdk.Coins
	for _, coin := range coins {
		result = result.Add(
			sdk.NewCoin(coin.Denom, coin.Amount.ToDec().Mul(multiple).RoundInt()),
		)
	}
	return result
}

package types

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	"math"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// An Accumulator handles calculating and tracking global reward distributions.
type Accumulator struct {
	PreviousAccumulationTime time.Time
	Indexes                  RewardIndexes
}

func NewAccumulator(previousAccrual time.Time, indexes RewardIndexes) *Accumulator {
	return &Accumulator{
		PreviousAccumulationTime: previousAccrual,
		Indexes:                  indexes,
	}
}

// Accumulate accrues rewards up to the current time.
//
// It calculates new rewards and adds them to the reward indexes for the period from PreviousAccumulationTime to currentTime.
// It stores the currentTime in PreviousAccumulationTime to be used for later accumulations.
//
// Rewards are not accrued for times outside of the start and end times of a reward period.
// If a period ends before currentTime, the PreviousAccrualTime is shortened to the end time. This allows accumulate to be called sequentially on consecutive reward periods.
//
// totalSourceShares is the sum of all users' source shares. For example:total btcb supplied to hard, total usdx borrowed from all bnb CDPs, or total shares in a swap pool.
func (acc *Accumulator) Accumulate(period MultiRewardPeriod, totalSourceShares sdkmath.LegacyDec, currentTime time.Time) {
	fmt.Println("Accumulate: ", period, totalSourceShares, currentTime)
	acc.AccumulateDecCoins(
		period.Start,
		period.End,
		sdk.NewDecCoinsFromCoins(period.RewardsPerSecond...),
		totalSourceShares,
		currentTime,
	)
}

// AccumulateDecCoins
func (acc *Accumulator) AccumulateDecCoins(
	periodStart time.Time,
	periodEnd time.Time,
	periodRewardsPerSecond sdk.DecCoins,
	totalSourceShares sdkmath.LegacyDec,
	currentTime time.Time,
) {
	fmt.Println("AccumulateDecCoins: ", periodStart, periodEnd, periodRewardsPerSecond, totalSourceShares, currentTime)
	accumulationDuration := acc.getTimeElapsedWithinLimits(acc.PreviousAccumulationTime, currentTime, periodStart, periodEnd)
	fmt.Println("accumulationDuration: ", accumulationDuration)

	indexesIncrement := acc.calculateNewRewards(periodRewardsPerSecond, totalSourceShares, accumulationDuration)
	fmt.Println("indexesIncrement: ", indexesIncrement)

	acc.Indexes = acc.Indexes.Add(indexesIncrement)
	acc.PreviousAccumulationTime = minTime(periodEnd, currentTime)
}

// getTimeElapsedWithinLimits returns the duration between start and end times, capped by min and max times.
// If the start and end range is outside the min to max time range then zero duration is returned.
func (*Accumulator) getTimeElapsedWithinLimits(start, end, limitMin, limitMax time.Time) time.Duration {
	fmt.Println("getTimeElapsedWithinLimits: ", start, end, limitMin, limitMax)
	if start.After(end) {
		panic(fmt.Sprintf("start time (%s) cannot be after end time (%s)", start, end))
	}
	if limitMin.After(limitMax) {
		panic(fmt.Sprintf("minimum limit time (%s) cannot be after maximum limit time (%s)", limitMin, limitMax))
	}
	if start.After(limitMax) || end.Before(limitMin) {
		// no intersection between the start-end and limitMin-limitMax time ranges
		return 0
	}
	return minTime(end, limitMax).Sub(maxTime(start, limitMin))
}

// calculateNewRewards calculates the amount to increase the global reward indexes by, for a given reward rate, duration, and number of source shares.
// The total rewards to distribute in this block are given by reward rate * duration. This value divided by the sum of all source shares to give
// total rewards per source share, which is what the indexes store.
// Note, duration is rounded to the nearest second to keep rewards calculation consistent with kava-7.
func (*Accumulator) calculateNewRewards(rewardsPerSecond sdk.DecCoins, totalSourceShares sdkmath.LegacyDec, duration time.Duration) RewardIndexes {
	fmt.Println("calculateNewRewards: ", rewardsPerSecond, totalSourceShares, duration)
	if totalSourceShares.LTE(sdkmath.LegacyZeroDec()) {
		// When there is zero source shares, there is no users with deposits/borrows/delegations to pay out the current block's rewards to.
		// So drop the rewards and pay out nothing.
		return nil
	}
	durationSeconds := int64(math.RoundToEven(duration.Seconds()))
	if durationSeconds <= 0 {
		// If the duration is zero, there will be no increment.
		// So return an empty increment instead of one full of zeros.
		return nil
	}
	increment := NewRewardIndexesFromCoins(rewardsPerSecond)
	fmt.Println("increment: ", increment)
	fmt.Println("durationSeconds: ", durationSeconds)
	fmt.Println("totalSourceShares: ", totalSourceShares)
	increment = increment.Mul(sdkmath.LegacyNewDec(durationSeconds)).Quo(totalSourceShares)
	return increment
}

// minTime returns the earliest of two times.
func minTime(t1, t2 time.Time) time.Time {
	if t2.Before(t1) {
		return t2
	}
	return t1
}

// maxTime returns the latest of two times.
func maxTime(t1, t2 time.Time) time.Time {
	if t2.After(t1) {
		return t2
	}
	return t1
}

// NewRewardIndexesFromCoins is a helper function to initialize a RewardIndexes slice with the values from a Coins slice.
func NewRewardIndexesFromCoins(coins sdk.DecCoins) RewardIndexes {
	var indexes RewardIndexes
	for _, coin := range coins {
		fmt.Println("NewRewardIndexesFromCoins: ", coin.Denom, coin.Amount)
		indexes = append(indexes, NewRewardIndex(coin.Denom, coin.Amount))
	}
	return indexes
}

func CalculatePerSecondRewards(
	periodStart time.Time,
	periodEnd time.Time,
	periodRewardsPerSecond sdk.DecCoins,
	previousTime, currentTime time.Time,
) (sdk.DecCoins, time.Time) {
	duration := (&Accumulator{}).getTimeElapsedWithinLimits(
		previousTime,
		currentTime,
		periodStart,
		periodEnd,
	)

	upTo := minTime(periodEnd, currentTime)

	durationSeconds := int64(math.RoundToEven(duration.Seconds()))
	if durationSeconds <= 0 {
		// If the duration is zero, there will be no increment.
		// So return an empty increment instead of one full of zeros.
		return nil, upTo // TODO
	}

	return periodRewardsPerSecond.MulDec(sdkmath.LegacyNewDec(durationSeconds)), upTo
}

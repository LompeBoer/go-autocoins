package binance

import (
	"log"
	"time"
)

const (
	TickerWeight           = 40.0
	ExchangeInfoWeight     = 10.0
	KlineWeight            = 3.0  // This is combined for all the Kline requests AutoCoins does per symbol.
	WeightEstimationBuffer = 1.2  // WeightEstimationBuffer percentage of `EstimatedWeightUsage` to use in rate limit.
	MinimumWeightLimit     = 0.5  // MinimumWeightLimit percentage for minimum weight limit warning.
	MaximumWeightLimit     = 0.75 // MaximumWeightLimit percentage for maximum weight limit warning.
)

// CheckForWeightLimit check if the rate limit estimation will not exceed the set limit.
// There should be a buffer left so the WH bot still has room to do its thing.
func (a *API) CheckForWeightLimit() float64 {
	limit := float64(a.EstimatedWeightUsage) * WeightEstimationBuffer
	if a.WeightLimit != 0 {
		maxLimit := float64(a.WeightLimit) * MaximumWeightLimit
		minLimit := float64(a.WeightLimit) * MinimumWeightLimit
		if limit > maxLimit {
			limit = maxLimit
		} else if limit < minLimit {
			limit = minLimit
		}
	}
	return limit
}

// PreCheckForWeightLimit determines ahead of time if the rate limit will be exceeded.
func (a *API) PreCheckForWeightLimit() bool {
	limit := a.CheckForWeightLimit()
	log.Printf("Binance API Weight - Used: %d Estimated: %d Limit: %.0f\n", a.UsedWeight, a.EstimatedWeightUsage, limit)
	return float64(a.EstimatedWeightUsage+a.UsedWeight) > limit
}

// PauseForWeightWarning sleeps until the rate limit is reset.
func (a *API) PauseForWeightWarning() {
	a.pauseRequest(time.Minute)
}

// RateLimitChecks sets the rate limit estimation and pauses execution when estimated weight will be exceeded.
func (a *API) RateLimitChecks(symbolCount int) {
	a.EstimatedWeightUsage = int((float64(symbolCount) * KlineWeight) + TickerWeight + ExchangeInfoWeight)
	if a.PreCheckForWeightLimit() {
		log.Println("Weight warning! Will pause for one minute")
		a.PauseForWeightWarning()
		log.Println("Finished weight wait")
	}
}

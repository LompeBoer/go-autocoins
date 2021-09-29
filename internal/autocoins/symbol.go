package autocoins

import (
	"math"
	"time"

	"github.com/LompeBoer/go-autocoins/internal/exchange/binance"
)

type ExchangeData struct {
	Prices24Hours *[]binance.Ticker
	Kline1Minute  []binance.KLine
	Kline1Month   []binance.KLine
	Candles       int
}

type SymbolDataValues struct {
	Percent1Hour  []float64 `json:"perc1hrVal"`
	Percent4Hour  float64   `json:"perc4hrVal"`
	Percent24Hour float64   `json:"perc24hrVal"`
	AllTimeHigh   float64   `json:"AthVal"`
	Age           int       `json:"AgeVal"`
}

type SymbolDataResult struct {
	Percent1Hour  bool `json:"perc1hr"`
	Percent4Hour  bool `json:"perc4hr"`
	Percent24Hour bool `json:"perc24hr"`
	AllTimeHigh   bool `json:"Ath"`
	Age           bool `json:"Age"`
}

type SymbolDataObject struct {
	Symbol    binance.Symbol   `json:"symbol"`
	Open      bool             `json:"Open"`
	Time      time.Time        `json:"dateTime"`
	APIFailed bool             `json:"apiFailed"`
	Excluded  bool             `json:"excluded"`
	Values    SymbolDataValues `json:"values"`
	Result    SymbolDataResult `json:"result"`
	data      ExchangeData
	settings  *SettingsAutoCoins
}

func (a *AutoCoins) RetrieveSymbolData(symbol binance.Symbol, prices24Hours *[]binance.Ticker, c chan SymbolDataObject) {
	minCandles := 4
	if a.Settings.AutoCoins.CooldownHours >= 4 {
		minCandles = a.Settings.AutoCoins.CooldownHours
	}
	dateTime := time.Now()
	limit := minCandles * 60
	kline1Minute, err := a.ExchangeAPI.GetKLine(symbol, limit, binance.OneMinute)
	if err != nil {
		c <- a.apiFailResult(symbol)
		return
	}

	start := time.Unix(symbol.OnboardDate/1000, 0)
	age := time.Since(start).Hours() / 24.0
	limit2 := math.Round((age / 30) + 1)

	kline1Month, err := a.ExchangeAPI.GetKLine(symbol, int(limit2), binance.OneMonth)
	if err != nil {
		c <- a.apiFailResult(symbol)
		return
	}

	object := SymbolDataObject{
		Symbol: symbol,
		Time:   dateTime,
		data: ExchangeData{
			Prices24Hours: prices24Hours,
			Kline1Minute:  kline1Minute,
			Kline1Month:   kline1Month,
			Candles:       minCandles,
		},
		settings: &a.Settings.AutoCoins,
		Values: SymbolDataValues{
			Age: int(age),
		},
	}
	object.Calculate()
	c <- object
}

func (s *SymbolDataObject) Calculate() {
	prices1Hour, err := binance.Get1HourPrices(s.data.Kline1Minute)
	if err != nil {
		s.APIFailed = true
		return
	}
	percent1Hour := []float64{}
	for i := 1; i < s.data.Candles+1; i++ {
		end := i*60 - 1
		start := end - 59
		if end > len(prices1Hour)-1 {
			s.APIFailed = true
			return
		}
		percent := ((prices1Hour[end] - prices1Hour[start]) * 100) / prices1Hour[end]
		percent1Hour = append(percent1Hour, percent)
	}
	current4HoursPercent := ((prices1Hour[239] - prices1Hour[0]) * 100) / prices1Hour[239]
	current24HoursPercent := binance.CalculateCurrent24HourPercent(*s.data.Prices24Hours, s.Symbol.Name)

	// Get age and max all time high
	ath := binance.GetMaximumAllTimeHigh(s.data.Kline1Month)
	currentPercentageATH := ((ath - prices1Hour[len(prices1Hour)-1]) * 100 / ath)

	s.calculateResults(percent1Hour, current4HoursPercent, current24HoursPercent, currentPercentageATH)
}

func (s *SymbolDataObject) calculateResults(percent1Hour []float64, current4HoursPercent float64, current24HoursPercent float64, currentPercentageATH float64) {
	// 1 hour percent
	x := s.settings.CooldownHours - 1
	values1HrPercent := percent1Hour[:x]
	max := 0.0
	for _, val := range values1HrPercent {
		v := math.Abs(val)
		if v > max {
			max = v
		}
	}
	result1HrPercent := max < float64(s.settings.Max1hrPercent)

	// 4 hour percent
	result4HrPercent := math.Abs(current4HoursPercent) < float64(s.settings.Max4hrPercent)
	values4HrPercent := current4HoursPercent

	// 24 hour percent
	result24HrPercent := math.Abs(current24HoursPercent) < float64(s.settings.Max24hrPercent)
	values24HrPercent := current24HoursPercent

	resultAth := currentPercentageATH > float64(s.settings.MinAthPercent)
	valuesAth := currentPercentageATH

	resultAge := s.Values.Age > s.settings.MinAge

	s.Result.Percent1Hour = result1HrPercent
	s.Values.Percent1Hour = values1HrPercent
	s.Result.Percent4Hour = result4HrPercent
	s.Values.Percent4Hour = values4HrPercent
	s.Result.Percent24Hour = result24HrPercent
	s.Values.Percent24Hour = values24HrPercent
	s.Result.AllTimeHigh = resultAth
	s.Values.AllTimeHigh = valuesAth
	s.Result.Age = resultAge
	s.Open = false
	s.APIFailed = false
}

func (s *SymbolDataObject) ShouldQuarantine() bool {
	return !s.Result.Percent1Hour || !s.Result.Percent24Hour || !s.Result.Percent4Hour || !s.Result.AllTimeHigh || !s.Result.Age
}

func (a *AutoCoins) apiFailResult(symbol binance.Symbol) SymbolDataObject {
	return SymbolDataObject{
		Symbol:    symbol,
		APIFailed: true,
	}
}

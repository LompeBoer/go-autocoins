package autocoins

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/LompeBoer/go-autocoins/internal/binance"
	"github.com/LompeBoer/go-autocoins/internal/pairslist"
	"github.com/LompeBoer/go-autocoins/internal/wickhunter"
)

type AutoCoins struct {
	Settings                   Settings
	ExchangeAPI                *binance.API
	BotAPI                     *wickhunter.API
	ctx                        context.Context
	cancel                     context.CancelFunc
	wg                         sync.WaitGroup
	IsRunning                  bool
	MaxFailedSymbolsPercentage float64
	StorageFilename            string
	DisableWrite               bool
	OutputWriter               OutputWriter
}

// GetInfo retrieves all symbol data and calculates market swing.
// It returns a list of permitted coins to trade.
func (a *AutoCoins) GetInfo(pairsList []pairslist.Pair) ([]SymbolDataObject, SymbolLists, error) {
	exchangeInfo, err := a.ExchangeAPI.GetExchangeInfo()
	if err != nil {
		return nil, SymbolLists{}, err
	}
	usedSymbols, err := a.BotAPI.GetPositions()
	if err != nil {
		return nil, SymbolLists{}, fmt.Errorf("botapi:getpositions: %s", err.Error())
	}
	symbols, err := a.filterSymbols(usedSymbols, exchangeInfo.Symbols, pairsList)
	if err != nil {
		return nil, SymbolLists{}, fmt.Errorf("unable to filter symbols: %s", err.Error())
	}
	sort.Sort(binance.BySymbolName(symbols))

	// Will pause execution when rate limit will be exceeded.
	a.ExchangeAPI.RateLimitChecks(len(symbols))

	prices24Hours, err := a.ExchangeAPI.GetTicker()
	if err != nil {
		return nil, SymbolLists{}, err
	}

	c := make(chan SymbolDataObject)
	count := a.CalculateAllSymbolData(symbols, prices24Hours, c)

	objects := []SymbolDataObject{}
	for i := 0; i < count; i++ {
		object := <-c
		objects = append(objects, object)
	}

	positions, err := a.BotAPI.GetPositions()
	if err != nil {
		return nil, SymbolLists{}, fmt.Errorf("botapi:getpositions: %s", err.Error())
	}

	lists, err := a.makeLists(objects, positions)
	if err != nil {
		return nil, SymbolLists{}, fmt.Errorf("unable to make list: %s", err.Error())
	}

	// If not enough symbol data is retrieved from the API fail this run.
	percentageFailed := float64(len(lists.FailedToProcess)) / float64(len(symbols))
	if percentageFailed > a.MaxFailedSymbolsPercentage {
		return nil, SymbolLists{}, fmt.Errorf("unable to retrieve enough data from Binance API (%.0f%% failed)", percentageFailed*100)
	}

	return objects, lists, nil
}

// SymbolLists contains all the calculated lists.
type SymbolLists struct {
	Quarantined          []string // Quarantined symbols to quarantine.
	QuarantinedNew       []string // QuarantinedNew newly added symbols to the quarantine list.
	QuarantinedSkipped   []string // QuarantinedSkipped should be quarantined but have currently open trade.
	QuarantinedCurrently []string // QuarantinedCurrently already quarantined.
	QuarantinedRemoved   []string // QuarantinedRemoved no longer quarantined.
	Permitted            []string // Permitted symbols allowed to trade.
	PermittedCurrently   []string // PermittedCurrently symbols that were already being traded.
	FailedToProcess      []string // FailedToProcess symbols that failed to retrieve enough data to make calculations.
	NotTrading           []string // NotTrading coins that are excluded from trading.
}

// makeLists makes the SymbolLists object, this groups all the symbols in a certain list.
func (a *AutoCoins) makeLists(objects []SymbolDataObject, positions []wickhunter.Position) (SymbolLists, error) {
	openPositions := []string{}
	permittedCurrently := []string{}
	quarantinedCurrently := []string{}
	notTrading := []string{}
	for _, p := range positions {
		if p.IsOpen() {
			openPositions = append(openPositions, p.Symbol)
		}
		if p.Permitted {
			permittedCurrently = append(permittedCurrently, p.Symbol)
		} else {
			quarantinedCurrently = append(quarantinedCurrently, p.Symbol)
		}
	}

	// Quarantined / Permitted
	quarantined := []string{}
	permitted := []string{}
	quarantinedSkipped := []string{} // Skipped because it is currently being traded.
	failed := []string{}
	for _, object := range objects {
		for _, openSymbol := range openPositions {
			if openSymbol == object.Symbol.Name {
				object.Open = true
				break
			}
		}
		if object.APIFailed {
			failed = append(failed, object.Symbol.Name)

			if object.Open {
				// Open trades should still be permitted.
				permitted = append(permitted, object.Symbol.Name)
			}
		} else if !object.Percent1Hour || !object.Percent24Hour || !object.Percent4Hour || !object.AllTimeHigh || !object.Age {
			if !object.Open {
				quarantined = append(quarantined, object.Symbol.Name)
			} else {
				quarantinedSkipped = append(quarantinedSkipped, object.Symbol.Name)

				// Open trades should still be permitted.
				permitted = append(permitted, object.Symbol.Name)
			}
		} else {
			permitted = append(permitted, object.Symbol.Name)
		}
	}

	quarantinedNew := []string{}
	for _, q := range quarantined {
		found := false
		for _, qc := range quarantinedCurrently {
			if q == qc {
				found = true
				break
			}
		}
		if !found {
			quarantinedNew = append(quarantinedNew, q)
		}
	}

	quarantinedRemoved := []string{}
	for _, p := range permitted {
		found := false
		for _, pc := range permittedCurrently {
			if p == pc {
				found = true
				break
			}
		}
		if !found {
			quarantinedRemoved = append(quarantinedRemoved, p)
		}
	}

	for _, p := range positions {
		found := false
		for _, pem := range permitted {
			if p.Symbol == pem {
				found = true
				break
			}
		}

		if !found {
			notTrading = append(notTrading, p.Symbol)
		}
	}

	sort.Strings(quarantined)
	sort.Strings(quarantinedNew)
	sort.Strings(quarantinedSkipped)
	sort.Strings(quarantinedCurrently)
	sort.Strings(quarantinedRemoved)
	sort.Strings(permitted)
	sort.Strings(permittedCurrently)
	sort.Strings(failed)
	sort.Strings(notTrading)

	return SymbolLists{
		Quarantined:          quarantined,
		QuarantinedNew:       quarantinedNew,
		QuarantinedSkipped:   quarantinedSkipped,
		QuarantinedCurrently: quarantinedCurrently,
		QuarantinedRemoved:   quarantinedRemoved,
		Permitted:            permitted,
		PermittedCurrently:   permittedCurrently,
		FailedToProcess:      failed,
		NotTrading:           notTrading,
	}, nil
}

func (a *AutoCoins) CalculateAllSymbolData(symbols []binance.Symbol, prices24Hours []binance.Ticker, c chan SymbolDataObject) int {
	count := 0
	for _, symbol := range symbols {
		go a.CalculateSymbolData(symbol, prices24Hours, c)
		count++
	}
	return count
}

func (a *AutoCoins) CalculateSymbolData(symbol binance.Symbol, prices24Hours []binance.Ticker, c chan SymbolDataObject) {
	minCandles := 4
	if a.Settings.CooldownHours >= 4 {
		minCandles = a.Settings.CooldownHours
	}
	dateTime := time.Now()
	limit := minCandles * 60
	kline1Minute, err := a.ExchangeAPI.GetKLine(symbol, limit, binance.OneMinute)
	if err != nil {
		c <- a.apiFailResult(symbol)
		return
	}
	prices1Hour, err := binance.Get1HourPrices(kline1Minute)
	if err != nil {
		c <- a.apiFailResult(symbol)
		return
	}
	percent1Hour := []float64{}
	for i := 1; i < minCandles+1; i++ {
		end := i*60 - 1
		start := end - 59
		if end > len(prices1Hour)-1 {
			c <- a.apiFailResult(symbol)
			return
		}
		percent := ((prices1Hour[end] - prices1Hour[start]) * 100) / prices1Hour[end]
		percent1Hour = append(percent1Hour, percent)
	}
	current4HoursPercent := ((prices1Hour[239] - prices1Hour[0]) * 100) / prices1Hour[239]
	current24HoursPercent := binance.CalculateCurrent24HourPercent(prices24Hours, symbol.Name)

	kline24Hours, err := a.ExchangeAPI.GetKLine(symbol, 1500, binance.OneDay)
	if err != nil {
		c <- a.apiFailResult(symbol)
		return
	}
	age := len(kline24Hours)
	limit2 := math.Round((float64(age) / 30) + 1)

	kline1Month, err := a.ExchangeAPI.GetKLine(symbol, int(limit2), binance.OneMonth)
	if err != nil {
		c <- a.apiFailResult(symbol)
		return
	}
	ath := binance.GetMaximumAllTimeHigh(kline1Month)
	currentPercentageATH := ((ath - prices1Hour[len(prices1Hour)-1]) * 100 / ath)

	c <- a.calculateSymbolResults(symbol, percent1Hour, current4HoursPercent, current24HoursPercent, currentPercentageATH, age, dateTime)
}

func (a *AutoCoins) calculateSymbolResults(symbol binance.Symbol, percent1Hour []float64, current4HoursPercent float64, current24HoursPercent float64, currentPercentageATH float64, age int, dateTime time.Time) SymbolDataObject {
	// 1 hour percent
	x := a.Settings.CooldownHours - 1
	values1HrPercent := percent1Hour[:x]
	max := 0.0
	for _, val := range values1HrPercent {
		v := math.Abs(val)
		if v > max {
			max = v
		}
	}
	result1HrPercent := max < float64(a.Settings.Max1hrPercent)

	// 4 hour percent
	result4HrPercent := math.Abs(current4HoursPercent) < float64(a.Settings.Max4hrPercent)
	values4HrPercent := current4HoursPercent

	// 24 hour percent
	result24HrPercent := math.Abs(current24HoursPercent) < float64(a.Settings.Max24hrPercent)
	values24HrPercent := current24HoursPercent

	resultAth := currentPercentageATH > float64(a.Settings.MinAthPercent)
	valuesAth := currentPercentageATH

	resultAge := age > a.Settings.MinAge
	valueAge := age

	return SymbolDataObject{
		Symbol:             symbol,
		Percent1Hour:       result1HrPercent,
		Percent1HourValue:  values1HrPercent,
		Percent4Hour:       result4HrPercent,
		Percent4HourValue:  values4HrPercent,
		Percent24Hour:      result24HrPercent,
		Percent24HourValue: values24HrPercent,
		AllTimeHigh:        resultAth,
		AllTimeHighValue:   valuesAth,
		Age:                resultAge,
		AgeValue:           valueAge,
		Open:               false,
		Time:               dateTime,
		APIFailed:          false,
	}
}

func (a *AutoCoins) apiFailResult(symbol binance.Symbol) SymbolDataObject {
	return SymbolDataObject{
		Symbol:    symbol,
		APIFailed: true,
	}
}

type SymbolDataObject struct {
	Symbol             binance.Symbol `json:"symbol"`
	Percent1Hour       bool           `json:"perc1hr"`
	Percent1HourValue  []float64      `json:"perc1hrVal"`
	Percent4Hour       bool           `json:"perc4hr"`
	Percent4HourValue  float64        `json:"perc4hrVal"`
	Percent24Hour      bool           `json:"perc24hr"`
	Percent24HourValue float64        `json:"perc24hrVal"`
	AllTimeHigh        bool           `json:"Ath"`
	AllTimeHighValue   float64        `json:"AthVal"`
	Age                bool           `json:"Age"`
	AgeValue           int            `json:"AgeVal"`
	Open               bool           `json:"Open"`
	Time               time.Time      `json:"dateTime"`
	APIFailed          bool           `json:"apiFailed"`
}

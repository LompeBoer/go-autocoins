package autocoins

import (
	"fmt"
	"math"
)

type MarketSwing struct {
	Timeframe  string
	SwingMood  string
	CountTotal int
	Swing      float64

	Positive MarketSwingValues
	Negative MarketSwingValues
}

type MarketSwingValues struct {
	Percent    float64
	CoinCount  int
	CountTotal float64
	Average    float64
	Max        float64
	MaxCoin    string
}

func CalculateMarketSwing(objects []SymbolDataObject) []MarketSwing {
	marketSwing1 := MarketSwing{Timeframe: "1hr"}
	marketSwing4 := MarketSwing{Timeframe: "4hrs"}
	marketSwing24 := MarketSwing{Timeframe: "24hrs"}

	for _, object := range objects {
		if object.APIFailed {
			continue
		}
		marketSwing1.processObject(object.Values.Percent1Hour[0], object.Symbol.Name)
		marketSwing4.processObject(object.Values.Percent4Hour, object.Symbol.Name)
		marketSwing24.processObject(object.Values.Percent24Hour, object.Symbol.Name)
	}

	marketSwing1.calculate()
	marketSwing4.calculate()
	marketSwing24.calculate()

	return []MarketSwing{marketSwing1, marketSwing4, marketSwing24}
}

func (m *MarketSwing) processObject(percentValue float64, symbol string) {
	if percentValue < 0 {
		m.Negative.CoinCount++
		m.Negative.CountTotal += percentValue
	} else {
		m.Positive.CoinCount++
		m.Positive.CountTotal += percentValue
	}

	if percentValue > m.Positive.Max {
		m.Positive.Max = percentValue
		m.Positive.MaxCoin = symbol
	}
	if percentValue < m.Negative.Max {
		m.Negative.Max = percentValue
		m.Negative.MaxCoin = symbol
	}
}

func (m *MarketSwing) calculate() {
	if m.Positive.CoinCount != 0 {
		m.Positive.Average = m.Positive.CountTotal / float64(m.Positive.CoinCount)
	}
	if m.Negative.CoinCount != 0 {
		m.Negative.Average = m.Negative.CountTotal / float64(m.Negative.CoinCount)
	}
	m.CountTotal = m.Positive.CoinCount + m.Negative.CoinCount
	if m.CountTotal != 0 {
		m.Positive.Percent = (float64(m.Positive.CoinCount) / float64(m.CountTotal)) * 100
	}
	m.Negative.Percent = 100 - m.Positive.Percent

	m.Swing = m.Positive.Percent - m.Negative.Percent
	if m.Swing < 0 {
		m.SwingMood = fmt.Sprintf("%.0f%% Bearish", math.Abs(m.Swing))
	} else {
		m.SwingMood = fmt.Sprintf("%.0f%% Bullish", m.Swing)
	}
}

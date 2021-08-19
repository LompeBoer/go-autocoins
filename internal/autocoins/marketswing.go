package autocoins

import (
	"fmt"
	"io"
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

func (a *AutoCoins) marketSwingValues(objects []SymbolDataObject) []MarketSwing {
	marketSwing1 := MarketSwing{Timeframe: "1hr"}
	marketSwing4 := MarketSwing{Timeframe: "4hrs"}
	marketSwing24 := MarketSwing{Timeframe: "24hrs"}

	for _, object := range objects {
		if object.APIFailed {
			continue
		}
		marketSwing1.ProcessObject(object.Percent1HourValue[0], object.Symbol.Name)
		marketSwing4.ProcessObject(object.Percent4HourValue, object.Symbol.Name)
		marketSwing24.ProcessObject(object.Percent24HourValue, object.Symbol.Name)
	}

	marketSwing1.Calculate()
	marketSwing4.Calculate()
	marketSwing24.Calculate()

	return []MarketSwing{marketSwing1, marketSwing4, marketSwing24}
}

func (m *MarketSwing) ProcessObject(percentValue float64, symbol string) {
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

func (m *MarketSwing) Calculate() {
	m.Positive.Average = m.Positive.CountTotal / float64(m.Positive.CoinCount)
	m.Negative.Average = m.Negative.CountTotal / float64(m.Negative.CoinCount)
	m.CountTotal = m.Positive.CoinCount + m.Negative.CoinCount
	m.Positive.Percent = (float64(m.Positive.CoinCount) / float64(m.CountTotal)) * 100
	m.Negative.Percent = 100 - m.Positive.Percent

	m.Swing = m.Positive.Percent - m.Negative.Percent
	if m.Swing < 0 {
		m.Swing = math.Abs(m.Swing)
		m.SwingMood = fmt.Sprintf("%.0f%% Bearish", m.Swing)
	} else {
		m.SwingMood = fmt.Sprintf("%.0f%% Bullish", m.Swing)
	}
}

func (m *MarketSwing) WriteString(w io.Writer, applyStyle bool) {
	bold := ""
	if applyStyle {
		bold = "**"
	}
	fmt.Fprintf(w, "%sMarketSwing - Last %s%s - %s\n", bold, m.Timeframe, bold, m.SwingMood)
	fmt.Fprintf(w, "| %.0f%% Long | %d Coins | Avg %.2f%% | Max %.2f%% %s\n", m.Positive.Percent, m.Positive.CoinCount, m.Positive.Average, m.Positive.Max, m.Positive.MaxCoin)
	fmt.Fprintf(w, "| %.0f%% Short | %d Coins | Avg %.2f%% | Max %.2f%% %s\n", m.Negative.Percent, m.Negative.CoinCount, m.Negative.Average, m.Negative.Max, m.Negative.MaxCoin)
}

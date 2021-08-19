package binance

import (
	"log"
	"strconv"
)

func Get1HourPrices(klines []KLine) ([]float64, error) {
	values := []float64{}
	for _, kline := range klines {
		fv, err := strconv.ParseFloat(kline.Open, 64)
		if err != nil {
			return nil, err
		}
		values = append(values, fv)
	}

	return values, nil
}

func CalculateCurrent24HourPercent(prices24Hours []Ticker, symbolName string) float64 {
	for _, t := range prices24Hours {
		if t.Symbol == symbolName {
			v, err := strconv.ParseFloat(t.PriceChangePercent, 64)
			if err != nil {
				log.Printf("Current 24hr percent parse error: %s\n", err.Error())
				continue
			}
			return v
		}
	}
	return 0
}

func GetMaximumAllTimeHigh(klines []KLine) float64 {
	ath := 0.0
	for _, kline := range klines {
		fv, err := strconv.ParseFloat(kline.High, 64)
		if err != nil {
			log.Printf("Max ATH kline parse error: %s\n", err.Error())
			continue
		}
		if fv > ath {
			ath = fv
		}
	}

	return ath
}

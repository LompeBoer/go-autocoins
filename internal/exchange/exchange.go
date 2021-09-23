package exchange

type Symbol struct {
	Name       string
	BaseAsset  string
	QuoteAsset string
}

type Ticker struct {
	PriceChangePercent24h float64
}

type Kline struct {
}

type ExchangeService interface {
	GetSymbols() ([]Symbol, error)
	GetTickers([]string) ([]Ticker, error)
	GetKline(string, string, int, int) ([]Kline, error)
}

package bytbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

type API struct {
	BaseURL string
	client  http.Client
	context context.Context
	cancel  context.CancelFunc
}

type APIParams struct {
	BaseURL       string // BaseURL the base url for the ByBit API.
	ProxyURL      string
	ProxyUser     string
	ProxyPassword string
}

// https://api.bybit.com
func NewAPI(params APIParams) *API {
	api := API{
		BaseURL: params.BaseURL,
	}
	client := http.Client{
		Timeout: time.Second * 10,
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100
	transport.MaxConnsPerHost = 100
	transport.MaxIdleConnsPerHost = 100
	transport.DialContext = (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext
	transport.TLSHandshakeTimeout = 5 * time.Second

	if params.ProxyURL != "" {
		proxyURL, err := url.Parse(params.ProxyURL)
		if err != nil {
			log.Fatalf("Invalid proxy url: %s\n", err.Error())
		}
		if params.ProxyUser != "" {
			proxyURL.User = url.UserPassword(params.ProxyUser, params.ProxyPassword)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	client.Transport = transport

	ctx, cancel := context.WithCancel(context.Background())
	api.context = ctx
	api.cancel = cancel
	api.client = client

	return &api
}

func (a *API) Cancel() {
	a.cancel()
}

type ByBitResponse struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	ExtInfo string `json:"ext_info"`
	TimeNow string `json:"time_now"`
}

type Symbol struct {
	Name          string `json:"name"`
	Alias         string `json:"alias"`
	Status        string `json:"status"`
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
	PriceScale    int    `json:"price_scale"`
	TakerFee      string `json:"taker_fee"`
	MakerFee      string `json:"maker_fee"`
}

type SymbolsResponse struct {
	ByBitResponse
	Result []Symbol `json:"result"`
}

func (a *API) GetSymbols() ([]Symbol, error) {
	url := a.BaseURL + "/v2/public/symbols"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.client.Do(req.WithContext(a.context))
	if err != nil {
		return nil, err
	}

	var symbols SymbolsResponse
	err = json.NewDecoder(resp.Body).Decode(&symbols)
	if err != nil {
		return nil, err
	}

	return symbols.Result, nil
}

type Kline struct {
	ID       int     `json:"id"`
	Symbol   string  `json:"symbol"`
	Period   string  `json:"period"`
	StartAt  int64   `json:"start_at"`
	Volume   float64 `json:"volume"`
	Open     float64 `json:"open"`
	High     float64 `json:"high"`
	Low      float64 `json:"low"`
	Close    float64 `json:"close"`
	Interval int     `json:"interval"`
	OpenTime int64   `json:"open_time"`
	Turnover float64 `json:"turnover"`
}

type KlineResponse struct {
	ByBitResponse
	Result []Kline
}

func (a *API) GetKline(symbol, interval string, from, limit int) ([]Kline, error) {
	url := fmt.Sprintf("%s/public/linear/kline?symbol=%s&interval=%sfrom=%dlimit=%d", a.BaseURL, symbol, interval, from, limit)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.client.Do(req.WithContext(a.context))
	if err != nil {
		return nil, err
	}

	var klines KlineResponse
	err = json.NewDecoder(resp.Body).Decode(&klines)
	if err != nil {
		return nil, err
	}

	return klines.Result, nil
}

type Ticker struct {
	Symbol                 string `json:"symbol"`
	BidPrice               string `json:"bid_price"`
	AskPrice               string `json:"ask_price"`
	LastPrice              string `json:"last_price"`
	LastTickDirection      string `json:"last_tick_direction"`
	PrevPrice24h           string `json:"prev_price_24h"`
	Price24hPcnt           string `json:"price_24h_pcnt"`
	HighPrice24h           string `json:"high_price_24h"`
	LowPrice24h            string `json:"low_price_24h"`
	PrevPrice1h            string `json:"prev_price_1h"`
	Price1hPcnt            string `json:"price_1h_pcnt"`
	MarkPrice              string `json:"mark_price"`
	IndexPrice             string `json:"index_price"`
	OpenInterest           int64  `json:"open_interest"`
	OpenValue              string `json:"open_value"`
	TotalTurnover          string `json:"total_turnover"`
	Turnover24h            string `json:"turnover_24h"`
	TotalVolume            int64  `json:"total_volume"`
	Volume24h              int64  `json:"volume_24h"`
	FundingRate            string `json:"funding_rate"`
	PredictedFundingRate   string `json:"predicted_funding_rate"`
	NextFundingTime        string `json:"next_funding_time"`
	CountdownHour          string `json:"countdown_hour"`
	DeliveryFeeRate        string `json:"delivery_fee_rate"`
	PredictedDeliveryPrice string `json:"predicted_delivery_price"`
	DeliveryTime           string `json:"delivery_time"`
}

type TickerResponse struct {
	ByBitResponse
	Result []Ticker `json:"result"`
}

func (a *API) GetTicker(symbol string) ([]Ticker, error) {
	url := fmt.Sprintf("%s/v2/public/tickers?symbol=%s", a.BaseURL, symbol)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.client.Do(req.WithContext(a.context))
	if err != nil {
		return nil, err
	}

	var ticker TickerResponse
	err = json.NewDecoder(resp.Body).Decode(&ticker)
	if err != nil {
		return nil, err
	}

	return ticker.Result, nil
}

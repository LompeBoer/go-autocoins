package binance

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type API struct {
	DebugSaveResponses   bool      // DebugSaveResponses saves API responses to disk in the ./data dir.
	DebugReadResponses   bool      // DebugReadResponses read API responses from disk.
	BaseURL              string    // BaseURL the base url for the Binance API.
	UsedWeight           int       // UsedWeight the last value from the `X-Mbx-Used-Weight-1m` header.
	WeightLimit          int       // WeightLimit the maximum weight allowed to be used.
	LastWeightUpdate     time.Time // LastWeightUpdate time when `X-Mbx-Used-Weight-1m` header was last read.
	EstimatedWeightUsage int       // EstimatedWeightUsage when this value is exceeded throttle the requests.
	WeightWarning        bool      // WeightWarning indicates to pause requests until the warning is over.
	client               http.Client
	context              context.Context
	cancel               context.CancelFunc
}

type APIParams struct {
	DebugSaveResponses bool   // DebugSaveResponses saves API responses to disk in the ./data dir.
	DebugReadResponses bool   // DebugReadResponses read API responses from disk.
	BaseURL            string // BaseURL the base url for the Binance API.
	ProxyURL           string
	ProxyUser          string
	ProxyPassword      string
}

func NewAPI(params APIParams) *API {
	api := API{
		DebugSaveResponses: params.DebugSaveResponses,
		DebugReadResponses: params.DebugReadResponses,
		BaseURL:            params.BaseURL,
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

type BySymbolName []Symbol

func (a BySymbolName) Len() int {
	return len(a)
}

func (a BySymbolName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

func (a BySymbolName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type ExchangeInfo struct {
	Timezone    string      `json:"timezone"`
	ServerTime  int64       `json:"serverTime"`
	FuturesType string      `json:"futuresType"`
	RateLimits  []RateLimit `json:"rateLimits"`
	Symbols     []Symbol    `json:"symbols"`
}

type Symbol struct {
	Name         string `json:"symbol"`
	ContractType string `json:"contractType"`
	BaseAsset    string `json:"baseAsset"`
	QuoteAsset   string `json:"quoteAsset"`
	MarginAsset  string `json:"marginAsset"`
	OnboardDate  int64  `json:"onboardDate"`
}

type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

type Ticker struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	LastPrice          string `json:"lastPrice"`
	LastQty            string `json:"lastQty"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	FirstID            int64  `json:"firstId"`
	LastID             int64  `json:"lastId"`
	Count              int64  `json:"count"`
}

type KLine struct {
	OpenTime                 int64
	Open                     string
	High                     string
	Low                      string
	Close                    string
	Volume                   string
	CloseTime                int64
	QuoteAssetVolume         string
	NumberOfTrades           int64
	TakerBuyBaseAssetVolume  string
	TakerBuyQuoteAssetVolume string
	Ignore                   string
}

func (a *API) GetExchangeInfo() (ExchangeInfo, error) {
	var exchangeInfo ExchangeInfo

	url := a.BaseURL + "/fapi/v1/exchangeInfo"
	r, err := a.requestGet(url, true)
	if err != nil {
		return exchangeInfo, err
	}

	data := a.handleResponse(url, r.Body)

	if err := json.Unmarshal(data, &exchangeInfo); err != nil {
		return exchangeInfo, err
	}

	for _, limit := range exchangeInfo.RateLimits {
		if limit.RateLimitType == "REQUEST_WEIGHT" {
			a.WeightLimit = limit.Limit
		}
	}

	return exchangeInfo, nil
}

// GetTicker get the 24hr ticker data
// https://binance-docs.github.io/apidocs/futures/en/#24hr-ticker-price-change-statistics
func (a *API) GetTicker() ([]Ticker, error) {
	url := a.BaseURL + "/fapi/v1/ticker/24hr"
	r, err := a.requestGet(url, false)
	if err != nil {
		return nil, err
	}

	data := a.handleResponse(url, r.Body)

	var binanceTicker []Ticker
	if err := json.Unmarshal(data, &binanceTicker); err != nil {
		return nil, err
	}

	return binanceTicker, nil
}

type KlineInterval string

const (
	OneMinute      KlineInterval = "1m"
	ThreeMinutes   KlineInterval = "3m"
	FiveMinutes    KlineInterval = "5m"
	FifteenMinutes KlineInterval = "15m"
	ThirtyMinute   KlineInterval = "30m"
	OneHour        KlineInterval = "1h"
	TwoHours       KlineInterval = "2h"
	FourHours      KlineInterval = "4h"
	SixHours       KlineInterval = "6h"
	EightHours     KlineInterval = "8h"
	TwelveHours    KlineInterval = "12h"
	OneDay         KlineInterval = "1d"
	ThreeDays      KlineInterval = "3d"
	OneWeek        KlineInterval = "1w"
	OneMonth       KlineInterval = "1M"
)

// GetKLine return the candlestick data
// https://binance-docs.github.io/apidocs/futures/en/#kline-candlestick-data
func (a *API) GetKLine(symbol Symbol, limit int, interval KlineInterval) ([]KLine, error) {
	l := strconv.Itoa(limit)
	url := fmt.Sprintf("%s/fapi/v1/klines?symbol=%s&interval=%s&limit=%s", a.BaseURL, symbol.Name, interval, l)
	r, err := a.requestGet(url, false)
	if err != nil {
		log.Printf("ERROR: GetKLine:requestGet: %s\n", err.Error())
		return nil, err
	}

	responseData := a.handleResponse(url, r.Body)

	var data [][]interface{}
	if err := json.Unmarshal(responseData, &data); err != nil {
		log.Printf("ERROR: GetKLine:Unmarshal: %s\n", err.Error())
		return nil, err
	}
	klines := []KLine{}
	for _, d := range data {
		kline := KLine{}
		for i, v := range d {
			switch i {
			case 0:
				kline.OpenTime = a.klineInt64Value(v)
			case 1:
				kline.Open = a.klineStringValue(v)
			case 2:
				kline.High = a.klineStringValue(v)
			case 3:
				kline.Low = a.klineStringValue(v)
			case 4:
				kline.Close = a.klineStringValue(v)
			case 5:
				kline.Volume = a.klineStringValue(v)
			case 6:
				kline.CloseTime = a.klineInt64Value(v)
			case 7:
				kline.QuoteAssetVolume = a.klineStringValue(v)
			case 8:
				kline.NumberOfTrades = a.klineInt64Value(v)
			case 9:
				kline.TakerBuyBaseAssetVolume = a.klineStringValue(v)
			case 10:
				kline.TakerBuyQuoteAssetVolume = a.klineStringValue(v)
			case 11:
				kline.Ignore = a.klineStringValue(v)
			}
		}
		klines = append(klines, kline)
	}

	return klines, nil
}

func (a *API) klineStringValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int64:
	case float64:
	default:
	}

	return ""
}

func (a *API) klineInt64Value(value interface{}) int64 {
	switch v := value.(type) {
	case int64:
		return v
	case float64:
		return int64(v)
	case int:
	case int32:
	case string:
	default:
	}

	return 0
}

// pauseRequest sleeps for specified time. Returns true when finished, false when cancelled.
func (a *API) pauseRequest(sleep time.Duration) bool {
	select {
	case <-a.context.Done():
		return false
	case <-time.After(sleep):
		return true
	}
}

func (a *API) requestGet(url string, skipWeightCheck bool) (*http.Response, error) {
	if !a.DebugReadResponses {
		// TODO: do reset time based on exchangeInfo api call.
		if time.Since(a.LastWeightUpdate) > time.Minute {
			a.UsedWeight = 0
		}
		if a.WeightLimit > 0 && !skipWeightCheck {
			limit := a.CheckForWeightLimit()
			if float64(a.UsedWeight) > limit {
				// TODO: not print warning for each request, too much spam.
				// TODO: put in a while?
				log.Printf("WARNING: exceeding weight limit: %d/%d\n", a.UsedWeight, a.EstimatedWeightUsage)
				if !a.pauseRequest(60 * time.Second) {
					return nil, errors.New("context cancelled")
				}
			}
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		resp, err := a.client.Do(req.WithContext(a.context))
		if err != nil {
			return resp, err
		}

		if weightHeader, ok := resp.Header["X-Mbx-Used-Weight-1m"]; ok {
			if len(weightHeader) > 0 {
				weightUsed, err := strconv.Atoi(weightHeader[0])
				if err != nil {
					return resp, err
				}
				a.UsedWeight = weightUsed
				a.LastWeightUpdate = time.Now()
			}
		}

		return resp, err
	}

	// This is for debugging purposes.
	// When `DebugReadResponses` is true this will read a saved HTTP response from disk.
	filename := FilenameForURL(url)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		Body: file,
	}, nil
}

func (a *API) handleResponse(url string, body io.ReadCloser) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)

	data := buf.Bytes()

	if !a.DebugSaveResponses {
		return data
	}

	// This is for debugging purposes.
	// When `DebugSaveResponses` is true this will save a copy of the HTTP response to disk.
	filename := FilenameForURL(url)
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

// filenameForURL hashes the given URL (using sha1) to be used as filename.
func FilenameForURL(url string) string {
	h := sha1.New()
	io.WriteString(h, url)
	hash := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("data/%s", hash)
}

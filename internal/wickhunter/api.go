package wickhunter

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type API struct {
	BaseURL string
	client  http.Client
	context context.Context
	cancel  context.CancelFunc
}

func NewAPI(apiBaseURL string) *API {
	api := API{
		BaseURL: apiBaseURL,
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
	client.Transport = transport

	ctx, cancel := context.WithCancel(context.Background())
	api.context = ctx
	api.cancel = cancel
	api.client = client

	return &api
}

func (a *API) SetSymbolTrading(symbol string, enabled bool) error {
	url := fmt.Sprintf("%s/symbols/%s/enable/%s", a.BaseURL, strings.ToLower(symbol), strconv.FormatBool(enabled))
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}
	resp, err := a.client.Do(req.WithContext(a.context))
	if err != nil {
		return err
	}

	// fmt.Printf("%s (%d)\n", url, resp.StatusCode)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("response status code is '%d' (%s)", resp.StatusCode, resp.Status)
	}

	return nil
}

type Position struct {
	Symbol    string `json:"symbol"`
	Permitted bool   `json:"permitted"`
	State     string `json:"state"`
}

func (a *API) GetPositions() ([]Position, error) {
	url := fmt.Sprintf("%s/bot/positions", a.BaseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.client.Do(req.WithContext(a.context))
	if err != nil {
		return nil, err
	}

	// fmt.Printf("%s (%d)\n", url, resp.StatusCode)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("response status code is '%d' (%s)", resp.StatusCode, resp.Status)
	}

	var positions []Position
	if err := json.NewDecoder(resp.Body).Decode(&positions); err != nil {
		return nil, err
	}

	return positions, nil
}

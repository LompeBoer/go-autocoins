package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

type API struct {
	Owner   string
	Repo    string
	client  http.Client
	context context.Context
	cancel  context.CancelFunc
}

func NewAPI(owner string, repo string) *API {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
	}).DialContext
	transport.TLSHandshakeTimeout = 5 * time.Second

	client.Transport = transport

	ctx, cancel := context.WithCancel(context.Background())

	return &API{
		Owner:   owner,
		Repo:    repo,
		client:  client,
		context: ctx,
		cancel:  cancel,
	}
}

type Release struct {
	HTMLURL string `json:"html_url"`
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
}

func (a *API) LatestRelease() (Release, error) {
	var release Release

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", a.Owner, a.Repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return release, err
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	resp, err := a.client.Do(req.WithContext(a.context))
	if err != nil {
		return release, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return release, err
	}

	return release, err
}

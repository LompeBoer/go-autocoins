package discord

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type DiscordWebHook struct {
	URL     string
	Enabled bool
}

type SendParams struct {
	Content string `json:"content"`
}

func (w *DiscordWebHook) SendMessage(message string) error {
	if !w.Enabled || w.URL == "" {
		return nil
	}

	params := SendParams{
		Content: message,
	}
	b, err := json.Marshal(params)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(b)
	_, err = http.Post(w.URL, "application/json", reader)
	if err != nil {
		return err
	}
	// fmt.Printf("response post status code: %d\n", r.StatusCode)
	return nil
}

func (w *DiscordWebHook) SendError(message string, mention bool) error {
	if !w.Enabled || w.URL == "" {
		return nil
	}

	if mention {
		message = "@here " + message
	}

	return w.SendMessage(message)
}

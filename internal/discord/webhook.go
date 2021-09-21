package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DiscordWebHook struct {
	URL     string
	Enabled bool
}

func (w *DiscordWebHook) SendTextMessage(message string) error {
	return w.SendMessage(DiscordWebhookMessage{
		Content: message,
	})
}

func (w *DiscordWebHook) SendMessage(message DiscordWebhookMessage) error {
	if !w.Enabled || w.URL == "" {
		return nil
	}

	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(b)
	resp, err := http.Post(w.URL, "application/json", reader)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("failed to post Discord webhook: %s", resp.Status)
	}

	return nil
}

func (w *DiscordWebHook) SendError(message string, mention bool) error {
	if !w.Enabled || w.URL == "" {
		return nil
	}

	if mention {
		message = "@here " + message
	}

	return w.SendTextMessage(message)
}

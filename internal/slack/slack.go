package slack

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

// Slack represents the metadata that is required for the posting to slack webhook
type Slack struct {
	Channel   string `json:"channel"`
	IconEmoji string `json:"icon_emoji"`
	Username  string `json:"username"`
	Text      string `json:"text"`
}

// SlackWebhook represents the config for the slack api webhook
type SlackWebhook struct {
	Data       Slack
	WebhookURL string
}

// Send will post a message to the targetted channel
func (s SlackWebhook) Send(msg string) error {
	s.Data.Text = msg
	b, err := json.Marshal(s.Data)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("payload", string(b))

	req, err := http.NewRequest("POST", s.WebhookURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	log.Println(resp.Status)
	return nil
}

package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"lambda/internal/weather"
)

// Discordに送信するpayloadの構造体
type DiscordWebhookPayload struct {
	Content string  `json:"content,omitempty"`
	Embeds  []Embed `json:"embeds,omitempty"`
}

// Discordの埋め込みメッセージの構造体
type Embed struct {
	Title  string       `json:"title"`
	Color  int          `json:"color"`
	URL    string       `json:"url,omitempty"`
	Fields []EmbedField `json:"fields,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

func SendDiscordMessage(webhookURL, city string, weatherInfo *weather.WeatherResponse) error {
	forecast := weatherInfo.Forecasts[0]
	embed := Embed{
		Title: fmt.Sprintf("今日 %s は雨が降るらしいよ！", city),
		Color: 0x00FF00,
		URL:   weatherInfo.Link,
		Fields: []EmbedField{
			{
				Name:  weatherInfo.Title,
				Value: forecast.Telop,
			},
		},
	}
	payload := DiscordWebhookPayload{
		Embeds: []Embed{embed},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("discord webhookエラー: ステータス %d, レスポンス: %s", resp.StatusCode, string(body))
	}
	return nil
}

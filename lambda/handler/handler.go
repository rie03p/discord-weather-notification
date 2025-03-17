package handler

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"lambda/internal/discord"
	"lambda/internal/weather"

	"github.com/aws/aws-lambda-go/lambda"
)

func SendRainAlert(discordWebhookURL string) error {
	for city := range weather.CityIDMap {
		weatherInfo, err := weather.GetWeather(city)
		if err != nil {
			log.Printf("都市 %s の天気取得エラー: %v\n", city, err)
			continue
		}
		if len(weatherInfo.Forecasts) == 0 {
			log.Printf("都市 %s の予報データがありません\n", city)
			continue
		}
		todayTelop := weatherInfo.Forecasts[0].Telop
		if strings.Contains(todayTelop, "雨") {
			if err := discord.SendDiscordMessage(discordWebhookURL, city, weatherInfo); err != nil {
				log.Printf("都市 %s のDiscord送信エラー: %v\n", city, err)
			}
		}
	}
	return nil
}

func Handler(ctx context.Context) error {
	discordWebhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if discordWebhookURL == "" {
		return errors.New("DISCORD_WEBHOOK_URL 環境変数が設定されていません")
	}
	return SendRainAlert(discordWebhookURL)
}

func StartLambda() {
	lambda.Start(Handler)
}

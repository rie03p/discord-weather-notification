package handler

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"sync"

	"lambda/internal/discord"
	"lambda/internal/weather"

	"github.com/aws/aws-lambda-go/lambda"
)

func SendRainAlert(discordWebhookURL string) error {
	var wg sync.WaitGroup

	for city := range weather.CityIDMap {
		wg.Add(1)
		go func(city string) {
			defer wg.Done()

			weatherInfo, err := weather.GetWeather(city)
			if err != nil {
				log.Printf("都市 %s の天気取得エラー: %v\n", city, err)
				return
			}
			if len(weatherInfo.Forecasts) == 0 {
				log.Printf("都市 %s の予報データがありません\n", city)
				return
			}
			todayTelop := weatherInfo.Forecasts[0].Telop
			if strings.Contains(todayTelop, "雨") {
				if err := discord.SendDiscordMessage(discordWebhookURL, city, weatherInfo); err != nil {
					log.Printf("都市 %s のDiscord送信エラー: %v\n", city, err)
				}
			}
		}(city)
	}

	wg.Wait()

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

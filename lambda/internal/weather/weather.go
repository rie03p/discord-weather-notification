package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// 天気予報APIのレスポンス構造体
type WeatherResponse struct {
	Title     string     `json:"title"`
	Forecasts []Forecast `json:"forecasts"`
	Link      string     `json:"link"`
}

// 各日の天気情報
type Forecast struct {
	Telop string        `json:"telop"`
	Image ForecastImage `json:"image"`
}

type ForecastImage struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// 天気予報APIのエンドポイント
var WeatherAPIURL = "https://weather.tsukumijima.net/api/forecast"

// 都市名とAPIに必要な都市IDの対応表
var CityIDMap = map[string]string{
	"東京":  "130010",
	"神奈川": "140010",
	"滋賀":  "250010",
}

func GetWeather(city string) (*WeatherResponse, error) {
	cityID, ok := CityIDMap[city]
	if !ok {
		return nil, errors.New("city not found")
	}
	url := fmt.Sprintf("%s?city=%s", WeatherAPIURL, cityID)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s %d", resp.Status, resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var weatherResp WeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return nil, err
	}
	return &weatherResp, nil
}

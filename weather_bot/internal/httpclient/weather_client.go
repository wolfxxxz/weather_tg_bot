package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"weather_bot/internal/apperrors"
	"weather_bot/internal/config"

	"github.com/sirupsen/logrus"
)

type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

type Clouds struct {
	All int `json:"all"`
}

type Sys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int64  `json:"sunrise"`
	Sunset  int64  `json:"sunset"`
}

type GetWeatherResponse struct {
	Coord      *Coord     `json:"coord"`
	Weather    []*Weather `json:"weather"`
	Base       string     `json:"base"`
	Main       *Main      `json:"main"`
	Visibility int        `json:"visibility"`
	Wind       *Wind      `json:"wind"`
	Clouds     *Clouds    `json:"clouds"`
	Dt         int64      `json:"dt"`
	Sys        *Sys       `json:"sys"`
	Timezone   int        `json:"timezone"`
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Cod        int        `json:"cod"`
}

type WeatherClient struct {
	config *config.Config
	client *http.Client
	log    *logrus.Logger
}

func NewWeatherClient(config *config.Config, httpClient *http.Client, log *logrus.Logger) *WeatherClient {
	return &WeatherClient{
		config: config,
		client: httpClient,
		log:    log,
	}
}

func (wcl *WeatherClient) GetWeatherForecast(latitude float64, longitude float64) (*GetWeatherResponse, error) {
	weatherApiURL := wcl.appendQueryParamsToWeatherApiHost(latitude, longitude)
	weather, err := wcl.getWeatherForecastResponse(weatherApiURL)
	if err != nil {
		wcl.log.Error(err)
		return nil, err
	}

	return weather, nil
}

func (wcl *WeatherClient) appendQueryParamsToWeatherApiHost(latitude float64, longitude float64) string {
	weatherApiURL := fmt.Sprintf("%vlat=%v&lon=%v&appid=%v&units=metric", wcl.config.Host, latitude, longitude, wcl.config.KeyApi)
	return weatherApiURL
}

func (wcl *WeatherClient) getWeatherForecastResponse(weatherApiURL string) (*GetWeatherResponse, error) {
	resp, err := wcl.client.Get(weatherApiURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			wcl.log.Error(err)
		}
		return nil, apperrors.WeatherClientError.AppendMessage(string(body))
	}

	var response GetWeatherResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

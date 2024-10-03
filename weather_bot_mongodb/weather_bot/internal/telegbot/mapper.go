package telegbot

import (
	"fmt"
	"weather_bot/internal/apperrors"
	"weather_bot/internal/httpclient"
)

func MapGetWeatherResponseToWeatherAnswer(data *httpclient.GetWeatherResponse) (string, error) {
	if data.Main == nil {
		return "", apperrors.BotMapperEncodingErr.AppendMessage("data.Main == nil")
	}

	temperature := fmt.Sprintf("Temperature Middle:%.2f°C Max:%.2f°C Min:%.2f°C",
		data.Main.Temp, data.Main.TempMax, data.Main.TempMin)
	weather := temperature
	return weather, nil

}

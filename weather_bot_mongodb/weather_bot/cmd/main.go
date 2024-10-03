package main

import (
	"context"
	"net/http"
	"weather_bot/internal/apperrors"
	"weather_bot/internal/config"
	"weather_bot/internal/database"
	"weather_bot/internal/httpclient"
	"weather_bot/internal/log"
	"weather_bot/internal/repositories"
	"weather_bot/internal/telegbot"
)

var sendWeatherForecast = 15

func main() {

	logger, err := log.NewLogAndSetLevel("info")
	if err != nil {
		logger.Debug(err)
	}

	conf := config.NewConfig()
	err = conf.ParseConfig(".env", logger)
	if err != nil {
		logger.Fatal(apperrors.EnvConfigLoadError.AppendMessage(err))
	}

	if err = log.SetLevel(logger, conf.LogLevel); err != nil {
		logger.Debug(err)
	}

	httpClient := &http.Client{}
	weatherClient := httpclient.NewWeatherClient(conf, httpClient, logger)

	ctx := context.Background()
	mongoDB, err := database.InitClient(ctx, conf, logger)
	if err != nil {
		logger.Debug(err)
	}

	userRepo := repositories.NewUserRepo(conf, logger, mongoDB)
	if err != nil {
		logger.Fatal(err)
	}

	bot, err := telegbot.NewBot(conf, httpClient, weatherClient, logger, userRepo)
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		if err := bot.SendingWeatherForecastNotifications(ctx, sendWeatherForecast, 1); err != nil {
			logger.Fatal(err)
		}
	}()

	logger.Info("BOT is replying on messages.")
	if err = bot.ReplyingOnMessages(ctx); err != nil {
		logger.Fatal(err)
	}

}

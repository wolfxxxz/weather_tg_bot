package telegbot

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"weather_bot/internal/apperrors"
	"weather_bot/internal/config"
	"weather_bot/internal/httpclient"
	"weather_bot/internal/models"
	"weather_bot/internal/repositories"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	botClient         *tgbotapi.BotAPI
	weatherClient     *httpclient.WeatherClient
	log               *logrus.Logger
	db                repositories.UserRepoInterface
	timeoutMongoQuery int
}

func NewBot(config *config.Config, httpClient *http.Client, wcl *httpclient.WeatherClient, log *logrus.Logger, db repositories.UserRepoInterface) (*Bot, error) {
	client, err := tgbotapi.NewBotAPIWithClient(config.Token, "https://api.telegram.org/bot%s/%s", httpClient)
	if err != nil {
		return nil, err
	}

	client.Token = config.Token
	log.Infof("Authorized on account %s", client.Self.UserName)
	timeoutMongoQuery, err := strconv.Atoi(config.TimeoutMongoQuery)
	if err != nil {
		return nil, err
	}
	return &Bot{botClient: client, log: log, weatherClient: wcl, db: db, timeoutMongoQuery: timeoutMongoQuery}, nil
}

func (bot *Bot) ReplyingOnMessages(ctx context.Context) error {
	updateConfig := tgbotapi.NewUpdate(allAvailableUpdates)
	updateConfig.Timeout = expectAnswerSec
	updates := bot.botClient.GetUpdatesChan(updateConfig)

	for update := range updates {
		answerMessage := bot.replyOnNewMessage(ctx, &update)
		if answerMessage == nil {
			bot.log.Debug("somebody sent geotranslation")
			continue
		}

		_, err := bot.botClient.Send(answerMessage)
		if err != nil {
			bot.log.Debugf("Error sending message: %v\n", err)
		}
	}

	return nil
}

func (bot *Bot) replyOnNewMessage(ctx context.Context, upd *tgbotapi.Update) *tgbotapi.MessageConfig {
	if upd.Message == nil {
		return nil
	}

	chatID := upd.Message.Chat.ID
	text := upd.Message.Text
	bot.log.Infof("Replying on message. Text: %s; ChatID %v\n", text, chatID)

	reply := "If you want to subscribe, Please send geolocation"
	if upd.Message.Location != nil {
		user := models.Create(chatID, upd.Message.Location.Latitude, upd.Message.Location.Longitude)

		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(bot.timeoutMongoQuery))
		defer cancel()

		err := bot.db.SaveUserIfNotExist(ctxWithTimeout, user)
		if err != nil {
			bot.log.Error(apperrors.BotSendMessageError.AppendMessage(err))
			reply = "Something went wrong, provider is quilty"
			msg := tgbotapi.NewMessage(chatID, reply)
			return &msg
		}

		reply = "you subscribe on the weather forecast"
	}

	message := tgbotapi.NewMessage(chatID, reply)
	return &message
}

func (bot *Bot) SendingWeatherForecastNotifications(ctx context.Context, hourRing, timeOutMinute int) error {
	ticker := time.NewTicker(time.Minute * time.Duration(timeOutMinute))
	defer ticker.Stop()

	for t := range ticker.C {
		if t.Hour() == hourRing {
			if err := bot.pushScheduledWeatherForecast(ctx, t); err != nil {
				bot.log.Debug(err)
				return err
			}
		}
	}
	return nil
}

func (bot *Bot) pushScheduledWeatherForecast(ctx context.Context, executedTime time.Time) error {
	users, err := bot.db.GetAllUsers(ctx)
	if err != nil {
		bot.log.Debug(err)
		return err
	}

	for _, user := range users {
		time.Sleep(time.Second * 1)
		reply := ""
		getWeatherForecastResponse, err := bot.weatherClient.GetWeatherForecast(user.Coordinates.Latitude, user.Coordinates.Longitude)
		if err != nil {
			bot.log.Error(apperrors.BotSendMessageError.AppendMessage(err))
			reply = "Something went wrong, provider is quilty"
			return err
		}

		bot.log.Infof("weather answer. GetWeatherForecastResponse: %+v", getWeatherForecastResponse)
		reply, err = MapGetWeatherResponseToWeatherAnswer(getWeatherForecastResponse)
		if err != nil {
			bot.log.Error(apperrors.BotSendMessageError.AppendMessage(err))
			reply = "Something went wrong, noone is to blame. Just a mistake"
			return err
		}

		message := tgbotapi.NewMessage(user.ChatID, reply)
		bot.log.Infof("Sending a reply to user. Reply: %+v; ChatID: %v;", reply, user.ChatID)
		_, err = bot.botClient.Send(message)
		if err != nil {
			bot.log.Debugf("Error sending message: %v\n", err)
			return err
		}

		user.Update()
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(bot.timeoutMongoQuery))
		defer cancel()
		if err = bot.db.UpdateModification(ctxWithTimeout, user); err != nil {
			bot.log.Debug(err)
			return err
		}
	}

	bot.log.Infof("Execute a scheduled task at %v sent users %v", executedTime.Format("15:04"), len(users))
	return nil
}

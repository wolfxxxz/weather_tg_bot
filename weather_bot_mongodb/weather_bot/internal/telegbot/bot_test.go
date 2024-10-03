package telegbot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
	"weather_bot/internal/config"
	"weather_bot/internal/httpclient"
	"weather_bot/internal/log"
	"weather_bot/mock"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const messageChatId = 1234
const testToken = "153667468:AAHlSHlMqSt1f_uFmVRJbm5gntu2HI4WW8I"

var logger, _ = log.NewLogAndSetLevel("debug")

type testTransport func(*http.Request) (*http.Response, error)

func (transport testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return transport(req)
}

func fakeHTTPBotClient(statusCode int, jsonResponse string) *http.Client {
	return &http.Client{
		Transport: testTransport(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: statusCode,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(jsonResponse)),
			}, nil
		}),
	}
}

func generateBotOkJsonApiResponse() (string, error) {
	testUser := tgbotapi.User{
		ID:        123,
		FirstName: "John",
		LastName:  "Doe",
	}

	user, err := json.Marshal(&testUser)
	if err != nil {
		return "", err
	}

	testApiResponse := tgbotapi.APIResponse{
		Ok:          true,
		Result:      user,
		ErrorCode:   0,
		Description: "",
		Parameters:  nil,
	}

	response, err := json.Marshal(&testApiResponse)
	if err != nil {
		return "", err
	}

	return string(response), nil
}

func generateBotFailJsonApiResponse(failErrorCode int, errorText string) (string, error) {
	testUser := tgbotapi.User{
		ID:        123,
		FirstName: "John",
		LastName:  "Doe",
	}

	user, err := json.Marshal(&testUser)
	if err != nil {
		return "", err
	}

	testApiResponse := tgbotapi.APIResponse{
		Ok:          false,
		Result:      user,
		ErrorCode:   failErrorCode,
		Description: errorText,
		Parameters:  nil,
	}

	response, err := json.Marshal(&testApiResponse)
	if err != nil {
		return "", err
	}

	return string(response), nil
}

func TestReplyOnNewMessageLocationErrMap(t *testing.T) {
	response, err := generateBotOkJsonApiResponse()
	if err != nil {
		t.Error(err)
		return
	}

	clientFake := fakeHTTPBotClient(200, response)
	msg := &tgbotapi.Location{
		Longitude: 30.627429,
		Latitude:  50.3906,
	}
	ttPass := []struct {
		name        string
		location    *tgbotapi.Location
		messageText string
		messageType string
		botReply    string
	}{
		{fmt.Sprintf("existing %s command", "loc"), msg, "msg", "bot_command", "you subscribe on the weather forecast"},
	}

	conf := config.Config{
		Token:             testToken,
		Host:              "https://api.openweathermap.org/data/2.5/weather?",
		KeyApi:            "sdfsfsfsfssfsdfdfsfsfsfergtrgt",
		TimeoutMongoQuery: "5",
	}

	weatherClient := httpclient.NewWeatherClient(&conf, clientFake, logger)
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUserRepo := mock.NewMockUserRepoInterface(ctl)
	mockUserRepo.EXPECT().SaveUserIfNotExist(gomock.Any(), gomock.Any()).Return(nil)
	bot, err := NewBot(&conf, clientFake, weatherClient, logger, mockUserRepo)
	if err != nil {
		t.Error(err)
		return
	}

	bot.botClient.Debug = true
	for _, tc := range ttPass {
		upd := tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: messageChatId,
				},
				Location: tc.location,
				Entities: []tgbotapi.MessageEntity{
					{
						Type:     tc.messageType,
						Offset:   0,
						Length:   0,
						URL:      "",
						User:     nil,
						Language: "",
					},
				},
			},
		}

		msg := bot.replyOnNewMessage(context.Background(), &upd)
		if msg != nil && msg.Text != tc.botReply {
			t.Errorf("bot reply should be %s, but got %s", tc.botReply, msg.Text)
		}
	}
}

func TestReplyOnNewMessageLocationIfNotExistErr(t *testing.T) {
	response, err := generateBotOkJsonApiResponse()
	if err != nil {
		t.Error(err)
		return
	}

	clientFake := fakeHTTPBotClient(200, response)
	msg := &tgbotapi.Location{
		Longitude: 30.627429,
		Latitude:  50.3906,
	}
	ttPass := []struct {
		name        string
		location    *tgbotapi.Location
		messageText string
		messageType string
		botReply    string
	}{
		{fmt.Sprintf("existing %s command", "loc"), msg, "msg", "bot_command", "Something went wrong, provider is quilty"},
	}

	conf := config.Config{
		Token:             testToken,
		Host:              "https://api.openweathermap.org/data/2.5/weather?",
		KeyApi:            "sdfsfsfsfssfsdfdfsfsfsfergtrgt",
		TimeoutMongoQuery: "5",
	}

	weatherClient := httpclient.NewWeatherClient(&conf, clientFake, logger)
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUserRepo := mock.NewMockUserRepoInterface(ctl)
	expectedError := errors.New("Some error")
	mockUserRepo.EXPECT().SaveUserIfNotExist(gomock.Any(), gomock.Any()).Return(expectedError)
	bot, err := NewBot(&conf, clientFake, weatherClient, logger, mockUserRepo)
	if err != nil {
		t.Error(err)
		return
	}

	bot.botClient.Debug = true
	for _, tc := range ttPass {
		upd := tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: messageChatId,
				},
				Location: tc.location,
				Entities: []tgbotapi.MessageEntity{
					{
						Type:     tc.messageType,
						Offset:   0,
						Length:   0,
						URL:      "",
						User:     nil,
						Language: "",
					},
				},
			},
		}

		msg := bot.replyOnNewMessage(context.TODO(), &upd)
		if msg != nil && msg.Text != tc.botReply {
			t.Errorf("bot reply should be %s, but got %s", tc.botReply, msg.Text)
		}
	}
}

func TestBotMessagingFailed(t *testing.T) {
	ttFail := []struct {
		name                   string
		statusCode             int
		apiResponseDescription string
	}{
		{"invalid token", http.StatusInternalServerError, "Test error"},
		{"init bot error", http.StatusInternalServerError, "Test error"},
	}

	for _, tc := range ttFail {
		response, err := generateBotFailJsonApiResponse(1, tc.apiResponseDescription)
		if err != nil {
			t.Error(err)
			return
		}

		conf := config.Config{
			Token: testToken,
		}

		clientFake := fakeHTTPBotClient(200, response)
		weatherClient := httpclient.NewWeatherClient(&conf, clientFake, logger)
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		mockUserRepo := mock.NewMockUserRepoInterface(ctl)
		_, err = NewBot(&conf, fakeHTTPBotClient(tc.statusCode, response), weatherClient, logger, mockUserRepo)
		require.Error(t, err)

		expectedErr := fmt.Errorf("%s", tc.apiResponseDescription)
		assert.EqualError(t, err, expectedErr.Error())
	}
}

func TestReplyOnTxtMessageSuccess(t *testing.T) {
	response, err := generateBotOkJsonApiResponse()
	if err != nil {
		t.Error(err)
		return
	}

	clientFake := fakeHTTPBotClient(200, response)
	ttPass := []struct {
		name        string
		messageText string
		messageType string
		botReply    string
	}{
		{fmt.Sprintf("existing %s command", "ua"), "Hello", "bot_command", "If you want to subscribe, Please send geolocation"},
	}

	conf := config.Config{
		Token:             testToken,
		Host:              "https://api.openweathermap.org/data/2.5/weather?",
		KeyApi:            "sdfsfsfsfssfsdfdfsfsfsfergtrgt",
		TimeoutMongoQuery: "5",
	}

	weatherClient := httpclient.NewWeatherClient(&conf, clientFake, logger)

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUserRepo := mock.NewMockUserRepoInterface(ctl)

	bot, err := NewBot(&conf, clientFake, weatherClient, logger, mockUserRepo)
	if err != nil {
		t.Error(err)
		return
	}

	bot.botClient.Debug = true

	for _, tc := range ttPass {
		upd := tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: messageChatId,
				},
				Text: tc.messageText,
				Entities: []tgbotapi.MessageEntity{
					{
						Type:     tc.messageType,
						Offset:   0,
						Length:   0,
						URL:      "",
						User:     nil,
						Language: "",
					},
				},
			},
		}

		msg := bot.replyOnNewMessage(context.TODO(), &upd)
		if msg != nil && msg.Text != tc.botReply {
			t.Errorf("bot reply should be %s, but got %s", tc.botReply, msg.Text)
		}
	}
}

func TestNewBot(t *testing.T) {
	apiToken := os.Getenv("TOKEN")
	testConfig := &config.Config{
		Token: apiToken,
	}

	response, err := generateBotOkJsonApiResponse()

	if err != nil {
		t.Error(err)
		return
	}

	testTgClientHttp := fakeHTTPBotClient(200, response)
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUserRepo := mock.NewMockUserRepoInterface(ctl)
	if _, err := NewBot(testConfig, testTgClientHttp, nil, logger, mockUserRepo); err != nil {
		t.Log(err)
	}
}

func TestReplyingOnMessages(t *testing.T) {
	authResponse, err := generateBotOkJsonApiResponse()
	if err != nil {
		t.Error(err)
		return
	}

	messageAPI := []*tgbotapi.Update{{
		UpdateID: 300,
		Message: &tgbotapi.Message{
			MessageID: 100,
			Chat: &tgbotapi.Chat{
				ID: 200,
			},
			Text: "",
			Location: &tgbotapi.Location{
				Longitude: 21.017532,
				Latitude:  52.237049,
			},
		},
	},
	}

	messageAPIResponse, err := json.Marshal(&messageAPI)
	if err != nil {
		t.Error(err)
		return
	}

	apiResponse := tgbotapi.APIResponse{
		Ok:          true,
		Result:      messageAPIResponse,
		ErrorCode:   0,
		Description: "",
		Parameters:  nil,
	}

	apiResponseJSON, err := json.Marshal(&apiResponse)
	if err != nil {
		t.Error(err)
		return
	}

	ttPass := []struct {
		name                 string
		givenWeatherResponse *httpclient.GetWeatherResponse
		botResponses         []*http.Response
	}{
		{
			"existing location command",
			&httpclient.GetWeatherResponse{
				Name: "Grutendberger",
				Main: &httpclient.Main{
					Temp:    25.5,
					TempMax: 30.0,
					TempMin: 20.0,
				},
				Weather: []*httpclient.Weather{
					{Description: "It's fail, baby :)"},
				},
			},
			[]*http.Response{
				{
					StatusCode: 200,
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(strings.NewReader(authResponse)),
				},
				{
					StatusCode: 200,
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(strings.NewReader(string(apiResponseJSON))),
				},
			},
		},
	}

	conf := &config.Config{
		Host:   "https://api.openweathermap.org/data/2.5/weather?",
		KeyApi: "sdfsfsfsfssfsdfdfsfsfsfergtrgt",
	}
	for _, tc := range ttPass {

		responseJSON, err := json.Marshal(tc.givenWeatherResponse)
		if err != nil {
			t.Fatal(err)
		}

		httpWeatherClient := fakeHTTPBotClient(200, string(responseJSON))
		weatherClient := httpclient.NewWeatherClient(conf, httpWeatherClient, logger)
		bot := fakeBotWithWeatherClientMultipleResponses(weatherClient, tc.botResponses)
		bot.botClient.Debug = true
		go func() {
			if err = bot.ReplyingOnMessages(context.TODO()); err != nil {
				logger.Fatal(err)
			}
		}()
		time.Sleep(time.Second * 2)
	}
}

type testTransport2 struct {
	responses []*http.Response
	index     int
}

func (t *testTransport2) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.index >= len(t.responses) {
		return nil, errors.New("no more responses")
	}

	response := t.responses[t.index]
	t.index++
	return response, nil
}

func fakeHTTPBotClientWithMultipleResponses(responses []*http.Response) *http.Client {
	return &http.Client{
		Transport: &testTransport2{
			responses: responses,
			index:     0,
		},
	}
}

func fakeBotWithWeatherClientMultipleResponses(weatherClient *httpclient.WeatherClient, responses []*http.Response) *Bot {
	apiToken := os.Getenv("TOKEN")
	testConfig := &config.Config{
		Token:             apiToken,
		TimeoutMongoQuery: "5",
	}

	testTgClientHttp := fakeHTTPBotClientWithMultipleResponses(responses)

	ctl := gomock.NewController(&testing.T{})
	defer ctl.Finish()

	mockUserRepo := mock.NewMockUserRepoInterface(ctl)

	bot, err := NewBot(testConfig, testTgClientHttp, weatherClient, logger, mockUserRepo)
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	return bot
}

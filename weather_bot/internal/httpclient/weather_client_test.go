package httpclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"weather_bot/internal/config"
	"weather_bot/internal/log"

	"github.com/stretchr/testify/assert"
)

func TestAppendQueryParamsToWeatherApiHost(t *testing.T) {
	a := config.Config{Host: "https://testHost", KeyApi: "test_key"}
	httpClient := &http.Client{}
	logger, _ := log.NewLogAndSetLevel("debug")
	wcl1 := NewWeatherClient(&a, httpClient, logger)
	b := config.Config{Host: "https://testHost2", KeyApi: "test_key2"}
	wcl2 := NewWeatherClient(&b, httpClient, logger)
	testCase := []struct {
		weatherCl *WeatherClient
		latitude  float64
		longitude float64
		expected  string
	}{
		{
			weatherCl: wcl1,
			latitude:  50.390708,
			longitude: 30.627441,
			expected:  fmt.Sprintf("%vlat=%v&lon=%v&appid=%v", "https://testHost", 50.390708, 30.627441, "test_key&units=metric"),
		},
		{
			weatherCl: wcl2,
			latitude:  50.390708,
			longitude: 30.627441,
			expected:  fmt.Sprintf("%vlat=%v&lon=%v&appid=%v", "https://testHost2", 50.390708, 30.627441, "test_key2&units=metric"),
		},
	}

	for _, tc := range testCase {
		result := tc.weatherCl.appendQueryParamsToWeatherApiHost(tc.latitude, tc.longitude)
		assert.Equal(t, result, tc.expected,
			fmt.Sprintf("Incorrect result. Expected %v, got %v", result, tc.expected))
	}
}

var testMessage = []byte(`{"coord":{"lon":30.6274,"lat":50.3907},"weather":[{"id":804,"main":"Clouds","description":"overcast clouds","icon":"04d"}],"base":"stations","main":{"temp":294.73,"feels_like":295.27,"temp_min":294.32,"temp_max":295.95,"pressure":1019,"humidity":89},"visibility":10000,"wind":{"speed":0.45,"deg":65,"gust":1.34},"clouds":{"all":100},"dt":1693814064,"sys":{"type":2,"id":2082717,"country":"UA","sunrise":1693797349,"sunset":1693845480},"timezone":10800,"id":696377,"name":"Pozniaky","cod":200}`)

func TestGetWeatherForecastResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/valid-path" {
			response := testMessage
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		} else if r.URL.Path == "/invalid-path" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`invalid-json`))
		}

	}))

	defer server.Close()

	var getWeather *GetWeatherResponse
	a := config.Config{Host: "https://testHost", KeyApi: "test_key"}
	httpClient := &http.Client{}
	logger, _ := log.NewLogAndSetLevel("debug")
	wcl := NewWeatherClient(&a, httpClient, logger)
	testCase := []struct {
		path       string
		expected   *GetWeatherResponse
		expectErr  bool
		errMessage string
	}{
		{
			path:       server.URL + "/valid-path",
			expected:   getWeather,
			expectErr:  false,
			errMessage: "",
		},
		{
			path:       server.URL + "/invalid-path",
			expected:   getWeather,
			expectErr:  false,
			errMessage: "HTTP_SEND_REQUEST_ERR: Failed send request : []",
		},
	}

	for _, tc := range testCase {
		result, err := wcl.getWeatherForecastResponse(tc.path)
		if err != nil {
			assert.EqualError(t, err, tc.errMessage)
		} else {
			assert.IsType(t, result, tc.expected, fmt.Sprintf("Incorrect result. Expected %v, result %v", tc.expected, result))
		}
	}
}

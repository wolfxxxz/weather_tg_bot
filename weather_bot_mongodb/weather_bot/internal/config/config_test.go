package config

import (
	"fmt"

	"os"
	"testing"
	"weather_bot/internal/apperrors"
	"weather_bot/internal/log"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	testCase := []struct {
		path string
		err  apperrors.AppError
	}{
		{path: ".env", err: *apperrors.NewAppError()},
		{path: ".enb", err: *apperrors.EnvConfigParseError.AppendMessage(" open .enb: no such file or directory")},
	}

	for _, tc := range testCase {
		conf := NewConfig()
		logger, _ := log.NewLogAndSetLevel("info")
		err := conf.ParseConfig(tc.path, logger)
		if err != nil {
			assert.EqualError(t, err, tc.err.Error(),
				fmt.Sprintf("Incorrect result. Expected %v, got %v", tc.err.Error(), tc.err))
		} else {
			expected := os.Getenv("TOKEN")
			assert.Equal(t, expected, conf.Token,
				fmt.Sprintf("Incorrect result. Expected %v, got %v", expected, conf.Token))
		}
	}
}

func TestParseConfigConfFail(t *testing.T) {
	testCase := []struct {
		path string
		err  apperrors.AppError
	}{
		{path: ".env", err: *apperrors.EnvConfigParseError.AppendMessage("Expected a pointer to a Struct\n")},
	}

	for _, tc := range testCase {
		conf := NewConfig()
		conf = nil
		logger, _ := log.NewLogAndSetLevel("info")
		err := conf.ParseConfig(tc.path, logger)
		if err != nil {
			assert.EqualError(t, err, tc.err.Error(),
				fmt.Sprintf("Incorrect result. Expected %v, got %v", tc.err.Error(), tc.err))
		} else {
			expected := os.Getenv("TOKEN")
			assert.Equal(t, expected, conf.Token,
				fmt.Sprintf("Incorrect result. Expected %v, got %v", expected, conf.Token))
		}
	}
}

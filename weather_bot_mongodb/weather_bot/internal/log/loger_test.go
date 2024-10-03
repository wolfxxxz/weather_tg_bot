package log

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogAndSetLevel(t *testing.T) {
	testCase := []struct {
		logLevel string
		err      error
	}{
		{logLevel: "info"},
		{logLevel: "err", err: fmt.Errorf("not a valid logrus Level: \"%v\"", "err")},
	}
	for _, tc := range testCase {
		expectedLevel, _ := logrus.ParseLevel(tc.logLevel)
		result, resErr := NewLogAndSetLevel(tc.logLevel)

		if resErr != nil {
			assert.EqualError(t, resErr, tc.err.Error(), fmt.Sprintf("Incorrect result. Expected %v, got %v", tc.err.Error(), resErr))
		} else {
			logger := logrus.New()
			assert.EqualValues(t, result.Log, expectedLevel,
				fmt.Sprintf("Incorrect result. Expected %v, got %v", result, expectedLevel))
			assert.IsType(t, result.Log, logger,
				fmt.Sprintf("Incorrect result. Expected %T, got %T", result.Log, logger))
		}
	}
}

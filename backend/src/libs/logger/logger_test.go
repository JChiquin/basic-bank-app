package logger

import (
	"testing"

	"bank-service/src/libs/env"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetInstance(t *testing.T) {

	//Action
	instance := GetInstance()

	//Data Assertion
	assert.IsType(t, instance, &logrus.Entry{})
}

func TestGetLogLevel(t *testing.T) {
	testCases := []struct {
		TestName string
		Level    string
		Expected logrus.Level
	}{
		{
			TestName: "Error level",
			Level:    "error",
			Expected: logrus.ErrorLevel,
		},
		{
			TestName: "Warn level",
			Level:    "warn",
			Expected: logrus.WarnLevel,
		},
		{
			TestName: "Info level",
			Level:    "info",
			Expected: logrus.InfoLevel,
		},
		{
			TestName: "Http level",
			Level:    "http",
			Expected: logrus.InfoLevel,
		},
		{
			TestName: "Verbose level",
			Level:    "verbose",
			Expected: logrus.InfoLevel,
		},
		{
			TestName: "Debug level",
			Level:    "debug",
			Expected: logrus.DebugLevel,
		},
		{
			TestName: "Silly level",
			Level:    "silly",
			Expected: logrus.TraceLevel,
		},
		{
			TestName: "Default when value is unexpected",
			Level:    "anything",
			Expected: logrus.TraceLevel,
		},
		{
			TestName: "Default when value is empty",
			Level:    "",
			Expected: logrus.TraceLevel,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.TestName, func(t *testing.T) {
			oldValue := env.LogLevel
			env.LogLevel = tC.Level
			level := getLogLevel()

			assert.Equal(t, tC.Expected, level)

			t.Cleanup(func() {
				env.LogLevel = oldValue
			})
		})
	}
}

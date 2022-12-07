package logger

import (
	"os"

	"bank-service/src/libs/env"

	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
)

var (
	logger *logrus.Entry = logrus.NewEntry(logrus.StandardLogger()) //Default console logger

	// match possible strings to log levels
	logLevels map[string]logrus.Level = map[string]logrus.Level{
		"error":   logrus.ErrorLevel,
		"warn":    logrus.WarnLevel,
		"info":    logrus.InfoLevel,
		"http":    logrus.InfoLevel,
		"verbose": logrus.InfoLevel,
		"debug":   logrus.DebugLevel,
		"silly":   logrus.TraceLevel,
	}
)

/*
GetInstance returns the singleton instance
*/
func GetInstance() *logrus.Entry {
	return logger
}

/*
SetupLogger setup function to instantiate file logger
*/
func SetupLogger() {
	hostname, _ := os.Hostname()

	log := logrus.New()
	log.SetFormatter(&ecslogrus.Formatter{})

	// Setup Log Level
	log.SetLevel(getLogLevel())

	//Saving logger with default fields
	logger = log.WithFields(logrus.Fields{
		"hostname": hostname,
	})

	logger.Debug("Logger succesfully connected")
}

/*
getLogLevel returns the log level according to LOG_LEVEL env var
Default is 'silly' (most verbose)
*/
func getLogLevel() logrus.Level {
	level, ok := logLevels[env.LogLevel]
	if !ok {
		return logLevels["silly"]
	}
	return level
}

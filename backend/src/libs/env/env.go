package env

import (
	"bank-service/src/utils/constant"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (

	// AppEnv Application Environment
	AppEnv string

	// BankServiceSecondsBetweenAttempts BankService Interval in Seconds between attempts
	BankServiceSecondsBetweenAttempts time.Duration

	// BankServicePostgresqlHost BankService PostgreSQL host
	BankServicePostgresqlHost string

	// BankServicePostgresqlPort BankService PostgreSQL port
	BankServicePostgresqlPort string

	// BankServicePostgresqlName BankService PostgreSQL name
	BankServicePostgresqlName string

	// BankServicePostgresqlNameTest BankService PostgreSQL name Test
	BankServicePostgresqlNameTest string

	// BankServicePostgresqlUsername BankService PostgreSQL app username
	BankServicePostgresqlUsername string

	// BankServicePostgresqlPassword BankService PostgreSQL app password
	BankServicePostgresqlPassword string

	// BankServicePostgresqlSSLMode BankService PostgreSQL ssl mode
	BankServicePostgresqlSSLMode string

	// BankServiceRestPort BankService Rest port
	BankServiceRestPort string

	// LogLevel Log level
	LogLevel string

	// WhiteList White List
	WhiteList string

	// JWTSecretKey jwt secret key
	JWTSecretKey string
)

func init() {
	setEnvVars()
	setupUpdateEnvVars()
}

/*
setupUpdateEnvVars prepares the update of environment variables in environments other than testing
*/
func setupUpdateEnvVars() {
	if AppEnv != "testing" {
		go updateEnvVars()
	}
}

func setEnvVars() {
	// App Environment
	AppEnv = os.Getenv("APP_ENV")

	// BankService - Rest
	BankServiceRestPort = ProcessCriticalEnvVar("BANK_SERVICE_REST_PORT")

	// BankService Interval in Seconds Between Attempts
	var seconds int
	ProcessIntEnvVar(&seconds, "BANK_SERVICE_SECONDS_BETWEEN_ATTEMPTS", 60)
	BankServiceSecondsBetweenAttempts = time.Duration(seconds) * time.Second

	// BankService - PostgreSQL
	BankServicePostgresqlHost = ProcessCriticalEnvVar("BANK_SERVICE_POSTGRESQL_HOST")
	BankServicePostgresqlPort = ProcessCriticalEnvVar("BANK_SERVICE_POSTGRESQL_PORT")
	BankServicePostgresqlName = ProcessCriticalEnvVar("BANK_SERVICE_POSTGRESQL_NAME")
	BankServicePostgresqlNameTest = os.Getenv("BANK_SERVICE_POSTGRESQL_NAME_TEST")
	BankServicePostgresqlUsername = ProcessCriticalEnvVar("BANK_SERVICE_POSTGRESQL_USERNAME")
	BankServicePostgresqlPassword = ProcessCriticalEnvVar("BANK_SERVICE_POSTGRESQL_PASSWORD")
	BankServicePostgresqlSSLMode = os.Getenv("BANK_SERVICE_POSTGRESQL_SSLMODE")

	// Log level
	LogLevel = os.Getenv("LOG_LEVEL")

	// WhiteList
	WhiteList = os.Getenv("WHITE_LIST")

	// JWTSecretKey
	JWTSecretKey = ProcessCriticalEnvVar("JWT_SECRET_KEY")
}

// updateEnvVars updates the env vars periodically
func updateEnvVars() {
	for {
		time.Sleep(constant.IntervalBetweenEnvVarUpdate)
		setEnvVars()
	}
}

// ProcessIntEnvVar gets environment variable from os and parses it to int
func ProcessIntEnvVar(intVar *int, envKey string, defaultValue int) {
	var err error
	*intVar, err = strconv.Atoi(os.Getenv(envKey))
	if err != nil {
		*intVar = defaultValue
	}
}

// ProcessCriticalEnvVar gets environment variable from os and panics if it's not set
func ProcessCriticalEnvVar(envKey string) string {
	envVar := os.Getenv(envKey)
	if envVar == "" {
		panic(fmt.Sprintf("%s not set", envKey)) // Logger isn't available yet
	}
	return envVar
}

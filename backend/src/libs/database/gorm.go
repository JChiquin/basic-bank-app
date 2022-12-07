package database

import (
	"bank-service/src/libs/env"
	"fmt"
	"sync"
	"time"

	"bank-service/src/libs/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	ormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db   *gorm.DB
	once sync.Once
)

//CreateBankConnectionString returns the connection string based on environment variables
func CreateBankConnectionString() string {
	//db config vars
	dbHost := env.BankServicePostgresqlHost
	dbPort := env.BankServicePostgresqlPort
	dbName := env.BankServicePostgresqlName
	dbUser := env.BankServicePostgresqlUsername
	dbPassword := env.BankServicePostgresqlPassword
	dbSSLMode := env.BankServicePostgresqlSSLMode
	if env.AppEnv == "testing" {
		dbName = env.BankServicePostgresqlNameTest
	}
	//Make connection string with interpolation
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)
	return connectionString
}

/*
SetupBankGormDB open the pool connection in db var and return it
*/
func SetupBankGormDB() *gorm.DB {
	once.Do(func() {
		config := &gorm.Config{
			Logger: ormlogger.Default.LogMode(ormlogger.Info),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		}
		//connect to db
		var dbError error
		db, dbError = gorm.Open(postgres.Open(CreateBankConnectionString()), config)
		for dbError != nil {
			logger.GetInstance().Error("Failed to connect to own-database")
			time.Sleep(env.BankServiceSecondsBetweenAttempts)
			logger.GetInstance().Info("Retrying...")
			db, dbError = gorm.Open(postgres.Open(CreateBankConnectionString()), config)
		}
		logger.GetInstance().Info("Connected to own-database!")
		setConnectionMaxLifetime(db, time.Minute*5)
	})
	return db
}

/*
GetBankGormConnection return db pointer which already have an open connection
*/
func GetBankGormConnection() *gorm.DB {
	return SetupBankGormDB()
}

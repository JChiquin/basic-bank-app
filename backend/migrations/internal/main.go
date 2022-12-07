package main

import (
	"bank-service/src/libs/env"
	"crypto/tls"
	"log"
	"os"

	"github.com/go-pg/pg/v9"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

const directory = "migrations/internal"

func main() {
	dbHost := env.BankServicePostgresqlHost
	dbPort := env.BankServicePostgresqlPort
	dbName := env.BankServicePostgresqlName
	if env.AppEnv == "testing" {
		dbName = env.BankServicePostgresqlNameTest
	}
	dbUser := env.BankServicePostgresqlUsername
	dbPassword := env.BankServicePostgresqlPassword
	dbSSLMode := env.BankServicePostgresqlSSLMode

	options := &pg.Options{
		Addr:     dbHost + ":" + dbPort,
		User:     dbUser,
		Database: dbName,
		Password: dbPassword,
	}
	if dbSSLMode != "disable" {
		options.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	db := pg.Connect(options)

	err := migrations.Run(db, directory, os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

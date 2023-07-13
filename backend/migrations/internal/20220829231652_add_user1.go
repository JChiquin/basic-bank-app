package main

import (
	"bank-service/src/libs/env"
	"bank-service/src/libs/password"
	"bank-service/src/utils/constant"

	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

const userIDuser1 = 10000

func init() {
	const documentNumber = "20000000"
	up := func(db orm.DB) error {
		plainPassword := env.ProcessCriticalEnvVar("BANK_SERVICE_USER1_PASSWORD")
		email := env.ProcessCriticalEnvVar("BANK_SERVICE_USER1_EMAIL")
		hasedPassword, err := password.HashPassword(plainPassword)
		if err != nil {
			return err
		}
		_, err = db.Exec(`
			INSERT INTO "user" (id, first_name, last_name, email, password, document_number, account_number, user_type) VALUES(?, ?, ?, ?, ?, ?, ?, ?);
		`, userIDuser1, "Elon", "Musk", email, hasedPassword, documentNumber, "01298319273126387343", constant.UserTypeClient)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(`DELETE FROM "user" WHERE document_number = ?`, documentNumber)
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20220829231652_add_user1", up, down, opts)
}

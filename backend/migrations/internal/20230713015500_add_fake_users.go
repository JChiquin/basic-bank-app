package main

import (
	"bank-service/src/environments/client/resources/entity"
	"bank-service/src/libs/password"
	"bank-service/src/utils/constant"
	"encoding/json"
	"io"
	"os"

	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	file, err := os.Open("./fixture/fake_clients_data.json")
	if err != nil {
		panic(err)
	}

	fileJSON, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	users := []entity.User{}

	err = json.Unmarshal(fileJSON, &users)
	if err != nil {
		panic(err)
	}
	up := func(db orm.DB) error {

		for _, user := range users {
			password, err := password.HashPassword("12345678")
			if err != nil {
				return err
			}

			_, err = db.Exec(`
				INSERT INTO "user" (first_name, last_name, email, document_number, account_number, phone_number, birth_date, password ,user_type) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?);
			`, user.FirstName, user.LastName, user.Email, user.DocumentNumber, user.AccountNumber, user.PhoneNumber, user.BirthDate, password, constant.UserTypeClient)

			if err != nil {
				return err
			}
		}

		return nil
	}

	down := func(db orm.DB) error {
		for _, user := range users {
			_, err = db.Exec(`
				DELETE FROM "user" WHERE account_number = ?;
			`, user.AccountNumber)

			if err != nil {
				return err
			}
		}

		return nil
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20230713013541_add_fake_users", up, down, opts)
}

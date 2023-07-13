package main

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
			CREATE TABLE "user" (
				"id" serial,
				"first_name" varchar(40),
				"last_name" varchar(40),
				"email" varchar(50) UNIQUE,
				"birth_date" TIMESTAMP WITH TIME ZONE,
				"phone_number" varchar(20),
				"document_number" varchar(20) UNIQUE,
				"account_number" varchar(20) UNIQUE,
				"password" varchar(255),
				"user_type" varchar(20) NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
				updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
				deleted_at TIMESTAMP WITH TIME ZONE,
				PRIMARY KEY ("id"))
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(`DROP TABLE "user"`)
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20220727011546_create-user-table", up, down, opts)
}

package main

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
			CREATE TABLE "movement" (
				"id" serial,
				"user_id" INT NOT NULL,
				"amount" NUMERIC(16, 4) NOT NULL,
				"balance" NUMERIC(16, 4) NOT NULL,
				"multiplier" INT NOT NULL,
				"account_number" varchar(20) NOT NULL,
				"description" varchar(100) NOT NULL,

				created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
				updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
				deleted_at TIMESTAMP WITH TIME ZONE,
				
				PRIMARY KEY ("id"),
				CONSTRAINT fk_user
					FOREIGN KEY(user_id) 
						REFERENCES "user"(id) ON DELETE CASCADE
			);
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(`DROP TABLE "movement"`)
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20220729143904_create-movement-table", up, down, opts)
}

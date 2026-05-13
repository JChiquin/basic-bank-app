package main

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
			CREATE TABLE "password_reset_code" (
				"id" serial,
				"user_id" INT NOT NULL,
				"code_hash" varchar(255) NOT NULL,
				"expires_at" TIMESTAMP WITH TIME ZONE NOT NULL,
				"attempts" INT NOT NULL DEFAULT 0,
				"used_at" TIMESTAMP WITH TIME ZONE,
				created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
				updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
				deleted_at TIMESTAMP WITH TIME ZONE,
				PRIMARY KEY ("id"),
				CONSTRAINT fk_password_reset_code_user
					FOREIGN KEY(user_id)
						REFERENCES "user"(id) ON DELETE CASCADE
			);

			CREATE INDEX password_reset_code_user_id_idx ON "password_reset_code" ("user_id");
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(`DROP TABLE "password_reset_code"`)
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20260513094500_create_password_reset_code_table", up, down, opts)
}

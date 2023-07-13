package main

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
			CREATE FUNCTION create_bonus_movement() RETURNS TRIGGER AS $BODY$
			BEGIN
				INSERT INTO "movement" (user_id, amount, multiplier, account_number, description) 
				VALUES (NEW.id, 5000, 1, '83927623726321398231', 'Bono de bienvenida');

				RETURN NEW;
			END;
			$BODY$ LANGUAGE PLPGSQL;

			CREATE TRIGGER create_bonus_movement_trigger
			AFTER
			INSERT ON "user"
			FOR EACH ROW EXECUTE PROCEDURE create_bonus_movement();
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(`
			DROP TRIGGER IF EXISTS create_bonus_movement_trigger ON "user";
			DROP FUNCTION IF EXISTS create_bonus_movement;
		`)
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20220730153904_trigger_create_bonus_movement", up, down, opts)
}

package main

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
			CREATE FUNCTION set_balance_movement() RETURNS TRIGGER AS $BODY$
			DECLARE
			lastBalance movement.balance%TYPE;
			BEGIN
				SELECT balance INTO lastBalance
				FROM movement
				WHERE user_id = NEW.user_id
				ORDER BY created_at DESC
				LIMIT 1;

				NEW.balance := coalesce(lastBalance,0) + (NEW.multiplier * NEW.amount);

				RETURN NEW;
			END;
			$BODY$ LANGUAGE PLPGSQL;

			CREATE TRIGGER set_balance_movement_trigger
			BEFORE
			INSERT ON movement
			FOR EACH ROW EXECUTE PROCEDURE set_balance_movement();
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(`
			DROP TRIGGER IF EXISTS set_balance_movement_trigger ON movement;
			DROP FUNCTION IF EXISTS set_balance_movement;
		`)
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20220730143904_trigger_set_balance_movement", up, down, opts)
}

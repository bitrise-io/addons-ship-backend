package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190911094410, down20190911094410)
}

func up20190911094410(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE apps
		ADD COLUMN android_errors text[] DEFAULT ARRAY[]::text[],
		ADD COLUMN ios_errors text[] DEFAULT ARRAY[]::text[];`)
	return err
}

func down20190911094410(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE apps
		DROP COLUMN android_errors,
		DROP COLUMN ios_errors;`)
	return err
}

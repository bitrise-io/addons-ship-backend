package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190628063252, down20190628063252)
}

func up20190628063252(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions ADD COLUMN is_published boolean DEFAULT false;`)
	return err
}

func down20190628063252(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions DROP COLUMN is_published;`)
	return err
}

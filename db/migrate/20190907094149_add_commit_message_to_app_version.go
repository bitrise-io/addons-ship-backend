package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190907094149, down20190907094149)
}

func up20190907094149(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions ADD COLUMN commit_message text;`)
	return err
}

func down20190907094149(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions DROP COLUMN commit_message;`)
	return err
}

package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190607083123, down20190607083123)
}

func up20190607083123(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions ADD COLUMN scheme text;` +
		` ALTER TABLE app_versions ADD COLUMN configuration text;`)
	return err
}

func down20190607083123(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions DROP COLUMN scheme;` +
		` ALTER TABLE app_versions DROP COLUMN configuration;`)
	return err
}

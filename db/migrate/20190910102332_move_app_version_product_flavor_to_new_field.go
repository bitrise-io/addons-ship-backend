package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190910102332, down20190910102332)
}

func up20190910102332(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions ADD COLUMN product_flavour text;`)
	return err
}

func down20190910102332(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions DROP COLUMN product_flavour;`)
	return err
}

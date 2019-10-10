package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20191008114513, down20191008114513)
}

func up20191008114513(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions RENAME COLUMN product_flavour TO product_flavor;`)
	return err
}

func down20191008114513(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions RENAME COLUMN product_flavor TO product_flavour;`)
	return err
}

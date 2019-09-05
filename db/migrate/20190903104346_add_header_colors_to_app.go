package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190903104346, down20190903104346)
}

func up20190903104346(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE apps ADD COLUMN header_color_1 text;
	ALTER TABLE apps ADD COLUMN header_color_2 text;`)
	return err
}

func down20190903104346(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE apps DROP COLUMN header_color_1;
	ALTER TABLE apps DROP COLUMN header_color_2;`)
	return err
}

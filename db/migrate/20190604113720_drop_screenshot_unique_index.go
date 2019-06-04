package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190604113720, down20190604113720)
}

func up20190604113720(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP INDEX screenshots_filename_idx;`)
	return err
}

func down20190604113720(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE UNIQUE INDEX screenshots_filename_idx ON screenshots(filename, device_type, screen_size);`)
	return err
}

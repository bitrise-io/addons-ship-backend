package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190516071347, down20190516071347)
}

func up20190516071347(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions ADD COLUMN build_slug text;` +
		` ALTER TABLE app_versions ADD COLUMN whats_new text;`)
	return err
}

func down20190516071347(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions DROP COLUMN build_slug;` +
		` ALTER TABLE app_versions DROP COLUMN whats_new;`)
	return err
}

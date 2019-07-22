package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190719063048, down20190719063048)
}

func up20190719063048(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions ADD COLUMN app_info json DEFAULT '{}'::json;`)
	return err
}

func down20190719063048(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions DROP COLUMN app_info;`)
	return err
}

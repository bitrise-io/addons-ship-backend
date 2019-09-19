package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190918073336, down20190918073336)
}

func up20190918073336(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_version_events ADD COLUMN is_log_available boolean DEFAULT false;`)
	return err
}

func down20190918073336(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_version_events DROP COLUMN is_log_available;`)
	return err
}

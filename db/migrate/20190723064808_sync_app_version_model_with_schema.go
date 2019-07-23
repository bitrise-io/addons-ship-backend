package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190723064808, down20190723064808)
}

func up20190723064808(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions DROP COLUMN version;
    ALTER TABLE app_versions DROP COLUMN description;
    ALTER TABLE app_versions RENAME COLUMN app_info TO artifact_info;
    ALTER TABLE app_versions DROP COLUMN whats_new;`)
	return err
}

func down20190723064808(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions ADD COLUMN version text;
    ALTER TABLE app_versions ADD COLUMN description text;
    ALTER TABLE app_versions RENAME COLUMN artifact_info TO app_info;
    ALTER TABLE app_versions ADD COLUMN whats_new text;`)
	return err
}

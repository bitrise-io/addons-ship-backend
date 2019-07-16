package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190716131513, down20190716131513)
}

func up20190716131513(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE apps DROP COLUMN deleted_at;
    ALTER TABLE app_versions DROP COLUMN deleted_at;
    ALTER TABLE app_version_events DROP COLUMN deleted_at;
    ALTER TABLE app_contacts DROP COLUMN deleted_at;
    ALTER TABLE feature_graphics DROP COLUMN deleted_at;
    ALTER TABLE publish_tasks DROP COLUMN deleted_at;
    ALTER TABLE screenshots DROP COLUMN deleted_at;
    ALTER TABLE app_settings DROP COLUMN deleted_at;`)
	return err
}

func down20190716131513(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE apps ADD COLUMN deleted_at timestamp with time zone;
    ALTER TABLE app_versions ADD COLUMN deleted_at timestamp with time zone;
    ALTER TABLE app_version_events ADD COLUMN deleted_at timestamp with time zone;
    ALTER TABLE app_contacts ADD COLUMN deleted_at timestamp with time zone;
    ALTER TABLE feature_graphics ADD COLUMN deleted_at timestamp with time zone;
    ALTER TABLE publish_tasks ADD COLUMN deleted_at timestamp with time zone;
    ALTER TABLE screenshots ADD COLUMN deleted_at timestamp with time zone;
    ALTER TABLE app_settings ADD COLUMN deleted_at timestamp with time zone;`)
	return err
}

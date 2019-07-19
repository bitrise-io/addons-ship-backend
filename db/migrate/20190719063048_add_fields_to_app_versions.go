package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190719063048, down20190719063048)
}

func up20190719063048(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions ADD COLUMN minimum_os text;
	ALTER TABLE app_versions ADD COLUMN minimum_sdk text;
	ALTER TABLE app_versions ADD COLUMN certificate_expires_at timestamp with time zone;
	ALTER TABLE app_versions ADD COLUMN distribution_type text;
	ALTER TABLE app_versions ADD COLUMN app_info json DEFAULT '{}'::json;
	ALTER TABLE app_versions ADD COLUMN provisioning_info json DEFAULT '{}'::json;
	ALTER TABLE app_versions ADD COLUMN supported_device_types varying[] DEFAULT '{}'::character varying[];`)
	return err
}

func down20190719063048(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}

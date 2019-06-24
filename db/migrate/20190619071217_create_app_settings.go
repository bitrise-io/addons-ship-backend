package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190619071217, down20190619071217)
}

func up20190619071217(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE app_settings (
        id uuid primary key NOT NULL,
        app_id uuid NOT NULL REFERENCES apps (id),
        ios_settings json not null default '{}'::json,
        android_settings json not null default '{}'::json,
        ios_workflow text,
        android_workflow text,
        created_at timestamp with time zone NOT NULL,
        updated_at timestamp with time zone NOT NULL,
        deleted_at timestamp with time zone
    );`)
	return err
}

func down20190619071217(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE app_settings;`)
	return err
}

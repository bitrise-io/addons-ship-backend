package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190704083132, down20190704083132)
}

func up20190704083132(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE app_events;
	CREATE TABLE app_version_events (
        id uuid primary key NOT NULL,
		app_version_id uuid NOT NULL REFERENCES app_versions (id),
        status text NOT NULL,
        event_text text,
        created_at timestamp with time zone NOT NULL,
        updated_at timestamp with time zone NOT NULL,
        deleted_at timestamp with time zone
    );`)
	return err
}

func down20190704083132(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE app_version_events;
	CREATE TABLE app_events (
        id uuid primary key NOT NULL,
        app_id uuid NOT NULL REFERENCES apps (id),
        status text NOT NULL,
        created_at timestamp with time zone NOT NULL,
        updated_at timestamp with time zone NOT NULL,
        deleted_at timestamp with time zone
    );`)
	return err
}

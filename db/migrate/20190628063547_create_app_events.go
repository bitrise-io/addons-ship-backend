package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190628063547, down20190628063547)
}

func up20190628063547(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE app_events (
        id uuid primary key NOT NULL,
        app_id uuid NOT NULL REFERENCES apps (id),
        status text NOT NULL,
        created_at timestamp with time zone NOT NULL,
        updated_at timestamp with time zone NOT NULL,
        deleted_at timestamp with time zone
    );`)
	return err
}

func down20190628063547(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE app_events;`)
	return err
}

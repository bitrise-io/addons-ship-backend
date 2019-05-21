package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190511090020, down20190511090020)
}

func up20190511090020(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE apps (
        id uuid primary key NOT NULL,
        app_slug text NOT NULL,
        plan text,
        bitrise_api_token text,
        api_token text,
        created_at timestamp with time zone NOT NULL,
        updated_at timestamp with time zone NOT NULL,
        deleted_at timestamp with time zone
    );

    CREATE UNIQUE INDEX apps_app_slug_idx ON apps(app_slug);`)
	return err
}

func down20190511090020(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE apps;`)
	return err
}

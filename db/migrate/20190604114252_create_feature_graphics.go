package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190604114252, down20190604114252)
}

func up20190604114252(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE feature_graphics (
        id uuid primary key NOT NULL,
        app_version_id uuid NOT NULL REFERENCES app_versions (id),
        filename text NOT NULL,
        filesize integer,
        uploaded boolean,
        created_at timestamp with time zone NOT NULL,
        updated_at timestamp with time zone NOT NULL,
        deleted_at timestamp with time zone
    );`)
	return err
}

func down20190604114252(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE feature_graphics;`)
	return err
}

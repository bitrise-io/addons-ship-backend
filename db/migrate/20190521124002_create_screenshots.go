package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190521124002, down20190521124002)
}

func up20190521124002(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE screenshots (
        id uuid primary key NOT NULL,
        filename text NOT NULL,
        filesize integer,
        uploaded boolean,
        device_type text,
        screen_size text,
        created_at timestamp with time zone NOT NULL,
        updated_at timestamp with time zone NOT NULL,
        deleted_at timestamp with time zone
    );

    CREATE UNIQUE INDEX screenshots_filename_idx ON screenshots(filename, device_type, screen_size);`)
	return err
}

func down20190521124002(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE screenshots;`)
	return err
}

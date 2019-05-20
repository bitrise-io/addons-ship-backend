package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190513135553, down20190513135553)
}

func up20190513135553(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE app_versions (
        id uuid primary key NOT NULL,
        app_id uuid NOT NULL REFERENCES apps (id),
        version text NOT NULL,
        platform text,
        build_number text,
        description text,
        last_update timestamp with time zone,
        created_at timestamp with time zone NOT NULL,
        updated_at timestamp with time zone NOT NULL,
        deleted_at timestamp with time zone
    );

    CREATE INDEX app_versions_platform_idx ON app_versions(platform);`)
	return err
}

func down20190513135553(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE app_versions;`)
	return err
}

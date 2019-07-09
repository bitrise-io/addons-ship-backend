package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20190709130304, Down20190709130304)
}

func Up20190709130304(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE app_contacts (
		id uuid primary key NOT NULL,
		app_id uuid NOT NULL REFERENCES apps (id),
		email text NOT NULL,
		notification_preferences json not null default '{}'::json,
		confirmation_token text,
		confirmed_at timestamp with time zone,
		created_at timestamp with time zone NOT NULL,
		updated_at timestamp with time zone NOT NULL,
		deleted_at timestamp with time zone
	);`)
	return err
}

func Down20190709130304(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE app_contacts;`)
	return err
}

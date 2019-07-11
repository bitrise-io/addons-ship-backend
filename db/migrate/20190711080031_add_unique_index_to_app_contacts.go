package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190711080031, down20190711080031)
}

func up20190711080031(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE UNIQUE INDEX app_contacts_confirmation_token_idx ON app_contacts(confirmation_token);
    CREATE UNIQUE INDEX app_contacts_email_idx ON app_contacts(email);`)
	return err
}

func down20190711080031(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP INDEX app_contacts_confirmation_token_idx;
    DROP INDEX app_contacts_email_idx;`)
	return err
}

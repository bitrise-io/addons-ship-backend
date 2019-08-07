package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190806072243, down20190806072243)
}

func up20190806072243(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP INDEX app_contacts_email_idx;
	CREATE UNIQUE INDEX app_contacts_email_idx ON app_contacts(email, app_id);`)
	return err
}

func down20190806072243(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP INDEX app_contacts_email_idx;
	CREATE UNIQUE INDEX app_contacts_email_idx ON app_contacts(email);`)
	return err
}

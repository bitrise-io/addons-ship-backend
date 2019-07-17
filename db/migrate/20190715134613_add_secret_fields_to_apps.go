package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190715134613, down20190715134613)
}

func up20190715134613(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE apps ADD COLUMN encrypted_secret bytea;
    ALTER TABLE apps ADD COLUMN encrypted_secret_iv bytea;
    CREATE UNIQUE INDEX apps_encrypted_secret_iv_idx ON apps(encrypted_secret_iv);`)
	return err
}

func down20190715134613(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE apps DROP COLUMN encrypted_secret;
    ALTER TABLE apps DROP COLUMN encrypted_secret_iv;
    DROP INDEX apps_encrypted_secret_iv_idx;`)
	return err
}

package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190702140946, down20190702140946)
}

func up20190702140946(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE publish_tasks DROP COLUMN task_id;
	ALTER TABLE publish_tasks ADD COLUMN task_id uuid NOT NULL;`)
	return err
}

func down20190702140946(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE publish_tasks DROP COLUMN task_id;
	ALTER TABLE publish_tasks ADD COLUMN task_id text NOT NULL;`)
	return err
}

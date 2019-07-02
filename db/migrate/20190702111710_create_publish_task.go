package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190702111710, down20190702111710)
}

func up20190702111710(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE publish_tasks (
        id uuid primary key NOT NULL,
        app_version_id uuid NOT NULL REFERENCES apps (id),
        task_id text NOT NULL,
        created_at timestamp with time zone NOT NULL,
        updated_at timestamp with time zone NOT NULL,
        deleted_at timestamp with time zone
    );

    CREATE UNIQUE INDEX publish_tasks_task_id_idx ON publish_tasks(task_id);`)
	return err
}

func down20190702111710(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE publish_tasks;`)
	return err
}

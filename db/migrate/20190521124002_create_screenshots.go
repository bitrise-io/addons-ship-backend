package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190521124002, down20190521124002)
}

func up20190521124002(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return nil
}

func down20190521124002(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}

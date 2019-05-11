package main

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateApps, downCreateApps)
}

func upCreateApps(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return nil
}

func downCreateApps(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}

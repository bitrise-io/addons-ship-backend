package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190527103103, down20190527103103)
}

func up20190527103103(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions ADD COLUMN promotional_text text;` +
		` ALTER TABLE app_versions ADD COLUMN keywords text;` +
		` ALTER TABLE app_versions ADD COLUMN review_notes text;` +
		` ALTER TABLE app_versions ADD COLUMN support_url text;` +
		` ALTER TABLE app_versions ADD COLUMN marketing_url text;` +
		` ALTER TABLE app_versions ADD COLUMN scheme text;` +
		` ALTER TABLE app_versions ADD COLUMN configuration text;`)
	return err
}

func down20190527103103(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions DROP COLUMN promotional_text;` +
		` ALTER TABLE app_versions DROP COLUMN keywords;` +
		` ALTER TABLE app_versions DROP COLUMN review_notes;` +
		` ALTER TABLE app_versions DROP COLUMN support_url;` +
		` ALTER TABLE app_versions DROP COLUMN marketing_url;` +
		` ALTER TABLE app_versions DROP COLUMN scheme;` +
		` ALTER TABLE app_versions DROP COLUMN configuration;`)
	return err
}

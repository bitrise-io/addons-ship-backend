package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up20190729134726, down20190729134726)
}

func up20190729134726(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions
	DROP CONSTRAINT app_versions_app_id_fkey,
	ADD CONSTRAINT app_versions_app_id_fkey FOREIGN KEY (app_id) REFERENCES apps (id)
		ON DELETE CASCADE;
	ALTER TABLE app_settings
		DROP CONSTRAINT app_settings_app_id_fkey,
		ADD CONSTRAINT app_settings_app_id_fkey FOREIGN KEY (app_id) REFERENCES apps (id)
			ON DELETE CASCADE;
	ALTER TABLE app_contacts
		DROP CONSTRAINT app_contacts_app_id_fkey,
		ADD CONSTRAINT app_contacts_app_id_fkey FOREIGN KEY (app_id) REFERENCES apps (id)
			ON DELETE CASCADE;

	ALTER TABLE app_version_events
		DROP CONSTRAINT app_version_events_app_version_id_fkey,
		ADD CONSTRAINT app_version_events_app_version_id_fkey FOREIGN KEY (app_version_id) REFERENCES app_versions (id)
			ON DELETE CASCADE;
	ALTER TABLE feature_graphics
		DROP CONSTRAINT feature_graphics_app_version_id_fkey,
		ADD CONSTRAINT feature_graphics_app_version_id_fkey FOREIGN KEY (app_version_id) REFERENCES app_versions (id)
			ON DELETE CASCADE;
	ALTER TABLE publish_tasks
		DROP CONSTRAINT publish_tasks_app_version_id_fkey,
		ADD CONSTRAINT publish_tasks_app_version_id_fkey FOREIGN KEY (app_version_id) REFERENCES app_versions (id)
			ON DELETE CASCADE;
	ALTER TABLE screenshots
		DROP CONSTRAINT screenshots_app_version_id_fkey,
		ADD CONSTRAINT screenshots_app_version_id_fkey FOREIGN KEY (app_version_id) REFERENCES app_versions (id)
			ON DELETE CASCADE;
`)
	return err
}

func down20190729134726(tx *sql.Tx) error {
	_, err := tx.Exec(`ALTER TABLE app_versions
	DROP CONSTRAINT app_versions_app_id_fkey,
	ADD CONSTRAINT app_versions_app_id_fkey FOREIGN KEY (app_id) REFERENCES apps (id);
	ALTER TABLE app_settings
		DROP CONSTRAINT app_settings_app_id_fkey,
		ADD CONSTRAINT app_settings_app_id_fkey FOREIGN KEY (app_id) REFERENCES apps (id);
	ALTER TABLE app_contacts
		DROP CONSTRAINT app_contacts_app_id_fkey,
		ADD CONSTRAINT app_contacts_app_id_fkey FOREIGN KEY (app_id) REFERENCES apps (id);

	ALTER TABLE app_version_events
		DROP CONSTRAINT app_version_events_app_version_id_fkey,
		ADD CONSTRAINT app_version_events_app_version_id_fkey FOREIGN KEY (app_version_id) REFERENCES app_versions (id);
	ALTER TABLE feature_graphics
		DROP CONSTRAINT feature_graphics_app_version_id_fkey,
		ADD CONSTRAINT feature_graphics_app_version_id_fkey FOREIGN KEY (app_version_id) REFERENCES app_versions (id);
	ALTER TABLE publish_tasks
		DROP CONSTRAINT publish_tasks_app_version_id_fkey,
		ADD CONSTRAINT publish_tasks_app_version_id_fkey FOREIGN KEY (app_version_id) REFERENCES app_versions (id);
	ALTER TABLE screenshots
		DROP CONSTRAINT screenshots_app_version_id_fkey,
		ADD CONSTRAINT screenshots_app_version_id_fkey FOREIGN KEY (app_version_id) REFERENCES app_versions (id);
`)
	return err
}

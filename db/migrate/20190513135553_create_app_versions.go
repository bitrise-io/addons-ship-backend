package migration

import (
	"database/sql"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20190513135553, Down20190513135553)
}

func Up20190513135553(tx *sql.Tx) error {
	var err error
	db := dataservices.GetDB()
	if !db.HasTable(&models.AppVersion{}) {
		err = db.CreateTable(&models.AppVersion{}).Error
	}
	return err
}

func Down20190513135553(tx *sql.Tx) error {
	var err error
	db := dataservices.GetDB()
	if db.HasTable(&models.AppVersion{}) {
		err = db.DropTable(&models.AppVersion{}).Error
	}
	return err
}

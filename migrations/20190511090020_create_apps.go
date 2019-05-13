package main

import (
	"database/sql"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateApps, downCreateApps)
}

func upCreateApps(tx *sql.Tx) error {
	var err error
	db := dataservices.GetDB()
	if !db.HasTable(&models.App{}) {
		err = db.CreateTable(&models.App{}).Error
	}
	return err
}

func downCreateApps(tx *sql.Tx) error {
	var err error
	db := dataservices.GetDB()
	if db.HasTable(&models.App{}) {
		err = db.DropTable(&models.App{}).Error
	}
	return err
}

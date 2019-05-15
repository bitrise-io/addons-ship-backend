package models_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/jinzhu/gorm"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func recreateTestDB(t *testing.T) {
	dataservices.Close()

	err := dataservices.InitializeConnection(dataservices.ConnectionParams{}, false)
	if err != nil {
		fmt.Printf("Failed to initialize DB connection: %#v", err)
	}

	db := dataservices.GetDB()
	testDBName := os.Getenv("TEST_DB_NAME")
	// (re-)create test db
	err = db.Exec("DROP DATABASE IF EXISTS " + testDBName).Error
	if err != nil {
		fmt.Printf(" (!) Failed to DROP / destroy test db (%s), error: %s", testDBName, err)
	}

	panicIfErr(db.Exec("CREATE DATABASE " + testDBName).Error)
	dataservices.Close()

	err = dataservices.InitializeConnection(dataservices.ConnectionParams{DBName: testDBName}, true)
	if err != nil {
		fmt.Printf("Failed to initialize DB connection: %#v", err)
	}
	defer dataservices.Close()

	db = dataservices.GetDB()

	runTestMigrations(t, db)
}

func recreateAndInitTestDB(t *testing.T) {
	// create an empty database for tests
	recreateTestDB(t)

	testDBName := os.Getenv("TEST_DB_NAME")
	// re-init db connection, this time with the test db name
	panicIfErr(dataservices.InitializeConnection(dataservices.ConnectionParams{DBName: testDBName}, true))
}

func closeTestDB() {
	// close test DB
	dataservices.Close()
}

func prepareDB(t *testing.T) func() {
	t.Log("prepare DB")
	recreateAndInitTestDB(t)
	return closeTestDB
}

func runTestMigrations(t *testing.T, db *gorm.DB) {
	for _, migration := range []struct {
		message string
		fn      func() error
	}{
		{
			message: "create apps table",
			fn: func() error {
				if !db.HasTable(&models.App{}) {
					return db.CreateTable(&models.App{}).Error
				}
				return nil
			},
		},
		{
			message: "create app versions table",
			fn: func() error {
				if !db.HasTable(&models.AppVersion{}) {
					return db.CreateTable(&models.AppVersion{}).Error
				}
				return nil
			},
		},
	} {
		t.Log(migration.message)
		panicIfErr(migration.fn())
	}
}

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

func recreateTestDB() {
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

	runTestMigrations(db)
}

func recreateAndInitTestDB() {
	// create an empty database for tests
	recreateTestDB()

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
	recreateAndInitTestDB()
	return closeTestDB
}

func runTestMigrations(db *gorm.DB) {
	for _, migration := range []struct {
		fn func() error
	}{
		{
			fn: func() error {
				if !db.HasTable(&models.App{}) {
					return db.CreateTable(&models.App{}).Error
				}
				return nil
			},
		},
	} {
		panicIfErr(migration.fn())
	}
}

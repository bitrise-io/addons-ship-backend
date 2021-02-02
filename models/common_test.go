package models_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/go-utils/envutil"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
)

//nolint:unused,deadcode
func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

//nolint:unused,deadcode
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

//nolint:unused,deadcode
func recreateAndInitTestDB(t *testing.T) {
	// create an empty database for tests
	recreateTestDB(t)

	testDBName := os.Getenv("TEST_DB_NAME")
	// re-init db connection, this time with the test db name
	panicIfErr(dataservices.InitializeConnection(dataservices.ConnectionParams{DBName: testDBName}, true))
}

//nolint:unused,deadcode
func closeTestDB() {
	// close test DB
	dataservices.Close()
}

//nolint:unused,deadcode
func exportEnvVarsForTests(t *testing.T) {
	_, err := envutil.RevokableSetenv("APP_WEBHOOK_SECRET_ENCRYPT_KEY", "06042e86a7bd421c642c8c3e4ab13840")
	require.NoError(t, err)
}

//nolint:unused,deadcode
func prepareDB(t *testing.T) func() {
	t.Log("prepare DB")
	recreateAndInitTestDB(t)
	exportEnvVarsForTests(t)
	return closeTestDB
}

//nolint:unused,deadcode
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
		{
			message: "create screenshots table",
			fn: func() error {
				if !db.HasTable(&models.Screenshot{}) {
					return db.CreateTable(&models.Screenshot{}).Error
				}
				return nil
			},
		},
		{
			message: "create feature_graphics table",
			fn: func() error {
				if !db.HasTable(&models.FeatureGraphic{}) {
					return db.CreateTable(&models.FeatureGraphic{}).Error
				}
				return nil
			},
		},
		{
			message: "create app_settings table",
			fn: func() error {
				if !db.HasTable(&models.AppSettings{}) {
					return db.CreateTable(&models.AppSettings{}).Error
				}
				return nil
			},
		},
		{
			message: "create app_events table",
			fn: func() error {
				if !db.HasTable(&models.AppVersionEvent{}) {
					return db.CreateTable(&models.AppVersionEvent{}).Error
				}
				return nil
			},
		},
		{
			message: "create publish_tasks table",
			fn: func() error {
				if !db.HasTable(&models.PublishTask{}) {
					return db.CreateTable(&models.PublishTask{}).Error
				}
				return nil
			},
		},
		{
			message: "create publish_tasks table",
			fn: func() error {
				if !db.HasTable(&models.PublishTask{}) {
					return db.CreateTable(&models.PublishTask{}).Error
				}
				return nil
			},
		},
		{
			message: "create app_contacts table",
			fn: func() error {
				if !db.HasTable(&models.AppContact{}) {
					return db.CreateTable(&models.AppContact{}).Error
				}
				return nil
			},
		},
	} {
		t.Log(migration.message)
		panicIfErr(migration.fn())
	}
}

// +build database

package models_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func compareAppSettings(t *testing.T, expected, actual models.AppSettings) {
	expected.CreatedAt = time.Time{}
	expected.UpdatedAt = time.Time{}
	expected.App = nil
	actual.CreatedAt = time.Time{}
	actual.UpdatedAt = time.Time{}
	actual.App = nil
	require.Equal(t, expected, actual)
}

func Test_AppSettingsService_Find(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appSettingsService := models.AppSettingsService{DB: dataservices.GetDB()}

	testApp := createTestApp(t, &models.App{AppSlug: "test-app-slug"})
	testAppSettings := createTestAppSettings(t, &models.AppSettings{App: testApp, AndroidWorkflow: "android-deploy"})

	t.Run("when querying app settings that belongs to an app", func(t *testing.T) {
		foundAppSettings, err := appSettingsService.Find(&models.AppSettings{Record: models.Record{ID: testAppSettings.ID}, AppID: testApp.ID})
		require.NoError(t, err)
		compareAppSettings(t, *testAppSettings, *foundAppSettings)
	})

	t.Run("error - when feature graphic is not found", func(t *testing.T) {
		otherTestApp := createTestApp(t, &models.App{AppSlug: "test-app-slug-2"})

		foundAppSettings, err := appSettingsService.Find(&models.AppSettings{Record: models.Record{ID: testAppSettings.ID}, AppID: otherTestApp.ID})
		require.Equal(t, errors.Cause(err), gorm.ErrRecordNotFound)
		require.Nil(t, foundAppSettings)
	})
}

func Test_AppSettingsService_Update(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appSettingsService := models.AppSettingsService{DB: dataservices.GetDB()}

	t.Run("ok", func(t *testing.T) {
		testAppSettings := []*models.AppSettings{
			createTestAppSettings(t, &models.AppSettings{IosWorkflow: "my-ios-wf"}),
			createTestAppSettings(t, &models.AppSettings{IosWorkflow: "awesome-ios-wf"}),
		}

		testAppSettings[0].IosSettingsData = json.RawMessage(`{"app_sku": "20180601"}`)
		verrs, err := appSettingsService.Update(testAppSettings[0], []string{"IosSettingsData"})
		require.Empty(t, verrs)
		require.NoError(t, err)

		t.Log("check if app setting got updated")
		foundAppSettings, err := appSettingsService.Find(&models.AppSettings{Record: models.Record{ID: testAppSettings[0].ID}})
		require.NoError(t, err)

		foundIosSettings, err := foundAppSettings.IosSettings()
		require.NoError(t, err)
		require.Equal(t, "20180601", foundIosSettings.AppSKU)

		t.Log("check if no other app settings were updated")
		foundAppSettings, err = appSettingsService.Find(&models.AppSettings{Record: models.Record{ID: testAppSettings[1].ID}})
		require.NoError(t, err)
		compareAppSettings(t, *testAppSettings[1], *foundAppSettings)
	})
}

// +build database

package models_test

import (
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
)

func Test_AppVersionService_Create(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appVersionService := models.AppVersionService{DB: dataservices.GetDB()}
	testAppVersion := &models.AppVersion{
		Version: "v1.0",
	}
	createdAppVersion, err := appVersionService.Create(testAppVersion)
	require.NoError(t, err)
	require.False(t, createdAppVersion.ID.String() == "")
	require.False(t, createdAppVersion.CreatedAt.String() == "")
	require.False(t, createdAppVersion.UpdatedAt.String() == "")
}

func Test_AppVersionService_Find(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appVersionService := models.AppVersionService{DB: dataservices.GetDB()}
	testAppVersion := testAppVersion(t, &models.AppVersion{
		App: *testApp(t, &models.App{}),
	})

	foundAppVersion, err := appVersionService.Find(testAppVersion)
	require.NoError(t, err)
	require.Equal(t, testAppVersion, foundAppVersion)
}

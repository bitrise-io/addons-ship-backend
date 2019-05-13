// +build database

package models_test

import (
	"testing"
	"reflect"

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
	testAppVersion := createTestAppVersion(t, &models.AppVersion{
		App: *createTestApp(t, &models.App{}),
	})

	foundAppVersion, err := appVersionService.Find(testAppVersion)
	require.NoError(t, err)
	require.Equal(t, testAppVersion, foundAppVersion)
}

func Test_AppVersionService_FindAll(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appVersionService := models.AppVersionService{DB: dataservices.GetDB()}
	testApp1 := createTestApp(t, &models.App{})
	testApp1VersionAndroid := createTestAppVersion(t, &models.AppVersion{
		App: *testApp1,
		Platform: "android",
	})
	testApp1VersionIOS := createTestAppVersion(t, &models.AppVersion{
		App: *testApp1,
		Platform: "ios",
	})

	t.Run("when query all versions of test app 1", func(t *testing.T){
		foundAppVersions, err := appVersionService.FindAll(testApp1, map[string]interface{}{})
		require.NoError(t, err)
		reflect.DeepEqual([]models.AppVersion{*testApp1VersionIOS, *testApp1VersionAndroid}, foundAppVersions)
	})

	testApp2 := createTestApp(t, &models.App{})
	createTestAppVersion(t, &models.AppVersion{
		App: *testApp2,
		Platform: "ios",
	})

	t.Run("when query ios versions of test app 1", func(t *testing.T){
		foundAppVersions, err := appVersionService.FindAll(testApp1, map[string]interface{}{})
		require.NoError(t, err)
		reflect.DeepEqual([]models.AppVersion{*testApp1VersionIOS}, foundAppVersions)
	})
}

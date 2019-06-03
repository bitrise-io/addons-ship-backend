// +build database

package models_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
	"github.com/pkg/errors"
)

func compareAppVersion(t *testing.T, expected, actual models.AppVersion) {
	expected.CreatedAt = time.Time{}
	expected.UpdatedAt = time.Time{}
	expected.LastUpdate = time.Time{}
	actual.CreatedAt = time.Time{}
	actual.UpdatedAt = time.Time{}
	actual.LastUpdate = time.Time{}
	require.Equal(t, expected, actual)
}

func Test_AppVersionService_Create(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appVersionService := models.AppVersionService{DB: dataservices.GetDB()}
	t.Run("ok", func(t *testing.T) {
		testAppVersion := &models.AppVersion{
			Version:          "v1.0",
			AppStoreInfoData: json.RawMessage(`{"short_description":"Some quite short description"}`),
		}
		expectedAppStoreInfo := models.AppStoreInfo{
			ShortDescription: "Some quite short description",
		}
		createdAppVersion, verrs, err := appVersionService.Create(testAppVersion)
		require.NoError(t, err)
		require.Empty(t, verrs)
		require.False(t, createdAppVersion.ID.String() == "")
		require.False(t, createdAppVersion.CreatedAt.String() == "")
		require.False(t, createdAppVersion.UpdatedAt.String() == "")

		createdAppVersionStoreInfo, err := testAppVersion.AppStoreInfo()
		require.NoError(t, err)
		require.Equal(t, expectedAppStoreInfo, createdAppVersionStoreInfo)
	})

	t.Run("when app store info is not a valid JSON", func(t *testing.T) {
		testAppVersion := &models.AppVersion{
			Platform:         "ios",
			AppStoreInfoData: json.RawMessage(`invalid json`),
		}
		createdAppVersion, verrs, err := appVersionService.Create(testAppVersion)
		require.Empty(t, verrs)
		require.EqualError(t, err, "invalid character 'i' looking for beginning of value")
		require.Nil(t, createdAppVersion)
	})

	t.Run("when platform is android", func(t *testing.T) {
		t.Run("when short description is longer, than 80 characters", func(t *testing.T) {
			testAppVersion := &models.AppVersion{
				Platform:         "android",
				AppStoreInfoData: json.RawMessage(`{"short_description":"Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula e"}`),
			}
			createdAppVersion, verrs, err := appVersionService.Create(testAppVersion)
			require.Equal(t, []error{errors.New("short_description: Mustn't be longer than 80 characters")}, verrs)
			require.NoError(t, err)
			require.Nil(t, createdAppVersion)
		})

		t.Run("when full description is longer, than 80 characters", func(t *testing.T) {
			testAppVersion := &models.AppVersion{
				Platform:         "android",
				AppStoreInfoData: json.RawMessage(`{"full_description":"Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula e"}`),
			}
			createdAppVersion, verrs, err := appVersionService.Create(testAppVersion)
			require.Equal(t, []error{errors.New("full_description: Mustn't be longer than 80 characters")}, verrs)
			require.NoError(t, err)
			require.Nil(t, createdAppVersion)
		})
	})

	t.Run("when platform is ios", func(t *testing.T) {
		t.Run("when short description is longer, than 255 characters", func(t *testing.T) {
			testAppVersion := &models.AppVersion{
				Platform:         "ios",
				AppStoreInfoData: json.RawMessage(`{"short_description":"Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis,."}`),
			}
			createdAppVersion, verrs, err := appVersionService.Create(testAppVersion)
			require.Equal(t, []error{errors.New("short_description: Mustn't be longer than 255 characters")}, verrs)
			require.NoError(t, err)
			require.Nil(t, createdAppVersion)
		})
	})
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
		App:      *testApp1,
		Platform: "android",
	})
	testApp1VersionIOS := createTestAppVersion(t, &models.AppVersion{
		App:      *testApp1,
		Platform: "ios",
	})

	t.Run("when query all versions of test app 1", func(t *testing.T) {
		foundAppVersions, err := appVersionService.FindAll(testApp1, map[string]interface{}{})
		require.NoError(t, err)
		reflect.DeepEqual([]models.AppVersion{*testApp1VersionIOS, *testApp1VersionAndroid}, foundAppVersions)
	})

	testApp2 := createTestApp(t, &models.App{})
	createTestAppVersion(t, &models.AppVersion{
		App:      *testApp2,
		Platform: "ios",
	})

	t.Run("when query ios versions of test app 1", func(t *testing.T) {
		foundAppVersions, err := appVersionService.FindAll(testApp1, map[string]interface{}{})
		require.NoError(t, err)
		reflect.DeepEqual([]models.AppVersion{*testApp1VersionIOS}, foundAppVersions)
	})
}

func Test_AppVersionService_Update(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appVersionService := models.AppVersionService{DB: dataservices.GetDB()}

	t.Run("ok", func(t *testing.T) {
		testAppVersions := []*models.AppVersion{
			createTestAppVersion(t, &models.AppVersion{Platform: "iOS", Version: "v1.0"}),
			createTestAppVersion(t, &models.AppVersion{Platform: "Android", Version: "v1.2"}),
		}

		testAppVersions[0].AppStoreInfoData = json.RawMessage(`{"short_description": "Some short description"}`)
		verrs, err := appVersionService.Update(testAppVersions[0], []string{"AppStoreInfoData"})
		require.Empty(t, verrs)
		require.NoError(t, err)

		t.Log("check if app version got updated")
		foundAppVersion, err := appVersionService.Find(&models.AppVersion{Record: models.Record{ID: testAppVersions[0].ID}})
		require.NoError(t, err)

		foundAppStoreInfo, err := foundAppVersion.AppStoreInfo()
		require.NoError(t, err)
		require.Equal(t, "Some short description", foundAppStoreInfo.ShortDescription)

		t.Log("check if no other app version were updated")
		foundAppVersion, err = appVersionService.Find(&models.AppVersion{Record: models.Record{ID: testAppVersions[1].ID}})
		require.NoError(t, err)
		compareAppVersion(t, *testAppVersions[1], *foundAppVersion)
	})

	t.Run("when short description is longer than 80 characters", func(t *testing.T) {
		testAppVersion := createTestAppVersion(t, &models.AppVersion{Platform: "android", Version: "v1.0"})
		testAppVersion.AppStoreInfoData = json.RawMessage(`{"short_description":"Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula e"}`)
		verrs, err := appVersionService.Update(testAppVersion, []string{"AppStoreInfoData"})
		require.Equal(t, 1, len(verrs))
		require.Equal(t, "short_description: Mustn't be longer than 80 characters", verrs[0].Error())
		require.NoError(t, err)
	})

	t.Run("when trying to update non-existing field", func(t *testing.T) {
		testAppVersion := createTestAppVersion(t, &models.AppVersion{Platform: "iOS", Version: "v1.0"})
		verrs, err := appVersionService.Update(testAppVersion, []string{"NonExistingField"})
		require.EqualError(t, err, "Attribute name doesn't exist in the model")
		require.Equal(t, 0, len(verrs))
	})
}

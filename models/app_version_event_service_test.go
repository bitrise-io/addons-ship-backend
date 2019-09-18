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
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func compareAppVersionEvent(t *testing.T, expected, actual models.AppVersionEvent) {
	expected.CreatedAt = time.Time{}
	expected.UpdatedAt = time.Time{}
	expected.AppVersion = (*models.AppVersion)(nil)
	actual.CreatedAt = time.Time{}
	actual.UpdatedAt = time.Time{}
	actual.AppVersion = (*models.AppVersion)(nil)
	require.Equal(t, expected, actual)
}

func Test_AppVersionEventService_Create(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appVersionEventService := models.AppVersionEventService{DB: dataservices.GetDB()}
	testAppVersionEvent := &models.AppVersionEvent{Text: "Some interesting event"}

	createdAppVersionEvent, err := appVersionEventService.Create(testAppVersionEvent)
	require.NoError(t, err)
	require.False(t, createdAppVersionEvent.ID.String() == "")
	require.False(t, createdAppVersionEvent.CreatedAt.String() == "")
	require.False(t, createdAppVersionEvent.UpdatedAt.String() == "")
}

func Test_AppVersionEventService_Find(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appVersionEventService := models.AppVersionEventService{DB: dataservices.GetDB()}
	testAppVersion := createTestAppVersion(t, &models.AppVersion{Platform: "ios", ArtifactInfoData: json.RawMessage(`{"version":"1.0"}`)})
	testAppVersionEvent := createTestAppVersionEvent(t, &models.AppVersionEvent{Text: "Some interesting event", AppVersion: *testAppVersion})

	t.Run("when querying a app event that belongs to an app", func(t *testing.T) {
		foundAppVersionEvent, err := appVersionEventService.Find(&models.AppVersionEvent{Record: models.Record{ID: testAppVersionEvent.ID}, AppVersionID: testAppVersion.ID})
		require.NoError(t, err)
		reflect.DeepEqual(testAppVersionEvent, foundAppVersionEvent)
	})

	t.Run("error - when app event is not found", func(t *testing.T) {
		otherTestAppVersion := createTestAppVersion(t, &models.AppVersion{Platform: "android", ArtifactInfoData: json.RawMessage(`{"version":"1.0"}`)})

		foundAppVersionEvent, err := appVersionEventService.Find(&models.AppVersionEvent{Record: models.Record{ID: testAppVersionEvent.ID}, AppVersionID: otherTestAppVersion.ID})
		require.Equal(t, errors.Cause(err), gorm.ErrRecordNotFound)
		require.Nil(t, foundAppVersionEvent)
	})
}

func Test_AppVersionEventService_FindAll(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appVersionEventService := models.AppVersionEventService{DB: dataservices.GetDB()}
	testAppVersion := createTestAppVersion(t, &models.AppVersion{Platform: "ios", ArtifactInfoData: json.RawMessage(`{"version":"1.0"}`)})
	otherTestAppVersion := createTestAppVersion(t, &models.AppVersion{Platform: "android", ArtifactInfoData: json.RawMessage(`{"version":"1.0"}`)})
	testAppVersionEvent1 := createTestAppVersionEvent(t, &models.AppVersionEvent{Text: "Some interesting event", AppVersion: *testAppVersion})
	testAppVersionEvent2 := createTestAppVersionEvent(t, &models.AppVersionEvent{Text: "Some other interesting event", AppVersion: *testAppVersion})
	createTestAppVersionEvent(t, &models.AppVersionEvent{Text: "Some other interesting event", AppVersion: *otherTestAppVersion})

	t.Run("when query all app events of test app", func(t *testing.T) {
		foundAppVersionEvents, err := appVersionEventService.FindAll(testAppVersion)
		require.NoError(t, err)
		reflect.DeepEqual([]models.AppVersionEvent{*testAppVersionEvent2, *testAppVersionEvent1}, foundAppVersionEvents)
	})
}

func Test_AppVersionEventService_Update(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appVersionEventService := models.AppVersionEventService{DB: dataservices.GetDB()}
	testAppVersion := createTestAppVersion(t, &models.AppVersion{Platform: "ios", ArtifactInfoData: json.RawMessage(`{"version":"1.0"}`)})

	t.Run("ok", func(t *testing.T) {
		testAppVersionEvents := []*models.AppVersionEvent{
			createTestAppVersionEvent(t, &models.AppVersionEvent{Text: "Some interesting event", AppVersion: *testAppVersion}),
			createTestAppVersionEvent(t, &models.AppVersionEvent{Text: "Another interesting event", AppVersion: *testAppVersion}),
		}

		testAppVersionEvents[0].IsLogAvailable = true
		verrs, err := appVersionEventService.Update(testAppVersionEvents[0], []string{"IsLogAvailable"})
		require.Empty(t, verrs)
		require.NoError(t, err)

		t.Log("check if app version event got updated")
		foundAppVersionEvent, err := appVersionEventService.Find(&models.AppVersionEvent{Record: models.Record{ID: testAppVersionEvents[0].ID}})
		require.NoError(t, err)
		require.Equal(t, true, foundAppVersionEvent.IsLogAvailable)

		t.Log("check if no other app version events were updated")
		foundAppVersionEvent, err = appVersionEventService.Find(&models.AppVersionEvent{Record: models.Record{ID: testAppVersionEvents[1].ID}})
		require.NoError(t, err)
		compareAppVersionEvent(t, *testAppVersionEvents[1], *foundAppVersionEvent)
	})

	t.Run("when trying to update non-existing field", func(t *testing.T) {
		testAppVersionEvent := createTestAppVersionEvent(t, &models.AppVersionEvent{})
		verrs, err := appVersionEventService.Update(testAppVersionEvent, []string{"NonExistingField"})
		require.EqualError(t, err, "Attribute name doesn't exist in the model")
		require.Equal(t, 0, len(verrs))
	})
}

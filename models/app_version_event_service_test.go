// +build database

package models_test

import (
	"reflect"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

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
	testAppVersion := createTestAppVersion(t, &models.AppVersion{Platform: "ios"})
	testAppVersionEvent := createTestAppVersionEvent(t, &models.AppVersionEvent{Text: "Some interesting event", AppVersion: *testAppVersion})

	t.Run("when querying a app event that belongs to an app", func(t *testing.T) {
		foundAppVersionEvent, err := appVersionEventService.Find(&models.AppVersionEvent{Record: models.Record{ID: testAppVersionEvent.ID}, AppVersionID: testAppVersion.ID})
		require.NoError(t, err)
		reflect.DeepEqual(testAppVersionEvent, foundAppVersionEvent)
	})

	t.Run("error - when app event is not found", func(t *testing.T) {
		otherTestAppVersion := createTestAppVersion(t, &models.AppVersion{Platform: "android"})

		foundAppVersionEvent, err := appVersionEventService.Find(&models.AppVersionEvent{Record: models.Record{ID: testAppVersionEvent.ID}, AppVersionID: otherTestAppVersion.ID})
		require.Equal(t, errors.Cause(err), gorm.ErrRecordNotFound)
		require.Nil(t, foundAppVersionEvent)
	})
}

func Test_AppVersionEventService_FindAll(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appVersionEventService := models.AppVersionEventService{DB: dataservices.GetDB()}
	testAppVersion := createTestAppVersion(t, &models.AppVersion{Platform: "ios"})
	otherTestAppVersion := createTestAppVersion(t, &models.AppVersion{Platform: "android"})
	testAppVersionEvent1 := createTestAppVersionEvent(t, &models.AppVersionEvent{Text: "Some interesting event", AppVersion: *testAppVersion})
	testAppVersionEvent2 := createTestAppVersionEvent(t, &models.AppVersionEvent{Text: "Some other interesting event", AppVersion: *testAppVersion})
	createTestAppVersionEvent(t, &models.AppVersionEvent{Text: "Some other interesting event", AppVersion: *otherTestAppVersion})

	t.Run("when query all app events of test app", func(t *testing.T) {
		foundAppVersionEvents, err := appVersionEventService.FindAll(testAppVersion)
		require.NoError(t, err)
		reflect.DeepEqual([]models.AppVersionEvent{*testAppVersionEvent2, *testAppVersionEvent1}, foundAppVersionEvents)
	})
}

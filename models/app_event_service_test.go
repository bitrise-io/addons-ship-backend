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

func Test_AppEventService_Create(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appEventService := models.AppEventService{DB: dataservices.GetDB()}
	testAppEvent := &models.AppEvent{Text: "Some interesting event"}

	createdAppEvent, err := appEventService.Create(testAppEvent)
	require.NoError(t, err)
	require.False(t, createdAppEvent.ID.String() == "")
	require.False(t, createdAppEvent.CreatedAt.String() == "")
	require.False(t, createdAppEvent.UpdatedAt.String() == "")
}

func Test_AppEventService_Find(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appEventService := models.AppEventService{DB: dataservices.GetDB()}
	testApp := createTestApp(t, &models.App{AppSlug: "test-app-slug"})
	testAppEvent := createTestAppEvent(t, &models.AppEvent{Text: "Some interesting event", App: *testApp})

	t.Run("when querying a app event that belongs to an app", func(t *testing.T) {
		foundAppEvent, err := appEventService.Find(&models.AppEvent{Record: models.Record{ID: testAppEvent.ID}, AppID: testApp.ID})
		require.NoError(t, err)
		reflect.DeepEqual(testAppEvent, foundAppEvent)
	})

	t.Run("error - when app event is not found", func(t *testing.T) {
		otherTestApp := createTestApp(t, &models.App{AppSlug: "test-app-slug-2"})

		foundAppEvent, err := appEventService.Find(&models.AppEvent{Record: models.Record{ID: testAppEvent.ID}, AppID: otherTestApp.ID})
		require.Equal(t, errors.Cause(err), gorm.ErrRecordNotFound)
		require.Nil(t, foundAppEvent)
	})
}

func Test_AppEventService_FindAll(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appEventService := models.AppEventService{DB: dataservices.GetDB()}
	testApp := createTestApp(t, &models.App{AppSlug: "test-app-slug"})
	otherTestApp := createTestApp(t, &models.App{AppSlug: "test-app-slug-2"})
	testAppEvent1 := createTestAppEvent(t, &models.AppEvent{Text: "Some interesting event", App: *testApp})
	testAppEvent2 := createTestAppEvent(t, &models.AppEvent{Text: "Some other interesting event", App: *testApp})
	createTestAppEvent(t, &models.AppEvent{Text: "Some other interesting event", App: *otherTestApp})

	t.Run("when query all app events of test app", func(t *testing.T) {
		foundAppEvents, err := appEventService.FindAll(testApp)
		require.NoError(t, err)
		reflect.DeepEqual([]models.AppEvent{*testAppEvent2, *testAppEvent1}, foundAppEvents)
	})
}

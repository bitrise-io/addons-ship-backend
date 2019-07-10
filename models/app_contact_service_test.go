// +build database

package models_test

import (
	"encoding/json"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_AppContactService_Create(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appContactService := models.AppContactService{DB: dataservices.GetDB()}
	testAppContact := &models.AppContact{Email: "an-email@addr.ess"}

	createdAppContact, err := appContactService.Create(testAppContact)
	require.NoError(t, err)
	require.False(t, createdAppContact.ID.String() == "")
	require.False(t, createdAppContact.CreatedAt.String() == "")
	require.False(t, createdAppContact.UpdatedAt.String() == "")
}

func Test_AppContactService_Find(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appContactService := models.AppContactService{DB: dataservices.GetDB()}

	testApp := createTestApp(t, &models.App{AppSlug: "test-app-slug"})
	testAppContact := createTestAppContact(t, &models.AppContact{App: testApp, Email: "an-email@addr.ess"})

	t.Run("when querying an app contact that belongs to an app", func(t *testing.T) {
		foundAppContact, err := appContactService.Find(&models.AppContact{Record: models.Record{ID: testAppContact.ID}, AppID: testApp.ID})
		require.NoError(t, err)
		compareAppContacts(t, *testAppContact, *foundAppContact)
	})

	t.Run("error - when app contact is not found", func(t *testing.T) {
		otherTestApp := createTestApp(t, &models.App{AppSlug: "test-app-slug-2"})

		foundAppContact, err := appContactService.Find(&models.AppContact{Record: models.Record{ID: testAppContact.ID}, AppID: otherTestApp.ID})
		require.Equal(t, errors.Cause(err), gorm.ErrRecordNotFound)
		require.Nil(t, foundAppContact)
	})
}

func Test_AppContactService_Update(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appContactService := models.AppContactService{DB: dataservices.GetDB()}

	t.Run("ok", func(t *testing.T) {
		testAppContacts := []*models.AppContact{
			createTestAppContact(t, &models.AppContact{Email: "an-email@addr.ess"}),
			createTestAppContact(t, &models.AppContact{Email: "other-email@addr.ess"}),
		}

		testAppContacts[0].NotificationPreferencesData = json.RawMessage(`{"new_version": true}`)

		err := appContactService.Update(testAppContacts[0], []string{"NotificationPreferencesData"})
		require.NoError(t, err)

		t.Log("Check if app contact got updated")
		foundAppContact, err := appContactService.Find(&models.AppContact{Record: models.Record{ID: testAppContacts[0].ID}})
		require.NoError(t, err)

		notificationPrefs, err := foundAppContact.NotificationPreferences()
		require.NoError(t, err)
		require.Equal(t, notificationPrefs.NewVersion, true)

		t.Log("check if no other app contact was updated")
		foundAppContact, err = appContactService.Find(&models.AppContact{Record: models.Record{ID: testAppContacts[1].ID}})
		require.NoError(t, err)
		compareAppContacts(t, *testAppContacts[1], *foundAppContact)
	})
}

func Test_AppContactService_Delete(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appContactService := models.AppContactService{DB: dataservices.GetDB()}
	appContact := createTestAppContact(t, &models.AppContact{Email: "an-email@addr.ess"})

	t.Run("ok", func(t *testing.T) {
		err := appContactService.Delete(&models.AppContact{Record: models.Record{ID: appContact.ID}})
		require.NoError(t, err)
	})

	t.Run("error - when app contact is not found", func(t *testing.T) {
		err := appContactService.Delete(&models.AppContact{Record: models.Record{ID: uuid.NewV4()}})
		require.Equal(t, err, gorm.ErrRecordNotFound)
	})
}

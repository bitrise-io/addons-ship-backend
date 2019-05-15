// +build database

package models_test

import (
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
)

func Test_AppService_Create(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appService := models.AppService{DB: dataservices.GetDB()}
	testApp := &models.App{
		AppSlug: "test-app_slug",
	}
	createdApp, err := appService.Create(testApp)
	require.NoError(t, err)
	require.False(t, createdApp.ID.String() == "")
	require.False(t, createdApp.CreatedAt.String() == "")
	require.False(t, createdApp.UpdatedAt.String() == "")
}

func Test_AppService_Find(t *testing.T) {
	dbCloseCallbackMethod := prepareDB(t)
	defer dbCloseCallbackMethod()

	appService := models.AppService{DB: dataservices.GetDB()}
	testApp := createTestApp(t, &models.App{
		AppSlug: "test-app_slug",
	})

	foundApp, err := appService.Find(testApp)
	require.NoError(t, err)
	require.Equal(t, testApp, foundApp)
}

package models_test

import (
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
)

func createTestAppSettings(t *testing.T, appSettings *models.AppSettings) *models.AppSettings {
	err := dataservices.GetDB().Create(appSettings).Error
	require.NoError(t, err)
	return appSettings
}

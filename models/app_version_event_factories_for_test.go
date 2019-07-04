package models_test

import (
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
)

func createTestAppVersionEvent(t *testing.T, appVersionEvent *models.AppVersionEvent) *models.AppVersionEvent {
	err := dataservices.GetDB().Create(appVersionEvent).Error
	require.NoError(t, err)
	return appVersionEvent
}

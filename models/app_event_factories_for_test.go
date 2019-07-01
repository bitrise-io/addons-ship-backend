package models_test

import (
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
)

func createTestAppEvent(t *testing.T, appEvent *models.AppEvent) *models.AppEvent {
	err := dataservices.GetDB().Create(appEvent).Error
	require.NoError(t, err)
	return appEvent
}

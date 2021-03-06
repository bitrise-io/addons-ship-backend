// +build database

package models_test

import (
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
)

func createTestScreenshot(t *testing.T, screenshot *models.Screenshot) *models.Screenshot {
	err := dataservices.GetDB().Create(screenshot).Error
	require.NoError(t, err)
	return screenshot
}

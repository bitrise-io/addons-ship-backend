package models_test

import (
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
)

//nolint:unused,deadcode
func createTestPublishTask(t *testing.T, publishTask *models.PublishTask) *models.PublishTask {
	err := dataservices.GetDB().Create(publishTask).Error
	require.NoError(t, err)
	return publishTask
}

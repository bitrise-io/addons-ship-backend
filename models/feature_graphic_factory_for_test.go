package models_test

import (
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
)

//nolint:unused,deadcode
func createTestFeatureGraphic(t *testing.T, featureGraphic *models.FeatureGraphic) *models.FeatureGraphic {
	err := dataservices.GetDB().Create(featureGraphic).Error
	require.NoError(t, err)
	return featureGraphic
}

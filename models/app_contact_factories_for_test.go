package models_test

import (
	"testing"
	"time"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
)

func createTestAppContact(t *testing.T, appContact *models.AppContact) *models.AppContact {
	err := dataservices.GetDB().Create(appContact).Error
	require.NoError(t, err)
	return appContact
}

func compareAppContacts(t *testing.T, expected, actual models.AppContact) {
	expected.CreatedAt = time.Time{}
	expected.UpdatedAt = time.Time{}
	expected.App = nil
	actual.CreatedAt = time.Time{}
	actual.UpdatedAt = time.Time{}
	actual.App = nil
	require.Equal(t, expected, actual)
}

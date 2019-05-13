package factories

import (
	"testing"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/c2fo/testify/require"
)

// TestApp ...
func TestApp(t *testing.T, app *models.App) *models.App {
	err := dataservices.GetDB().Create(app).Error
	require.NoError(t, err)
	return app
}

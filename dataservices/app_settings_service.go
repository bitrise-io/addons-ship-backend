package dataservices

import "github.com/bitrise-io/addons-ship-backend/models"

// AppSettingsService ...
type AppSettingsService interface {
	Find(*models.AppSettings) (*models.AppSettings, error)
	Update(appSettings *models.AppSettings, whitelist []string) (validationErrors []error, dbErr error)
}

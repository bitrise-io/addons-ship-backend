package dataservices

import "github.com/bitrise-io/addons-ship-backend/models"

// AppVersionService ...
type AppVersionService interface {
	Create(*models.AppVersion) (*models.AppVersion, error)
	Find(*models.AppVersion) (*models.AppVersion, error)
	FindAll(app *models.App, filterParams map[string]interface{}) ([]models.AppVersion, error)
	Update(appVersion *models.AppVersion, whitelist []string) (validationErrors []error, dbErr error)
}

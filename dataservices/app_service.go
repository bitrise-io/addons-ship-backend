package dataservices

import "github.com/bitrise-io/addons-ship-backend/models"

// AppService ...
type AppService interface {
	Create(*models.App) (*models.App, error)
	Find(*models.App) (*models.App, error)
	Update(app *models.App, whitelist []string) (validationErrors []error, dbErr error)
	Delete(app *models.App) error
}

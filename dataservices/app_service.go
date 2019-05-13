package dataservices

import "github.com/bitrise-io/addons-ship-backend/models"

// AppService ...
type AppService interface {
	Create(*models.App) (*models.App, error)
	Find(*models.App) (*models.App, error)
}

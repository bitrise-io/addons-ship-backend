package dataservices

import (
	"github.com/bitrise-io/addons-ship-backend/models"
)

// AppEventService ...
type AppEventService interface {
	Create(appEvent *models.AppEvent) (*models.AppEvent, error)
	Find(appEvent *models.AppEvent) (*models.AppEvent, error)
	FindAll(app *models.App) ([]models.AppEvent, error)
}

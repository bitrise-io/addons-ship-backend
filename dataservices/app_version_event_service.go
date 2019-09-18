package dataservices

import (
	"github.com/bitrise-io/addons-ship-backend/models"
)

// AppVersionEventService ...
type AppVersionEventService interface {
	Create(appVersionEvent *models.AppVersionEvent) (*models.AppVersionEvent, error)
	Find(appVersionEvent *models.AppVersionEvent) (*models.AppVersionEvent, error)
	FindAll(appVersion *models.AppVersion) ([]models.AppVersionEvent, error)
	Update(appVersionEvent *models.AppVersionEvent, whitelist []string) (validationErrors []error, dbErr error)
}

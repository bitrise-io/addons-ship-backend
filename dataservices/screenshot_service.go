package dataservices

import "github.com/bitrise-io/addons-ship-backend/models"

// ScreenshotService ...
type ScreenshotService interface {
	Create(screenshot *models.Screenshot) (*models.Screenshot, []error, error)
	FindAll(appVersion *models.AppVersion) ([]models.Screenshot, error)
}

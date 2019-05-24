package dataservices

import "github.com/bitrise-io/addons-ship-backend/models"

// ScreenshotService ...
type ScreenshotService interface {
	FindAll(appVersion *models.AppVersion) ([]models.Screenshot, error)
}

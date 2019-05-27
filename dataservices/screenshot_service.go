package dataservices

import "github.com/bitrise-io/addons-ship-backend/models"

// ScreenshotService ...
type ScreenshotService interface {
	BatchCreate(screenshot []*models.Screenshot) ([]*models.Screenshot, []error, error)
	FindAll(appVersion *models.AppVersion) ([]models.Screenshot, error)
	BatchUpdate(screenshots []models.Screenshot, whitelist []string) ([]error, error)
}

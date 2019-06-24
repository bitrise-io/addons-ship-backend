package services_test

import "github.com/bitrise-io/addons-ship-backend/models"

type testAppSettingsService struct {
	findFn   func(*models.AppSettings) (*models.AppSettings, error)
	updateFn func(*models.AppSettings, []string) (validationErrors []error, dbErr error)
}

func (a *testAppSettingsService) Find(appSettings *models.AppSettings) (*models.AppSettings, error) {
	if a.findFn != nil {
		return a.findFn(appSettings)
	}
	panic("You have to override Find function in tests")
}

func (a *testAppSettingsService) Update(appSettings *models.AppSettings, whitelist []string) (validationErrors []error, dbErr error) {
	if a.updateFn != nil {
		return a.updateFn(appSettings, whitelist)
	}
	panic("You have to override Update function in tests")
}

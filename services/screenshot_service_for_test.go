package services_test

import "github.com/bitrise-io/addons-ship-backend/models"

type testScreenshotService struct {
	findAllFn func(*models.AppVersion) ([]models.Screenshot, error)
}

func (s *testScreenshotService) FindAll(appVersion *models.AppVersion) ([]models.Screenshot, error) {
	if s.findAllFn != nil {
		return s.findAllFn(appVersion)
	}
	panic("You have to override FindAll function in tests")
}

package services_test

import "github.com/bitrise-io/addons-ship-backend/models"

type testScreenshotService struct {
	createFn  func(*models.Screenshot) (*models.Screenshot, []error, error)
	findAllFn func(*models.AppVersion) ([]models.Screenshot, error)
}

func (s *testScreenshotService) Create(screenshot *models.Screenshot) (*models.Screenshot, []error, error) {
	if s.createFn != nil {
		return s.createFn(screenshot)
	}
	panic("You have to override Create function in tests")
}
func (s *testScreenshotService) FindAll(appVersion *models.AppVersion) ([]models.Screenshot, error) {
	if s.findAllFn != nil {
		return s.findAllFn(appVersion)
	}
	panic("You have to override FindAll function in tests")
}

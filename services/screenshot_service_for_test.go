package services_test

import "github.com/bitrise-io/addons-ship-backend/models"

type testScreenshotService struct {
	batchCreateFn func([]*models.Screenshot) ([]*models.Screenshot, []error, error)
	findAllFn     func(*models.AppVersion) ([]models.Screenshot, error)
	batchUpdateFn func([]models.Screenshot, []string) ([]error, error)
}

func (s *testScreenshotService) BatchCreate(screenshot []*models.Screenshot) ([]*models.Screenshot, []error, error) {
	if s.batchCreateFn != nil {
		return s.batchCreateFn(screenshot)
	}
	panic("You have to override BatchCreate function in tests")
}
func (s *testScreenshotService) FindAll(appVersion *models.AppVersion) ([]models.Screenshot, error) {
	if s.findAllFn != nil {
		return s.findAllFn(appVersion)
	}
	panic("You have to override FindAll function in tests")
}

func (s *testScreenshotService) BatchUpdate(screenshots []models.Screenshot, whitelist []string) ([]error, error) {
	if s.batchUpdateFn != nil {
		return s.batchUpdateFn(screenshots, whitelist)
	}
	panic("You have to override BatchUpdate function in tests")
}

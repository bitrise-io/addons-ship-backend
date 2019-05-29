package services_test

import "github.com/bitrise-io/addons-ship-backend/models"

type testScreenshotService struct {
	batchCreateFn func([]*models.Screenshot) ([]*models.Screenshot, []error, error)
	findFn        func(screenshot *models.Screenshot) (*models.Screenshot, error)
	findAllFn     func(*models.AppVersion) ([]models.Screenshot, error)
	batchUpdateFn func([]models.Screenshot, []string) ([]error, error)
	deleteFn      func(screenshot *models.Screenshot) (validationErrors []error, dbError error)
}

func (s *testScreenshotService) BatchCreate(screenshot []*models.Screenshot) ([]*models.Screenshot, []error, error) {
	if s.batchCreateFn != nil {
		return s.batchCreateFn(screenshot)
	}
	panic("You have to override BatchCreate function in tests")
}

func (s *testScreenshotService) Find(screenshot *models.Screenshot) (*models.Screenshot, error) {
	if s.findFn != nil {
		return s.findFn(screenshot)
	}
	panic("You have to override Find function in tests")
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

func (s *testScreenshotService) Delete(screenshot *models.Screenshot) (validationErrors []error, dbError error) {
	if s.deleteFn != nil {
		return s.deleteFn(screenshot)
	}
	panic("You have to override the Delete function in tests")
}

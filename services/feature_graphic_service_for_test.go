package services_test

import "github.com/bitrise-io/addons-ship-backend/models"

type testFeatureGraphicService struct {
	createFn func(*models.FeatureGraphic) (featureGraphic *models.FeatureGraphic, validationError []error, dbErr error)
	findFn   func(*models.FeatureGraphic) (*models.FeatureGraphic, error)
	updateFn func(models.FeatureGraphic, []string) (validationError []error, dbErr error)
	deleteFn func(screenshot *models.FeatureGraphic) error
}

func (s *testFeatureGraphicService) Create(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, []error, error) {
	if s.createFn != nil {
		return s.createFn(featureGraphic)
	}
	panic("You have to override Create function in tests")
}

func (s *testFeatureGraphicService) Find(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, error) {
	if s.findFn != nil {
		return s.findFn(featureGraphic)
	}
	panic("You have to override Find function in tests")
}

func (s *testFeatureGraphicService) Update(featureGraphic models.FeatureGraphic, whitelist []string) ([]error, error) {
	if s.updateFn != nil {
		return s.updateFn(featureGraphic, whitelist)
	}
	panic("You have to override BatchUpdate function in tests")
}

func (s *testFeatureGraphicService) Delete(featureGraphic *models.FeatureGraphic) error {
	if s.deleteFn != nil {
		return s.deleteFn(featureGraphic)
	}
	panic("You have to override the Delete function in tests")
}

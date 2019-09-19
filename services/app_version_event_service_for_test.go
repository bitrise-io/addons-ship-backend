package services_test

import "github.com/bitrise-io/addons-ship-backend/models"

type testAppVersionEventService struct {
	createFn  func(*models.AppVersionEvent) (*models.AppVersionEvent, error)
	findFn    func(*models.AppVersionEvent) (*models.AppVersionEvent, error)
	findAllFn func(*models.AppVersion) ([]models.AppVersionEvent, error)
	updateFn  func(*models.AppVersionEvent) ([]error, error)
}

func (a *testAppVersionEventService) Create(appVersionEvent *models.AppVersionEvent) (*models.AppVersionEvent, error) {
	if a.createFn != nil {
		return a.createFn(appVersionEvent)
	}
	panic("You have to override Create function in tests")
}

func (a *testAppVersionEventService) Find(appVersionEvent *models.AppVersionEvent) (*models.AppVersionEvent, error) {
	if a.findFn != nil {
		return a.findFn(appVersionEvent)
	}
	panic("You have to override Find function in tests")
}

func (a *testAppVersionEventService) FindAll(appVersion *models.AppVersion) ([]models.AppVersionEvent, error) {
	if a.findAllFn != nil {
		return a.findAllFn(appVersion)
	}
	panic("You have to override FindAll function in tests")
}

func (a *testAppVersionEventService) Update(appVersionEvent *models.AppVersionEvent, whitelist []string) (validationErrors []error, dbErr error) {
	if a.updateFn != nil {
		return a.updateFn(appVersionEvent)
	}
	panic("You have to override Update function in tests")
}

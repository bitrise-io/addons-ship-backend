package services_test

import "github.com/bitrise-io/addons-ship-backend/models"

type testAppVersionService struct {
	createFn  func(*models.AppVersion) (*models.AppVersion, []error, error)
	findFn    func(*models.AppVersion) (*models.AppVersion, error)
	findAllFn func(*models.App, map[string]interface{}) ([]models.AppVersion, error)
	updateFn  func(*models.AppVersion, []string) (validationErrors []error, dbErr error)
}

func (a *testAppVersionService) Create(appVersion *models.AppVersion) (*models.AppVersion, []error, error) {
	if a.createFn != nil {
		return a.createFn(appVersion)
	}
	panic("You have to override Create function in tests")
}
func (a *testAppVersionService) Find(appVersion *models.AppVersion) (*models.AppVersion, error) {
	if a.findFn != nil {
		return a.findFn(appVersion)
	}
	panic("You have to override Find function in tests")
}
func (a *testAppVersionService) FindAll(app *models.App, filterParams map[string]interface{}) ([]models.AppVersion, error) {
	if a.findAllFn != nil {
		return a.findAllFn(app, filterParams)
	}
	panic("You have to override FindAll function in tests")
}

func (a *testAppVersionService) Update(appVersion *models.AppVersion, whitelist []string) (validationErrors []error, dbErr error) {
	if a.updateFn != nil {
		return a.updateFn(appVersion, whitelist)
	}
	panic("You have to override Update function in tests")
}

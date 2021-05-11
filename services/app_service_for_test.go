package services_test

import "github.com/bitrise-io/addons-ship-backend/models"

type testAppService struct {
	createFn func(*models.App) (*models.App, error)
	findFn   func(*models.App) (*models.App, error)
	updateFn func(*models.App, []string) ([]error, error)
	deleteFn func(*models.App) error
}

func (a *testAppService) Create(app *models.App) (*models.App, error) {
	if a.createFn != nil {
		return a.createFn(app)
	}
	panic("You have to override Create function in tests")
}

func (a *testAppService) Find(app *models.App) (*models.App, error) {
	if a.findFn != nil {
		return a.findFn(app)
	}
	panic("You have to override Find function in tests")
}

func (a *testAppService) Update(app *models.App, whitelist []string) (validationErrors []error, dbErr error) {
	if a.updateFn != nil {
		return a.updateFn(app, whitelist)
	}
	panic("You have to override Update function in tests")
}

func (a *testAppService) Delete(app *models.App) error {
	if a.deleteFn != nil {
		return a.deleteFn(app)
	}
	panic("You have to override Delete function in tests")
}

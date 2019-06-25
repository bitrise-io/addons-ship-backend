package services_test

import "github.com/bitrise-io/addons-ship-backend/models"

type testAppService struct {
	createFn func(*models.App) (*models.App, error)
	findFn   func(*models.App) (*models.App, error)
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

func (a *testAppService) Delete(app *models.App) error {
	if a.deleteFn != nil {
		return a.deleteFn(app)
	}
	panic("You have to override Delete function in tests")
}

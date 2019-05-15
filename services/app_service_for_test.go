package services_test

import "github.com/bitrise-io/addons-ship-backend/models"

type testAppService struct {
	createFn func(*models.App) (*models.App, error)
	findFn   func(*models.App) (*models.App, error)
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

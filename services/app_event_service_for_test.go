package services_test

import "github.com/bitrise-io/addons-ship-backend/models"

type testAppEventService struct {
	createFn  func(*models.AppEvent) (*models.AppEvent, error)
	findFn    func(*models.AppEvent) (*models.AppEvent, error)
	findAllFn func(*models.App) ([]models.AppEvent, error)
}

func (a *testAppEventService) Create(appEvent *models.AppEvent) (*models.AppEvent, error) {
	if a.createFn != nil {
		return a.createFn(appEvent)
	}
	panic("You have to override Create function in tests")
}

func (a *testAppEventService) Find(appEvent *models.AppEvent) (*models.AppEvent, error) {
	if a.findFn != nil {
		return a.findFn(appEvent)
	}
	panic("You have to override Find function in tests")
}

func (a *testAppEventService) FindAll(app *models.App) ([]models.AppEvent, error) {
	if a.findAllFn != nil {
		return a.findAllFn(app)
	}
	panic("You have to override FindAll function in tests")
}

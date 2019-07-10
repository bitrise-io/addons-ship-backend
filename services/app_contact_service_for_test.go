package services_test

import "github.com/bitrise-io/addons-ship-backend/models"

type testAppContactService struct {
	createFn func(*models.AppContact) (*models.AppContact, error)
	findFn   func(*models.AppContact) (*models.AppContact, error)
	updateFn func(*models.AppContact, []string) error
	deleteFn func(*models.AppContact) error
}

func (a *testAppContactService) Create(appContact *models.AppContact) (*models.AppContact, error) {
	if a.createFn != nil {
		return a.createFn(appContact)
	}
	panic("You have to override AppContactService.Create function in tests")
}

func (a *testAppContactService) Find(appContact *models.AppContact) (*models.AppContact, error) {
	if a.findFn != nil {
		return a.findFn(appContact)
	}
	panic("You have to override AppContactService.Find function in tests")
}

func (a *testAppContactService) Update(appContact *models.AppContact, whitelist []string) error {
	if a.updateFn != nil {
		return a.updateFn(appContact, whitelist)
	}
	panic("You have to override AppContactService.Update function in tests")
}

func (a *testAppContactService) Delete(appContact *models.AppContact) error {
	if a.deleteFn != nil {
		return a.deleteFn(appContact)
	}
	panic("You have to override AppContactService.Delete function in tests")
}

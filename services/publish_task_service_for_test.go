package services_test

import "github.com/bitrise-io/addons-ship-backend/models"

type testPublishTaskService struct {
	createFn func(*models.PublishTask) (*models.PublishTask, error)
	findFn   func(*models.PublishTask) (*models.PublishTask, error)
}

func (a *testPublishTaskService) Create(publishTask *models.PublishTask) (*models.PublishTask, error) {
	if a.createFn != nil {
		return a.createFn(publishTask)
	}
	panic("You have to override Create function in tests")
}

func (a *testPublishTaskService) Find(publishTask *models.PublishTask) (*models.PublishTask, error) {
	if a.findFn != nil {
		return a.findFn(publishTask)
	}
	panic("You have to override Find function in tests")
}

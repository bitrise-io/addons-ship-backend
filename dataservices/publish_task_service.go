package dataservices

import (
	"github.com/bitrise-io/addons-ship-backend/models"
)

// PublishTaskService ...
type PublishTaskService interface {
	Create(publishTask *models.PublishTask) (*models.PublishTask, error)
	Find(publishTask *models.PublishTask) (*models.PublishTask, error)
}

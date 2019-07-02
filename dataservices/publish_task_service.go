package dataservices

// PublishTaskService ...
type PublishTaskService interface {
	Create(publishTask *models.PublishTask) (*models.PublishTask, error)
}

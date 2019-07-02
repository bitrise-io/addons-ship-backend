package models

import "github.com/jinzhu/gorm"

// PublishTaskService ...
type PublishTaskService struct {
	UpdatableModelService
	DB *gorm.DB
}

// Create ...
func (t *PublishTaskService) Create(publishTask *PublishTask) (*PublishTask, error) {
	return publishTask, t.DB.Create(&publishTask).Error
}

// Find ...
func (t *PublishTaskService) Find(publishTask *PublishTask) (*PublishTask, error) {
	err := t.DB.Where(publishTask).Preload("AppVersion").First(publishTask).Error
	if err != nil {
		return nil, err
	}
	return publishTask, nil
}

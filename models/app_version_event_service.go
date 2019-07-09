package models

import "github.com/jinzhu/gorm"

// AppVersionEventService ...
type AppVersionEventService struct {
	DB *gorm.DB
}

// Create ...
func (a *AppVersionEventService) Create(appVersionEvent *AppVersionEvent) (*AppVersionEvent, error) {
	result := a.DB.Create(appVersionEvent)
	if result.Error != nil {
		return nil, result.Error
	}

	return appVersionEvent, a.DB.Where(appVersionEvent).Preload("AppVersion").Preload("AppVersion.App").First(appVersionEvent).Error
}

// Find ...
func (a *AppVersionEventService) Find(appVersionEvent *AppVersionEvent) (*AppVersionEvent, error) {
	err := a.DB.Where(appVersionEvent).Preload("AppVersion").Preload("AppVersion.App").First(appVersionEvent).Error
	if err != nil {
		return nil, err
	}
	return appVersionEvent, nil
}

// FindAll ...
func (a *AppVersionEventService) FindAll(appVersion *AppVersion) ([]AppVersionEvent, error) {
	var appVersionEvents []AppVersionEvent
	err := a.DB.Preload("AppVersion").Preload("AppVersion.App").Where(map[string]interface{}{"app_version_id": appVersion.ID}).Find(&appVersionEvents).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return appVersionEvents, nil
}

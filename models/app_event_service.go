package models

import "github.com/jinzhu/gorm"

// AppEventService ...
type AppEventService struct {
	DB *gorm.DB
}

// Create ...
func (a *AppEventService) Create(appEvent *AppEvent) (*AppEvent, error) {
	return appEvent, a.DB.Create(&appEvent).Error
}

// Find ...
func (a *AppEventService) Find(appEvent *AppEvent) (*AppEvent, error) {
	err := a.DB.Where(appEvent).Preload("App").First(appEvent).Error
	if err != nil {
		return nil, err
	}
	return appEvent, nil
}

// FindAll ...
func (a *AppEventService) FindAll(app *App) ([]AppEvent, error) {
	var appEvents []AppEvent
	err := a.DB.Preload("App").Where(map[string]interface{}{"app_id": app.ID}).Find(&appEvents).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return appEvents, nil
}

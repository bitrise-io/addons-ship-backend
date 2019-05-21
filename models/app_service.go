package models

import "github.com/jinzhu/gorm"

// AppService ...
type AppService struct {
	DB *gorm.DB
}

// Create ...
func (a *AppService) Create(app *App) (*App, error) {
	return app, a.DB.Create(app).Error
}

// Find ...
func (a *AppService) Find(app *App) (*App, error) {
	err := a.DB.Preload("AppVersions").Where(app).First(app).Error
	if err != nil {
		return nil, err
	}
	return app, nil
}

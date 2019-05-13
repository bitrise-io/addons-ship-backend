package models

import "github.com/jinzhu/gorm"

// AppVersionService ...
type AppVersionService struct {
	DB *gorm.DB
}

// Create ...
func (a *AppVersionService) Create(app *AppVersion) (*AppVersion, error) {
	return app, a.DB.Create(app).Error
}

// Find ...
func (a *AppVersionService) Find(app *AppVersion) (*AppVersion, error) {
	err := a.DB.First(app).Error
	if err != nil {
		return nil, err
	}
	return app, nil
}

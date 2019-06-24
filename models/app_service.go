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

// Delete ...
func (a *AppService) Delete(app *App) error {
	result := a.DB.Delete(&app)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

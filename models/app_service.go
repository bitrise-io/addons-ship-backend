package models

import "github.com/jinzhu/gorm"

// AppService ...
type AppService struct {
	DB *gorm.DB
	UpdatableModelService
}

// Create ...
func (a *AppService) Create(app *App) (*App, error) {
	result := a.DB.Create(app)
	if result.Error != nil {
		return nil, result.Error
	}
	app.AppSettings.App = app
	return app, a.DB.Create(&app.AppSettings).Error
}

// Find ...
func (a *AppService) Find(app *App) (*App, error) {
	err := a.DB.Preload("AppVersions").Where(app).First(app).Error
	if err != nil {
		return nil, err
	}
	return app, nil
}

// Update ...
func (a *AppService) Update(app *App, whitelist []string) (validationErrors []error, dbErr error) {
	updateData, err := a.UpdateData(*app, whitelist)
	if err != nil {
		return nil, err
	}
	result := a.DB.Model(app).Updates(updateData)
	verrs := ValidationErrors(result.GetErrors())
	if len(verrs) > 0 {
		return verrs, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return nil, nil
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

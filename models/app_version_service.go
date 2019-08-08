package models

import "github.com/jinzhu/gorm"

// AppVersionService ...
type AppVersionService struct {
	DB *gorm.DB
	UpdatableModelService
}

// Create ...
func (a *AppVersionService) Create(appVersion *AppVersion) (*AppVersion, []error, error) {
	result := a.DB.Create(appVersion)
	verrs := ValidationErrors(result.GetErrors())
	if len(verrs) > 0 {
		return nil, verrs, nil
	}
	if result.Error != nil {
		return nil, nil, result.Error
	}
	return appVersion, nil, nil
}

// Find ...
func (a *AppVersionService) Find(appVersion *AppVersion) (*AppVersion, error) {
	err := a.DB.Preload("App").First(appVersion).Error
	if err != nil {
		return nil, err
	}
	return appVersion, nil
}

// FindAll ...
func (a *AppVersionService) FindAll(app *App, filterParams map[string]interface{}) ([]AppVersion, error) {
	var appVersions []AppVersion
	filterParams["app_id"] = app.ID
	err := a.DB.Where(filterParams).Find(&appVersions).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return appVersions, nil
}

// Update ...
func (a *AppVersionService) Update(appVersion *AppVersion, whitelist []string) (validationErrors []error, dbErr error) {
	updateData, err := a.UpdateData(*appVersion, whitelist)
	if err != nil {
		return nil, err
	}
	result := a.DB.Model(appVersion).Updates(updateData)
	verrs := ValidationErrors(result.GetErrors())
	if len(verrs) > 0 {
		return verrs, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return nil, nil
}

// Latest ...
func (a *AppVersionService) Latest(appVersion *AppVersion) (*AppVersion, error) {
	err := a.DB.Preload("App").Order("created_at DESC").First(appVersion, appVersion).Error
	if err != nil {
		return nil, err
	}
	return appVersion, nil
}

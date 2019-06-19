package models

import "github.com/jinzhu/gorm"

// AppSettingsService ...
type AppSettingsService struct {
	DB *gorm.DB
	UpdatableModelService
}

// Find ...
func (s *AppSettingsService) Find(appSettings *AppSettings) (*AppSettings, error) {
	err := s.DB.Where(appSettings).First(appSettings).Error
	if err != nil {
		return nil, err
	}

	return appSettings, nil
}

// Update ...
func (s *AppSettingsService) Update(appSettings *AppSettings, whitelist []string) (validationErrors []error, dbErr error) {
	updateData, err := s.UpdateData(*appSettings, whitelist)
	if err != nil {
		return nil, err
	}
	result := s.DB.Model(appSettings).Updates(updateData)
	verrs := ValidationErrors(result.GetErrors())
	if len(verrs) > 0 {
		return verrs, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return nil, nil
}

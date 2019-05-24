package models

import "github.com/jinzhu/gorm"

// ScreenshotService ...
type ScreenshotService struct {
	DB *gorm.DB
}

// Create ...
func (s *ScreenshotService) Create(screenshot *Screenshot) (*Screenshot, []error, error) {
	result := s.DB.Create(screenshot)
	verrs := result.GetErrors()
	if len(verrs) > 0 {
		return nil, verrs, nil
	}
	if result.Error != nil {
		return nil, nil, result.Error
	}
	return screenshot, nil, nil
}

// FindAll ...
func (s *ScreenshotService) FindAll(appVersion *AppVersion) ([]Screenshot, error) {
	var screenshots []Screenshot
	err := s.DB.Preload("AppVersion").Preload("AppVersion.App").Where(map[string]interface{}{"app_version_id": appVersion.ID}).
		Find(&screenshots).
		Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return screenshots, nil
}

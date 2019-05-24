package models

import "github.com/jinzhu/gorm"

// ScreenshotService ...
type ScreenshotService struct {
	DB *gorm.DB
}

// FindAll ...
func (a *ScreenshotService) FindAll(appVersion *AppVersion) ([]Screenshot, error) {
	var screenshots []Screenshot
	err := a.DB.Preload("AppVersion").Preload("AppVersion.App").Where(map[string]interface{}{"app_version_id": appVersion.ID}).
		Find(&screenshots).
		Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return screenshots, nil
}

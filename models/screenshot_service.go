package models

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// ScreenshotService ...
type ScreenshotService struct {
	DB *gorm.DB
	UpdatabeModelService
}

// BatchCreate ...
func (s *ScreenshotService) BatchCreate(screenshots []*Screenshot) ([]*Screenshot, []error, error) {
	tx := s.DB.Begin()
	for _, screenshot := range screenshots {
		result := tx.Create(screenshot)
		verrs := result.GetErrors()
		if len(verrs) > 0 {
			tx.Rollback()
			return nil, verrs, nil
		}
		if result.Error != nil {
			tx.Rollback()
			return nil, nil, result.Error
		}
	}
	return screenshots, nil, tx.Commit().Error
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

// BatchUpdate ...
func (s *ScreenshotService) BatchUpdate(screenshots []*Screenshot, whitelist []string) ([]*Screenshot, []error, error) {
	updatedupdatedScreenshot := []*Screenshot{}
	for _, screenshot := range screenshots {
		updatedScreenshotI, verrs, err := s.WhiteListedUpdate(s.DB, &screenshot, whitelist)
		updatedScreenshot, ok := updatedScreenshotI.(Screenshot)
		if !ok {
			return nil, nil, errors.New("Updated record has different type, than original")
		}
		if len(verrs) > 0 {
			return nil, verrs, nil
		}
		if err != nil {
			return nil, nil, err
		}
		updatedupdatedScreenshot = append(updatedupdatedScreenshot, &updatedScreenshot)
	}
	return updatedupdatedScreenshot, nil, nil
}

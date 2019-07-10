package models

import (
	"github.com/jinzhu/gorm"
)

// AppContactService ...
type AppContactService struct {
	DB *gorm.DB
	UpdatableModelService
}

// Create ...
func (s *AppContactService) Create(appContact *AppContact) (*AppContact, error) {
	result := s.DB.Create(appContact)
	if result.Error != nil {
		return nil, result.Error
	}
	return appContact, nil
}

// Find ...
func (s *AppContactService) Find(appContact *AppContact) (*AppContact, error) {
	err := s.DB.Preload("App").Where(appContact).First(appContact).Error
	if err != nil {
		return nil, err
	}

	return appContact, nil
}

// Update ...
func (s *AppContactService) Update(appContact *AppContact, whitelist []string) error {
	updateData, err := s.UpdateData(*appContact, whitelist)
	if err != nil {
		return err
	}
	result := s.DB.Model(appContact).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete ...
func (s *AppContactService) Delete(appContact *AppContact) error {
	result := s.DB.Delete(&appContact)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

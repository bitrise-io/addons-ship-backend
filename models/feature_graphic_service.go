package models

import "github.com/jinzhu/gorm"

// FeatureGraphicService ...
type FeatureGraphicService struct {
	DB *gorm.DB
	UpdatableModelService
}

// Create ...
func (s *FeatureGraphicService) Create(featureGraphic *FeatureGraphic) (*FeatureGraphic, []error, error) {
	result := s.DB.Create(featureGraphic)
	verrs := ValidationErrors(result.GetErrors())
	if len(verrs) > 0 {
		return nil, verrs, nil
	}
	if result.Error != nil {
		return nil, nil, result.Error
	}
	return featureGraphic, nil, nil
}

// Find ...
func (s *FeatureGraphicService) Find(featureGraphic *FeatureGraphic) (*FeatureGraphic, error) {
	err := s.DB.Where(featureGraphic).First(featureGraphic).Error
	if err != nil {
		return nil, err
	}

	return featureGraphic, nil
}

// Update ...
func (s *FeatureGraphicService) Update(featureGraphic FeatureGraphic, whitelist []string) ([]error, error) {
	updateData, err := s.UpdateData(featureGraphic, whitelist)
	if err != nil {
		return nil, err
	}
	result := s.DB.Model(&featureGraphic).Updates(updateData)
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
func (s *FeatureGraphicService) Delete(featureGraphic *FeatureGraphic) error {
	result := s.DB.Delete(&featureGraphic)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

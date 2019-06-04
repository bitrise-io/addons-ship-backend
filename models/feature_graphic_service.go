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

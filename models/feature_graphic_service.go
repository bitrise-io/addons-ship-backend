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

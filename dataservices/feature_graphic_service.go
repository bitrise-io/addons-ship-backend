package dataservices

import (
	"github.com/bitrise-io/addons-ship-backend/models"
)

// FeatureGraphicService ...
type FeatureGraphicService interface {
	Create(screenshot *models.FeatureGraphic) (*models.FeatureGraphic, []error, error)
	Find(screenshot *models.FeatureGraphic) (*models.FeatureGraphic, error)
	Update(screenshot models.FeatureGraphic, whitelist []string) (validationErrors []error, dbError error)
	Delete(screenshot *models.FeatureGraphic) error
}

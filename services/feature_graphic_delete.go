package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

// FeatureGraphicDeleteResponse ...
type FeatureGraphicDeleteResponse struct {
	Data *models.FeatureGraphic `json:"data"`
}

// FeatureGraphicDeleteHandler ...
func FeatureGraphicDeleteHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppVersionID, err := GetAuthorizedAppVersionIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	if env.FeatureGraphicService == nil {
		return errors.New("No Feature Graphic Service defined for handler")
	}

	if env.AWS == nil {
		return errors.New("No AWS Provider defined for handler")
	}

	featureGraphic, err := env.FeatureGraphicService.Find(
		&models.FeatureGraphic{AppVersionID: authorizedAppVersionID})
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	err = env.AWS.DeleteObject(featureGraphic.AWSPath())
	if err != nil {
		return errors.WithStack(err)
	}

	err = env.FeatureGraphicService.Delete(featureGraphic)
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	return httpresponse.RespondWithSuccess(w, FeatureGraphicDeleteResponse{
		Data: featureGraphic,
	})
}

package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// FeatureGraphicData ...
type FeatureGraphicData struct {
	models.FeatureGraphic
	DownloadURL string `json:"download_url,omitempty"`
	UploadURL   string `json:"upload_url,omitempty"`
}

// FeatureGraphicGetResponse ...
type FeatureGraphicGetResponse struct {
	Data FeatureGraphicData `json:"data"`
}

// FeatureGraphicGetHandler ...
func FeatureGraphicGetHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppVersionID, err := GetAuthorizedAppVersionIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if env.FeatureGraphicService == nil {
		return errors.New("No Feature Graphic Service defined for handler")
	}

	featureGraphic, err := env.FeatureGraphicService.Find(
		&models.FeatureGraphic{AppVersionID: authorizedAppVersionID},
	)
	switch {
	case errors.Cause(err) == gorm.ErrRecordNotFound:
		return httpresponse.RespondWithNotFoundError(w)
	case err != nil:
		return errors.Wrap(err, "SQL Error")
	}

	if env.AWS == nil {
		return errors.New("No AWS Provider defined for handler")
	}
	presignedURL, err := env.AWS.GeneratePresignedGETURL(featureGraphic.AWSPath(), presignedURLExpirationInterval)
	if err != nil {
		return errors.WithStack(err)
	}

	return httpresponse.RespondWithSuccess(w, FeatureGraphicGetResponse{
		Data: FeatureGraphicData{
			FeatureGraphic: *featureGraphic,
			DownloadURL:    presignedURL,
		},
	})
}

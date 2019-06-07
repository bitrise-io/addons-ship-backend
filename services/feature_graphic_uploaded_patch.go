package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

// FeatureGraphicUploadedPatchResponse ...
type FeatureGraphicUploadedPatchResponse struct {
	Data FeatureGraphicData `json:"data"`
}

// FeatureGraphicUploadedPatchHandler ...
func FeatureGraphicUploadedPatchHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppVersionID, err := GetAuthorizedAppVersionIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if env.FeatureGraphicService == nil {
		return errors.New("No Feature Graphic Service defined for handler")
	}

	featureGraphicToUpdate, err := env.FeatureGraphicService.Find(
		&models.FeatureGraphic{AppVersionID: authorizedAppVersionID},
	)
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	featureGraphicToUpdate.Uploaded = true
	verrs, err := env.FeatureGraphicService.Update(*featureGraphicToUpdate, []string{"Uploaded"})
	if len(verrs) > 0 {
		return httpresponse.RespondWithUnprocessableEntity(w, verrs)
	}
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	if env.AWS == nil {
		return errors.New("No AWS Provider defined for handler")
	}
	presignedURL, err := env.AWS.GeneratePresignedGETURL(featureGraphicToUpdate.AWSPath(), presignedURLExpirationInterval)
	if err != nil {
		return errors.WithStack(err)
	}
	return httpresponse.RespondWithSuccess(w, FeatureGraphicUploadedPatchResponse{
		Data: FeatureGraphicData{
			FeatureGraphic: *featureGraphicToUpdate,
			DownloadURL:    presignedURL,
		},
	})
}

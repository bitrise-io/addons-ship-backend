package services

import (
	"encoding/json"
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httprequest"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

type featureGraphicPostParamsElement struct {
	Filename string `json:"filename"`
	Filesize int64  `json:"filesize"`
}

// FeatureGraphicPostResponse ...
type FeatureGraphicPostResponse struct {
	Data FeatureGraphicData `json:"data"`
}

// FeatureGraphicPostHandler ...
func FeatureGraphicPostHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppVersionID, err := GetAuthorizedAppVersionIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	var params featureGraphicPostParamsElement
	defer httprequest.BodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}

	if env.FeatureGraphicService == nil {
		return errors.New("No Feature Graphic Service defined for handler")
	}

	createdFeatureGraphic, verrs, err := env.FeatureGraphicService.Create(&models.FeatureGraphic{
		UploadableObject: models.UploadableObject{
			Filename: params.Filename,
			Filesize: params.Filesize,
		},
		AppVersionID: authorizedAppVersionID,
	})
	if len(verrs) > 0 {
		return httpresponse.RespondWithUnprocessableEntity(w, verrs)
	}
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	if env.AWS == nil {
		return errors.New("No AWS Provider defined for handler")
	}
	presignedURL, err := env.AWS.GeneratePresignedPUTURL(createdFeatureGraphic.AWSPath(), presignedURLExpirationInterval, createdFeatureGraphic.Filesize)
	if err != nil {
		return errors.WithStack(err)
	}

	return httpresponse.RespondWithSuccess(w, FeatureGraphicPostResponse{
		Data: FeatureGraphicData{
			FeatureGraphic: *createdFeatureGraphic,
			UploadURL:      presignedURL,
		},
	})
}

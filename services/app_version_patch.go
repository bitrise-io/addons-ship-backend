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

// AppVersionPatchResponseData ...
type AppVersionPatchResponseData struct {
	*models.AppVersion
	AppStoreInfo models.AppStoreInfo `json:"app_store_info"`
}

// AppVersionPatchResponse ...
type AppVersionPatchResponse struct {
	Data AppVersionPatchResponseData `json:"data"`
}

// AppVersionPatchHandler ...
func AppVersionPatchHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppVersionID, err := GetAuthorizedAppVersionIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if env.AppVersionService == nil {
		return errors.New("No App Version Service defined for handler")
	}

	var params models.AppStoreInfo
	defer httprequest.BodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}
	appStoreInfo, err := json.Marshal(params)
	if err != nil {
		return errors.WithStack(err)
	}
	appVersionToUpdate := &models.AppVersion{
		Record:           models.Record{ID: authorizedAppVersionID},
		AppStoreInfoData: appStoreInfo,
	}
	verr, err := env.AppVersionService.Update(appVersionToUpdate, []string{"AppStoreInfoData"})
	if len(verr) > 0 {
		return httpresponse.RespondWithUnprocessableEntity(w, verr)
	}
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}
	response, err := newArtifactVersionPatchResponse(appVersionToUpdate)

	return httpresponse.RespondWithSuccess(w, AppVersionPatchResponse{
		Data: response,
	})
}

func newArtifactVersionPatchResponse(appVersion *models.AppVersion) (AppVersionPatchResponseData, error) {
	appStoreInfo, err := appVersion.AppStoreInfo()
	if err != nil {
		return AppVersionPatchResponseData{}, err
	}
	return AppVersionPatchResponseData{
		AppVersion:   appVersion,
		AppStoreInfo: appStoreInfo,
	}, nil
}

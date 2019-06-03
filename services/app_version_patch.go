package services

import (
	"encoding/json"
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httprequest"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// AppVersionPutResponseData ...
type AppVersionPutResponseData struct {
	*models.AppVersion
	AppStoreInfo models.AppStoreInfo `json:"app_store_info"`
}

// AppVersionPutResponse ...
type AppVersionPutResponse struct {
	Data AppVersionPutResponseData `json:"data"`
}

// AppVersionPutHandler ...
func AppVersionPutHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
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
	appVersionToUpdate, err := env.AppVersionService.Find(&models.AppVersion{Record: models.Record{ID: authorizedAppVersionID}})
	switch {
	case errors.Cause(err) == gorm.ErrRecordNotFound:
		return httpresponse.RespondWithNotFoundError(w)
	case err != nil:
		return errors.Wrap(err, "SQL Error")
	}
	appVersionToUpdate.AppStoreInfoData = appStoreInfo
	verr, err := env.AppVersionService.Update(appVersionToUpdate, []string{"AppStoreInfoData"})
	if len(verr) > 0 {
		return httpresponse.RespondWithUnprocessableEntity(w, verr)
	}
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}
	response, err := newArtifactVersionPatchResponse(appVersionToUpdate)
	if err != nil {
		return errors.WithStack(err)
	}

	return httpresponse.RespondWithSuccess(w, AppVersionPutResponse{
		Data: response,
	})
}

func newArtifactVersionPatchResponse(appVersion *models.AppVersion) (AppVersionPutResponseData, error) {
	appStoreInfo, err := appVersion.AppStoreInfo()
	if err != nil {
		return AppVersionPutResponseData{}, err
	}
	return AppVersionPutResponseData{
		AppVersion:   appVersion,
		AppStoreInfo: appStoreInfo,
	}, nil
}

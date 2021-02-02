package services

import (
	"encoding/json"
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httprequest"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// AppSettingsPatchParams ...
type AppSettingsPatchParams struct {
	IosSettings     models.IosSettings     `json:"ios_settings"`
	AndroidSettings models.AndroidSettings `json:"android_settings"`
	IosWorkflow     string                 `json:"ios_workflow"`
	AndroidWorkflow string                 `json:"android_workflow"`
}

// AppSettingsPatchResponseData ...
type AppSettingsPatchResponseData struct {
	*models.AppSettings
	IosSettings     models.IosSettings     `json:"ios_settings"`
	AndroidSettings models.AndroidSettings `json:"android_settings"`
}

// AppSettingsPatchResponse ...
type AppSettingsPatchResponse struct {
	Data AppSettingsPatchResponseData `json:"data"`
}

// AppSettingsPatchHandler ...
func AppSettingsPatchHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppID, err := GetAuthorizedAppIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if env.AppSettingsService == nil {
		return errors.New("No App Settings Service defined for handler")
	}

	var params AppSettingsPatchParams
	defer httprequest.BodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}

	appSettingsToUpdate, err := env.AppSettingsService.Find(&models.AppSettings{AppID: authorizedAppID})
	switch {
	case errors.Cause(err) == gorm.ErrRecordNotFound:
		return httpresponse.RespondWithNotFoundError(w)
	case err != nil:
		return errors.Wrap(err, "SQL Error")
	}

	appSettingsToUpdate, updateWhiteList, err := prepareAppSettingsToUpdate(env.BitriseAPI, appSettingsToUpdate, params)
	if err != nil {
		return errors.WithStack(err)
	}

	verr, err := env.AppSettingsService.Update(appSettingsToUpdate, updateWhiteList)
	if len(verr) > 0 {
		return httpresponse.RespondWithUnprocessableEntity(w, verr)
	}
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	response, err := newAppSettingsPatchResponse(appSettingsToUpdate)
	if err != nil {
		return errors.WithStack(err)
	}

	return httpresponse.RespondWithSuccess(w, AppSettingsPatchResponse{
		Data: response,
	})
}

func prepareAppSettingsToUpdate(api bitrise.APIInterface, appSettingsToUpdate *models.AppSettings, params AppSettingsPatchParams) (*models.AppSettings, []string, error) {
	updateWhiteList := []string{}
	if params.IosSettings.Valid() {
		if len(params.IosSettings.SelectedAppStoreProvisioningProfiles) > 0 {
			existingProvProfiles, err := api.GetProvisioningProfiles(appSettingsToUpdate.App.BitriseAPIToken, appSettingsToUpdate.App.AppSlug)
			if err != nil {
				return nil, []string{}, err
			}
			existingProvProfileSlugs := []string{}
			for _, provProfile := range existingProvProfiles {
				existingProvProfileSlugs = append(existingProvProfileSlugs, provProfile.Slug)
			}
			params.IosSettings.ValidateSelectedProvisioningProfileSlugs(existingProvProfileSlugs)
		}
		iosSettings, err := json.Marshal(params.IosSettings)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
		appSettingsToUpdate.IosSettingsData = iosSettings
		updateWhiteList = append(updateWhiteList, "IosSettingsData")
	}
	if params.AndroidSettings.Valid() {
		androidSettings, err := json.Marshal(params.AndroidSettings)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
		appSettingsToUpdate.AndroidSettingsData = androidSettings
		updateWhiteList = append(updateWhiteList, "AndroidSettingsData")
	}

	appSettingsToUpdate.IosWorkflow = params.IosWorkflow
	appSettingsToUpdate.AndroidWorkflow = params.AndroidWorkflow
	updateWhiteList = append(updateWhiteList, "IosWorkflow", "AndroidWorkflow")

	return appSettingsToUpdate, updateWhiteList, nil
}

func newAppSettingsPatchResponse(appSettings *models.AppSettings) (AppSettingsPatchResponseData, error) {
	iosSettings, err := appSettings.IosSettings()
	if err != nil {
		return AppSettingsPatchResponseData{}, err
	}
	androidSettings, err := appSettings.AndroidSettings()
	if err != nil {
		return AppSettingsPatchResponseData{}, err
	}
	return AppSettingsPatchResponseData{
		AppSettings:     appSettings,
		IosSettings:     iosSettings,
		AndroidSettings: androidSettings,
	}, nil
}

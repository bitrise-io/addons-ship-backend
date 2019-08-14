package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

// AppVersionsGetResponseElement ...
type AppVersionsGetResponseElement struct {
	models.AppVersion
	AppInfo              AppData  `json:"app_info"`
	DistributionType     string   `json:"distribution_type"`
	Version              string   `json:"version"`
	MinimumOS            string   `json:"minimum_os,omitempty"`
	MinimumSDK           string   `json:"minimum_sdk,omitempty"`
	Size                 int64    `json:"size"`
	SupportedDeviceTypes []string `json:"supported_device_types"`
}

// AppVersionsGetResponse ...
type AppVersionsGetResponse struct {
	Data []AppVersionsGetResponseElement `json:"data"`
}

// AppVersionsGetHandler ...
func AppVersionsGetHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppID, err := GetAuthorizedAppIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if env.AppVersionService == nil {
		return errors.New("No App Version Service defined for handler")
	}

	filterParams := map[string]interface{}{}
	if platformFilter := r.URL.Query().Get("platform"); platformFilter != "" {
		filterParams["platform"] = platformFilter
	}

	appVersions, err := env.AppVersionService.FindAll(
		&models.App{Record: models.Record{ID: authorizedAppID}},
		filterParams,
	)
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	if env.BitriseAPI == nil {
		return errors.New("No Bitrise API Service defined for handler")
	}

	response, err := newAppVersionsGetResponse(appVersions, env)
	if err != nil {
		return errors.WithStack(err)
	}

	return httpresponse.RespondWithSuccess(w, AppVersionsGetResponse{
		Data: response,
	})
}

func newAppVersionsGetResponse(appVersions []models.AppVersion, env *env.AppEnv) ([]AppVersionsGetResponseElement, error) {
	elements := []AppVersionsGetResponseElement{}

	var appData AppData
	if len(appVersions) > 0 {
		appDetails, err := env.BitriseAPI.GetAppDetails(appVersions[0].App.BitriseAPIToken, appVersions[0].App.AppSlug)
		if err != nil {
			return nil, err
		}
		appData = AppData{
			Title:       appDetails.Title,
			AppIconURL:  appDetails.AvatarURL,
			ProjectType: appDetails.ProjectType,
		}
	}

	for _, appVersion := range appVersions {
		artifactInfo, err := appVersion.ArtifactInfo()
		if err != nil {
			return nil, err
		}
		elements = append(elements, AppVersionsGetResponseElement{
			AppInfo:              appData,
			AppVersion:           appVersion,
			Version:              artifactInfo.Version,
			MinimumOS:            artifactInfo.MinimumOS,
			MinimumSDK:           artifactInfo.MinimumSDK,
			Size:                 artifactInfo.Size,
			DistributionType:     artifactInfo.DistributionType,
			SupportedDeviceTypes: artifactInfo.SupportedDeviceTypes,
		})
	}
	return elements, nil
}

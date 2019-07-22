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
	DistributionType     string   `json:"distributuin_type"`
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

	response, err := newAppVersionsGetResponse(appVersions)
	if err != nil {
		return errors.WithStack(err)
	}

	return httpresponse.RespondWithSuccess(w, AppVersionsGetResponse{
		Data: response,
	})
}

func newAppVersionsGetResponse(appVersions []models.AppVersion) ([]AppVersionsGetResponseElement, error) {
	elements := []AppVersionsGetResponseElement{}
	for _, appVersion := range appVersions {
		artifactInfo, err := appVersion.ArtifactInfo()
		if err != nil {
			return nil, err
		}
		elements = append(elements, AppVersionsGetResponseElement{
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

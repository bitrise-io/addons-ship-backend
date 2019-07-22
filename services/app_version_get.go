package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// AppData ...
type AppData struct {
	Title       string  `json:"title"`
	AppIconURL  *string `json:"app_icon_url"`
	ProjectType string  `json:"project_type"`
}

// AppVersionGetResponseData ...
type AppVersionGetResponseData struct {
	*models.AppVersion
	PublicInstallPageURL string              `json:"public_install_page_url"`
	AppStoreInfo         models.AppStoreInfo `json:"app_store_info"`
	PublishEnabled       bool                `json:"publish_enabled"`
	AppInfo              AppData             `json:"app_info"`
	DistributionType     string              `json:"distributuin_type,omitempty"`
	Version              string              `json:"version"`
	MinimumOS            string              `json:"minimum_os,omitempty"`
	MinimumSDK           string              `json:"minimum_sdk,omitempty"`
	BundleID             string              `json:"bundle_id,omitempty"`
	PackageName          string              `json:"package_name,omitempty"`
	Size                 int64               `json:"size"`
	SupportedDeviceTypes []string            `json:"supported_device_types"`
}

// AppVersionGetResponse ...
type AppVersionGetResponse struct {
	Data AppVersionGetResponseData `json:"data"`
}

// AppVersionGetHandler ...
func AppVersionGetHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppVersionID, err := GetAuthorizedAppVersionIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if env.AppVersionService == nil {
		return errors.New("No App Version Service defined for handler")
	}

	appVersion, err := env.AppVersionService.Find(
		&models.AppVersion{Record: models.Record{ID: authorizedAppVersionID}},
	)
	switch {
	case errors.Cause(err) == gorm.ErrRecordNotFound:
		return httpresponse.RespondWithNotFoundError(w)
	case err != nil:
		return errors.Wrap(err, "SQL Error")
	}

	if env.BitriseAPI == nil {
		return errors.New("No Bitrise API Service defined for handler")
	}

	artifacts, err := env.BitriseAPI.GetArtifacts(appVersion.App.APIToken, appVersion.App.AppSlug, appVersion.BuildSlug)
	if err != nil {
		return errors.WithStack(err)
	}

	responseData, err := newArtifactVersionGetResponse(appVersion, env, artifacts)
	if err != nil {
		return errors.WithStack(err)
	}

	return httpresponse.RespondWithSuccess(w, AppVersionGetResponse{
		Data: responseData,
	})
}

func newArtifactVersionGetResponse(appVersion *models.AppVersion, env *env.AppEnv, artifacts []bitrise.ArtifactListElementResponseModel) (AppVersionGetResponseData, error) {
	var publishEnabled, publicInstallPageEnabled bool
	var publicInstallPageArtifactSlug string
	switch appVersion.Platform {
	case "ios":
		_, publishEnabled, publicInstallPageEnabled, publicInstallPageArtifactSlug = selectIosArtifact(artifacts)
	case "android":
		_, publishEnabled, publicInstallPageEnabled, publicInstallPageArtifactSlug = selectAndroidArtifact(artifacts)
	default:
		return AppVersionGetResponseData{}, errors.Errorf("Invalid platform type of app version: %s", appVersion.Platform)
	}

	var artifactPublicInstallPageURL string
	if publicInstallPageEnabled {
		var err error
		artifactPublicInstallPageURL, err = env.BitriseAPI.GetArtifactPublicInstallPageURL(
			appVersion.App.BitriseAPIToken,
			appVersion.App.AppSlug,
			appVersion.BuildSlug,
			publicInstallPageArtifactSlug,
		)
		if err != nil {
			return AppVersionGetResponseData{}, errors.WithStack(err)
		}
	}

	appDetails, err := env.BitriseAPI.GetAppDetails(appVersion.App.APIToken, appVersion.App.AppSlug)
	if err != nil {
		return AppVersionGetResponseData{}, errors.WithStack(err)
	}
	appData := AppData{
		Title:       appDetails.Title,
		AppIconURL:  appDetails.AvatarURL,
		ProjectType: appDetails.ProjectType,
	}

	appStoreInfo, err := appVersion.AppStoreInfo()
	if err != nil {
		return AppVersionGetResponseData{}, err
	}
	artifactInfo, err := appVersion.ArtifactInfo()
	if err != nil {
		return AppVersionGetResponseData{}, err
	}
	return AppVersionGetResponseData{
		AppVersion:           appVersion,
		PublicInstallPageURL: artifactPublicInstallPageURL,
		AppStoreInfo:         appStoreInfo,
		PublishEnabled:       publishEnabled,
		AppInfo:              appData,
		DistributionType:     artifactInfo.DistributionType,
		Version:              artifactInfo.Version,
		MinimumOS:            artifactInfo.MinimumOS,
		MinimumSDK:           artifactInfo.MinimumSDK,
		Size:                 artifactInfo.Size,
		SupportedDeviceTypes: artifactInfo.SupportedDeviceTypes,
	}, nil
}

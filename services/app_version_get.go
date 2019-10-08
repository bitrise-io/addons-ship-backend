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
	Split                bool                `json:"split"`
	UniversalAvailable   bool                `json:"universal_available"`
	AppInfo              AppData             `json:"app_info"`
	IPAExportMethod      string              `json:"ipa_export_method,omitempty"`
	Version              string              `json:"version"`
	VersionCode          string              `json:"version_code"`
	MinimumOS            string              `json:"minimum_os,omitempty"`
	MinimumSDK           string              `json:"minimum_sdk,omitempty"`
	BundleID             string              `json:"bundle_id,omitempty"`
	PackageName          string              `json:"package_name,omitempty"`
	SupportedDeviceTypes []string            `json:"supported_device_types"`
	Module               string              `json:"module"`
	ProductFlavour       string              `json:"product_flavour"`
	BuildType            string              `json:"build_type"`
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

	artifacts, err := env.BitriseAPI.GetArtifacts(appVersion.App.BitriseAPIToken, appVersion.App.AppSlug, appVersion.BuildSlug)
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
	var ipaExportMethod string
	var publicInstallPageArtifactSlug string
	var publishAndShareInfo bitrise.PublishAndShareInfo
	switch appVersion.Platform {
	case "ios":
		_, publishEnabled, publicInstallPageEnabled, ipaExportMethod, publicInstallPageArtifactSlug = selectIosArtifact(artifacts)
	case "android":
		var err error
		artifactSelector := bitrise.NewArtifactSelector(artifacts)
		publishAndShareInfo, err = artifactSelector.PublishAndShareInfo(appVersion)
		if err != nil {
			return AppVersionGetResponseData{}, errors.WithStack(err)
		}
		publishEnabled = publishAndShareInfo.PublishEnabled
		publicInstallPageEnabled = publishAndShareInfo.PublicInstallPageEnabled
		publicInstallPageArtifactSlug = publishAndShareInfo.PublicInstallPageArtifactSlug
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

	appDetails, err := env.BitriseAPI.GetAppDetails(appVersion.App.BitriseAPIToken, appVersion.App.AppSlug)
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
		Split:                publishAndShareInfo.Split,
		UniversalAvailable:   publishAndShareInfo.UniversalAvailable,
		AppInfo:              appData,
		IPAExportMethod:      ipaExportMethod,
		Version:              artifactInfo.Version,
		MinimumOS:            artifactInfo.MinimumOS,
		MinimumSDK:           artifactInfo.MinimumSDK,
		SupportedDeviceTypes: artifactInfo.SupportedDeviceTypes,
		BundleID:             artifactInfo.BundleID,
		PackageName:          artifactInfo.PackageName,
		VersionCode:          artifactInfo.VersionCode,
		Module:               artifactInfo.Module,
		ProductFlavour:       appVersion.ProductFlavour,
		BuildType:            artifactInfo.BuildType,
	}, nil
}

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

	artifacts, err := env.BitriseAPI.GetArtifacts(
		appVersion.App.BitriseAPIToken,
		appVersion.App.AppSlug,
		appVersion.BuildSlug,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	switch appVersion.Platform {
	case "ios":
		return appVersionGetIosHelper(env, w, r, appVersion, artifacts)
	case "android":
		return appVersionGetAndroidHelper(env, w, r, appVersion, artifacts)
	default:
		return errors.Errorf("Invalid platform type of app version: %s", appVersion.Platform)
	}
}

func newArtifactVersionGetResponse(appVersion *models.AppVersion, artifact bitrise.ArtifactListElementResponseModel,
	publicInstallPageURL string, appDetails *bitrise.AppDetails, publishEnabled bool) (AppVersionGetResponseData, error) {
	var supportedDeviceTypes []string
	artifactMeta := artifact.ArtifactMeta
	if artifactMeta == nil {
		return AppVersionGetResponseData{}, errors.New("No artifact meta data found for artifact")
	}
	for _, familyID := range artifactMeta.AppInfo.DeviceFamilyList {
		switch familyID {
		case 1:
			supportedDeviceTypes = append(supportedDeviceTypes, "iPhone", "iPod Touch")
		case 2:
			supportedDeviceTypes = append(supportedDeviceTypes, "iPad")
		default:
			supportedDeviceTypes = append(supportedDeviceTypes, "Unknown")
		}
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
	var size int64
	if artifact.FileSizeBytes != nil {
		size = *artifact.FileSizeBytes
	}
	return AppVersionGetResponseData{
		AppVersion:           appVersion,
		PublicInstallPageURL: publicInstallPageURL,
		AppStoreInfo:         appStoreInfo,
		PublishEnabled:       publishEnabled,
	}, nil
}

package services

import (
	"net/http"
	"strconv"
	"time"

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
	MinimumOS            string              `json:"minimum_os"`
	MinimumSDK           string              `json:"minimum_sdk"`
	PackageName          string              `json:"package_name"`
	CertificateExpires   time.Time           `json:"certificate_expires"`
	BundleID             string              `json:"bundle_id"`
	Size                 int64               `json:"size"`
	SupportedDeviceTypes []string            `json:"supported_device_types"`
	DistributionType     string              `json:"distribution_type"`
	PublicInstallPageURL string              `json:"public_install_page_url"`
	AppInfo              AppData             `json:"app_info"`
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

	if appVersion.Platform == "ios" {
		return appVersionGetIosHelper(env, w, r, appVersion, artifacts)
	} else if appVersion.Platform == "android" {
		return appVersionGetAndroidHelper(env, w, r, appVersion, artifacts)
	}
	return errors.Errorf("Invalid platform type of app version: %s", appVersion.Platform)
}

func newArtifactVersionGetResponse(appVersion *models.AppVersion, artifact *bitrise.ArtifactMeta, publicInstallPageURL string, appDetails *bitrise.AppDetails) (AppVersionGetResponseData, error) {
	var supportedDeviceTypes []string
	for _, familyID := range artifact.AppInfo.DeviceFamilyList {
		switch familyID {
		case 1:
			supportedDeviceTypes = append(supportedDeviceTypes, "iPhone", "iPod Touch")
		case 2:
			supportedDeviceTypes = append(supportedDeviceTypes, "iPad")
		default:
			supportedDeviceTypes = append(supportedDeviceTypes, "Unknown")
		}
	}
	var size int64
	if artifact.Size != "" {
		var err error
		floatSize, err := strconv.ParseFloat(artifact.Size, 64)
		if err != nil {
			return AppVersionGetResponseData{}, err
		}
		size = int64(floatSize)
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
	return AppVersionGetResponseData{
		AppVersion:           appVersion,
		MinimumOS:            artifact.AppInfo.MinimumOS,
		MinimumSDK:           artifact.AppInfo.MinimumSDKVersion,
		PackageName:          artifact.AppInfo.PackageName,
		CertificateExpires:   artifact.ProvisioningInfo.ExpireDate,
		BundleID:             artifact.AppInfo.BundleID,
		Size:                 size,
		SupportedDeviceTypes: supportedDeviceTypes,
		DistributionType:     artifact.ProvisioningInfo.DistributionType,
		PublicInstallPageURL: publicInstallPageURL,
		AppInfo:              appData,
		AppStoreInfo:         appStoreInfo,
	}, nil
}

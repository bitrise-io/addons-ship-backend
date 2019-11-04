package services

import (
	"net/http"
	"strings"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

// AppVersionIosConfigGetResponse ...
type AppVersionIosConfigGetResponse struct {
	MetaData  IosConfigMetaData `json:"meta_data"`
	Artifacts []string          `json:"artifacts"`
}

// AppVersionIosConfigGetHandler ...
func AppVersionIosConfigGetHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppVersionID, err := GetAuthorizedAppVersionIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	if env.AppVersionService == nil {
		return errors.New("No App Version Service defined for handler")
	}
	if env.AppSettingsService == nil {
		return errors.New("No App Settings Service defined for handler")
	}
	if env.AWS == nil {
		return errors.New("No AWS Provider defined for handler")
	}
	if env.BitriseAPI == nil {
		return errors.New("No Bitrise API Service defined for handler")
	}
	if env.ScreenshotService == nil {
		return errors.New("No Screenshot Service defined for handler")
	}

	config := AppVersionIosConfigGetResponse{MetaData: IosConfigMetaData{}}

	appVersion, err := env.AppVersionService.Find(&models.AppVersion{Record: models.Record{ID: authorizedAppVersionID}})
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	storeInfo, err := appVersion.AppStoreInfo()
	if err != nil {
		return errors.WithStack(err)
	}
	screenshots, err := env.ScreenshotService.FindAll(appVersion)
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}
	scs, err := newIosScreenshotsResponse(screenshots, env)
	if err != nil {
		return errors.WithStack(err)
	}
	listingInfo := IosListingInfo{
		Screenshots:     scs,
		Description:     storeInfo.FullDescription,
		PromotionalText: storeInfo.PromotionalText,
		SupportURL:      storeInfo.SupportURL,
		SoftwareURL:     storeInfo.MarketingURL,
	}
	if len(storeInfo.Keywords) > 0 {
		listingInfo.Keywords = strings.Split(storeInfo.Keywords, ",")
	}
	config.MetaData.ListingInfoMap = map[string]IosListingInfo{"en-US": listingInfo}

	appSettings, err := env.AppSettingsService.Find(&models.AppSettings{AppID: appVersion.AppID})
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	iosSettings, err := appSettings.IosSettings()
	if err != nil {
		return errors.WithStack(err)
	}

	for _, provProfileSlug := range iosSettings.SelectedAppStoreProvisioningProfile {
		selectedProvisioningProfile, err := env.BitriseAPI.GetProvisioningProfile(appVersion.App.BitriseAPIToken, appVersion.App.AppSlug, provProfileSlug)
		if err != nil {
			return errors.WithStack(err)
		}
		config.MetaData.Signing.AppStoreProfileURL = append(config.MetaData.Signing.AppStoreProfileURL, selectedProvisioningProfile.DownloadURL)
	}

	codeSigningID, err := env.BitriseAPI.GetCodeSigningIdentity(appVersion.App.BitriseAPIToken, appVersion.App.AppSlug, iosSettings.SelectedCodeSigningIdentity)
	if err != nil {
		return errors.WithStack(err)
	}
	config.MetaData.Signing.DistributionCertificateURL = codeSigningID.DownloadURL
	config.MetaData.Signing.DistributionCertificatePasshprase = codeSigningID.CertificatePassword

	config.MetaData.ExportOptions = ExportOptions{IncludeBitcode: iosSettings.IncludeBitCode}
	config.MetaData.SKU = iosSettings.AppSKU
	config.MetaData.AppleUser = iosSettings.AppleDeveloperAccountEmail
	config.MetaData.AppleAppSpecificPassword = iosSettings.ApplSpecificPassword

	artifacts, err := env.BitriseAPI.GetArtifacts(appVersion.App.BitriseAPIToken, appVersion.App.AppSlug, appVersion.BuildSlug)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, artifact := range artifacts {
		if artifact.IsXCodeArchive() {
			artifactData, err := env.BitriseAPI.GetArtifact(appVersion.App.BitriseAPIToken, appVersion.App.AppSlug, appVersion.BuildSlug, artifact.Slug)
			if err != nil {
				return errors.WithStack(err)
			}
			if artifactData.DownloadPath == nil {
				return errors.New("Failed to get download URL for artifact")
			}
			config.Artifacts = append(config.Artifacts, *artifactData.DownloadPath)
		}
	}

	return httpresponse.RespondWithSuccess(w, config)
}

func newIosScreenshotsResponse(screenshotData []models.Screenshot, env *env.AppEnv) (map[string][]string, error) {
	scs := map[string][]string{}
	for _, sc := range screenshotData {
		url, err := env.AWS.GeneratePresignedGETURL(sc.AWSPath(), presignedURLExpirationInterval)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		scs[sc.ScreenSize] = append(scs[sc.ScreenSize], url)
	}
	return scs, nil
}

// IosListingInfo ...
type IosListingInfo struct {
	Screenshots     map[string][]string `json:"screenshots" yaml:"screenshots"`
	Description     string              `json:"description" yaml:"description"`
	PromotionalText string              `json:"promotional_text" yaml:"promotional_text"`
	Keywords        []string            `json:"keywords" yaml:"keywords"`
	SupportURL      string              `json:"support_url" yaml:"support_url"`
	SoftwareURL     string              `json:"software_url" yaml:"software_url"`
}

// Signing ...
type Signing struct {
	DistributionCertificatePasshprase string   `json:"distribution_certificate_passhprase"`
	DistributionCertificateURL        string   `json:"distribution_certificate_url"`
	AppStoreProfileURL                []string `json:"app_store_profile_url"`
}

// ExportOptions ...
type ExportOptions struct {
	IncludeBitcode bool `json:"include_bitcode,string"`
}

// IosConfigMetaData ...
type IosConfigMetaData struct {
	ListingInfoMap           map[string]IosListingInfo `json:"listing_info"`
	Signing                  Signing                   `json:"signing"`
	ExportOptions            ExportOptions             `json:"export_options"`
	SKU                      string                    `json:"sku"`
	AppleUser                string                    `json:"apple_user"`
	AppleAppSpecificPassword string                    `json:"apple_app_specific_password"`
}

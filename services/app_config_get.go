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

// AppVersionConfigGetResponse ...
type AppVersionConfigGetResponse struct {
	MetaData  metaData `json:"meta_data"`
	Artifacts []string `json:"artifacts"`
}

// AppVersionConfigGetHandler ...
func AppVersionConfigGetHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppVersionID, err := GetAuthorizedAppVersionIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	config := AppVersionConfigGetResponse{MetaData: metaData{}}

	appVersion, err := env.AppVersionService.Find(&models.AppVersion{Record: models.Record{ID: authorizedAppVersionID}})
	if err != nil {
		return errors.WithStack(err)
	}

	artifactInfo, err := appVersion.ArtifactInfo()
	if err != nil {
		return errors.WithStack(err)
	}
	config.MetaData.PackageName = artifactInfo.PackageName

	featureGraphic, err := env.FeatureGraphicService.Find(&models.FeatureGraphic{AppVersionID: authorizedAppVersionID})
	switch {
	case errors.Cause(err) == gorm.ErrRecordNotFound:
		return httpresponse.RespondWithNotFoundError(w)
	case err != nil:
		return errors.WithStack(err)
	}

	featureGraphicPresignedURL, err := env.AWS.GeneratePresignedGETURL(featureGraphic.AWSPath(), presignedURLExpirationInterval)
	if err != nil {
		return errors.WithStack(err)
	}

	storeInfo, err := appVersion.AppStoreInfo()
	if err != nil {
		return errors.WithStack(err)
	}
	appData, err := env.BitriseAPI.GetAppDetails(appVersion.App.APIToken, appVersion.App.AppSlug)
	if err != nil {
		return errors.WithStack(err)
	}
	config.MetaData.ListingInfo = listingInfo{
		ShortDescription: storeInfo.ShortDescription,
		FullDescription:  storeInfo.FullDescription,
		WhatsNew:         storeInfo.WhatsNew,
		FeatureGraphic:   featureGraphicPresignedURL,
		Title:            appData.Title,
	}

	appSettings, err := env.AppSettingsService.Find(&models.AppSettings{AppID: appVersion.AppID})
	if err != nil {
		return errors.WithStack(err)
	}

	androidSettings, err := appSettings.AndroidSettings()
	if err != nil {
		return errors.WithStack(err)
	}
	config.MetaData.Track = androidSettings.Track

	var selectedServiceAccount bitrise.GenericProjectFile
	serviceAccounts, err := env.BitriseAPI.GetServiceAccountFiles(appVersion.App.APIToken, appVersion.App.AppSlug)
	for _, serviceAccount := range serviceAccounts {
		if serviceAccount.Slug == androidSettings.SelectedServiceAccount {
			selectedServiceAccount = serviceAccount
			break
		}
	}
	if selectedServiceAccount == (bitrise.GenericProjectFile{}) {
		return httpresponse.RespondWithNotFoundError(w)
	}

	var selectedAndroidKeystore bitrise.AndroidKeystoreFile
	androidKeystoreFiles, err := env.BitriseAPI.GetAndroidKeystoreFiles(appVersion.App.APIToken, appVersion.App.AppSlug)
	for _, keystore := range androidKeystoreFiles {
		if keystore.Slug == androidSettings.SelectedKeystoreFile && keystore.UserEnvKey == "ANDROID_KEYSTORE" {
			selectedAndroidKeystore = keystore
			break
		}
	}
	if selectedAndroidKeystore == (bitrise.AndroidKeystoreFile{}) {
		return httpresponse.RespondWithNotFoundError(w)
	}
	config.MetaData.ServiceAccountJSON = selectedAndroidKeystore.DownloadURL

	ks := keystore{
		URL:         selectedAndroidKeystore.DownloadURL,
		Password:    selectedAndroidKeystore.ExposedMetadataStore.Password,
		Alias:       selectedAndroidKeystore.ExposedMetadataStore.Alias,
		KeyPassword: selectedAndroidKeystore.ExposedMetadataStore.PrivateKeyPassword,
	}

	config.MetaData.Keystore = ks

	// TODO: screenshots

	return httpresponse.RespondWithSuccess(w, AppVersionConfigGetResponse{})
}

type screenshots struct {
	Tv        []string `json:"tv,omitempty"`
	Wear      []string `json:"wear,omitempty"`
	Phone     []string `json:"phone,omitempty"`
	TenInch   []string `json:"ten_inch,omitempty"`
	SevenInch []string `json:"seven_inch,omitempty"`
}

type listingInfo struct {
	Screenshots      screenshots `json:"screenshots,omitempty"`
	Icon             string      `json:"icon,omitempty"`
	Video            string      `json:"video,omitempty"`
	Title            string      `json:"title,omitempty"`
	TvBanner         string      `json:"tv_banner,omitempty"`
	WhatsNew         string      `json:"whats_new,omitempty"`
	PromoGraphic     string      `json:"promo_graphic,omitempty"`
	FeatureGraphic   string      `json:"feature_graphic,omitempty"`
	FullDescription  string      `json:"full_description,omitempty"`
	ShortDescription string      `json:"short_description,omitempty"`
}

type keystore struct {
	URL         string `json:"url"`
	Password    string `json:"password"`
	Alias       string `json:"alias"`
	KeyPassword string `json:"key_password"`
}

type metaData struct {
	ListingInfo        listingInfo `json:"listing_info"`
	Track              string      `json:"track"`
	PackageName        string      `json:"package_name"`
	ServiceAccountJSON string      `json:"service_account_json"`
	Keystore           keystore    `json:"keystore"`
}

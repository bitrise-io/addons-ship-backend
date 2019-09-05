package services

import (
	"net/http"
	"reflect"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// AppVersionAndroidConfigGetResponse ...
type AppVersionAndroidConfigGetResponse struct {
	MetaData  MetaData `json:"meta_data"`
	Artifacts []string `json:"artifacts"`
}

// AppVersionAndroidConfigGetHandler ...
func AppVersionAndroidConfigGetHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
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
	if env.FeatureGraphicService == nil {
		return errors.New("No Feature Graphic Service defined for handler")
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

	config := AppVersionAndroidConfigGetResponse{MetaData: MetaData{}}

	appVersion, err := env.AppVersionService.Find(&models.AppVersion{Record: models.Record{ID: authorizedAppVersionID}})
	if err != nil {
		return errors.Wrap(err, "SQL Error")
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
		return errors.Wrap(err, "SQL Error")
	}

	featureGraphicPresignedURL, err := env.AWS.GeneratePresignedGETURL(featureGraphic.AWSPath(), presignedURLExpirationInterval)
	if err != nil {
		return errors.WithStack(err)
	}

	storeInfo, err := appVersion.AppStoreInfo()
	if err != nil {
		return errors.WithStack(err)
	}
	appData, err := env.BitriseAPI.GetAppDetails(appVersion.App.BitriseAPIToken, appVersion.App.AppSlug)
	if err != nil {
		return errors.WithStack(err)
	}

	screenshots, err := env.ScreenshotService.FindAll(appVersion)
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	scs, err := newScreenshotsResponse(screenshots, env)
	if err != nil {
		return errors.WithStack(err)
	}
	config.MetaData.ListingInfo = ListingInfos{
		"en_US": ListingInfo{
			ShortDescription: storeInfo.ShortDescription,
			FullDescription:  storeInfo.FullDescription,
			WhatsNew:         storeInfo.WhatsNew,
			FeatureGraphic:   featureGraphicPresignedURL,
			Title:            appData.Title,
			Screenshots:      scs,
		},
	}

	appSettings, err := env.AppSettingsService.Find(&models.AppSettings{AppID: appVersion.AppID})
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	androidSettings, err := appSettings.AndroidSettings()
	if err != nil {
		return errors.WithStack(err)
	}
	config.MetaData.Track = androidSettings.Track

	selectedServiceAccount, err := env.BitriseAPI.GetServiceAccountFile(appVersion.App.BitriseAPIToken, appVersion.App.AppSlug, androidSettings.SelectedServiceAccount)
	if err != nil {
		return errors.WithStack(err)
	}
	config.MetaData.ServiceAccountJSON = selectedServiceAccount.DownloadURL

	selectedAndroidKeystore, err := env.BitriseAPI.GetAndroidKeystoreFile(appVersion.App.BitriseAPIToken, appVersion.App.AppSlug, androidSettings.SelectedKeystoreFile)
	if err != nil {
		return errors.WithStack(err)
	}

	config.MetaData.Keystore = Keystore{
		URL:         selectedAndroidKeystore.DownloadURL,
		Password:    selectedAndroidKeystore.ExposedMetadataStore.Password,
		Alias:       selectedAndroidKeystore.ExposedMetadataStore.Alias,
		KeyPassword: selectedAndroidKeystore.ExposedMetadataStore.PrivateKeyPassword,
	}

	artifacts, err := env.BitriseAPI.GetArtifacts(appVersion.App.BitriseAPIToken, appVersion.App.AppSlug, appVersion.BuildSlug)
	if err != nil {
		return errors.WithStack(err)
	}

	artifactList, err := newArtifactResponse(env, appVersion.App.BitriseAPIToken, appVersion.App.AppSlug, appVersion.BuildSlug, artifacts)
	if err != nil {
		return errors.WithStack(err)
	}
	config.Artifacts = artifactList

	return httpresponse.RespondWithSuccess(w, config)
}

func newScreenshotsResponse(screenshotData []models.Screenshot, env *env.AppEnv) (Screenshots, error) {
	scs := Screenshots{}
	for _, sc := range screenshotData {
		url, err := env.AWS.GeneratePresignedGETURL(sc.AWSPath(), presignedURLExpirationInterval)
		if err != nil {
			return Screenshots{}, errors.WithStack(err)
		}
		switch sc.ScreenSize {
		case "tv":
			scs.Tv = append(scs.Tv, url)
		case "wear":
			scs.Wear = append(scs.Wear, url)
		case "phone":
			scs.Phone = append(scs.Phone, url)
		case "ten_inch":
			scs.TenInch = append(scs.TenInch, url)
		case "seven_inch":
			scs.SevenInch = append(scs.SevenInch, url)
		}
	}
	return scs, nil
}

func newArtifactResponse(env *env.AppEnv, apiToken, appSlug, buildSlug string, artifacts []bitrise.ArtifactListElementResponseModel) ([]string, error) {
	artifactURLs := []string{}
	selectedArtifact, _, _, _ := selectAndroidArtifact(artifacts)
	if selectedArtifact != nil && !reflect.DeepEqual(*selectedArtifact, bitrise.ArtifactListElementResponseModel{}) {
		artifactData, err := env.BitriseAPI.GetArtifact(apiToken, appSlug, buildSlug, selectedArtifact.Slug)
		if err != nil {
			return []string{}, errors.WithStack(err)
		}
		if artifactData.DownloadPath == nil {
			return []string{}, errors.New("Failed to get download URL for artifact")
		}
		artifactURLs = append(artifactURLs, *artifactData.DownloadPath)
		return artifactURLs, nil
	}
	splitAPKs := checkForSplitAPKs(artifacts)
	if len(splitAPKs) == 0 {
		return []string{}, nil
	}
	for _, artifact := range splitAPKs {
		artifactData, err := env.BitriseAPI.GetArtifact(apiToken, appSlug, buildSlug, artifact.Slug)
		if err != nil {
			return []string{}, errors.WithStack(err)
		}
		if artifactData.DownloadPath == nil {
			return []string{}, errors.New("Failed to get download URL for artifact")
		}
		artifactURLs = append(artifactURLs, *artifactData.DownloadPath)
	}
	return artifactURLs, nil
}

// Screenshots ...
type Screenshots struct {
	Tv        []string `json:"tv,omitempty"`
	Wear      []string `json:"wear,omitempty"`
	Phone     []string `json:"phone,omitempty"`
	TenInch   []string `json:"ten_inch,omitempty"`
	SevenInch []string `json:"seven_inch,omitempty"`
}

// ListingInfo ...
type ListingInfo struct {
	Screenshots      Screenshots `json:"screenshots,omitempty"`
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

// ListingInfos ...
type ListingInfos map[string]ListingInfo

// Keystore ...
type Keystore struct {
	URL         string `json:"url"`
	Password    string `json:"password"`
	Alias       string `json:"alias"`
	KeyPassword string `json:"key_password"`
}

// MetaData ...
type MetaData struct {
	ListingInfo        ListingInfos `json:"listing_info"`
	Track              string       `json:"track"`
	PackageName        string       `json:"package_name"`
	ServiceAccountJSON string       `json:"service_account_json"`
	Keystore           Keystore     `json:"keystore"`
}

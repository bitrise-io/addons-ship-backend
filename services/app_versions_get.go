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
	IPAExportMethod      string   `json:"ipa_export_method"`
	Version              string   `json:"version"`
	MinimumOS            string   `json:"minimum_os,omitempty"`
	MinimumSDK           string   `json:"minimum_sdk,omitempty"`
	Size                 int64    `json:"size"`
	SupportedDeviceTypes []string `json:"supported_device_types"`
	VersionCode          string   `json:"version_code"`
	BundleID             string   `json:"bundle_id,omitempty"`
	PackageName          string   `json:"package_name,omitempty"`
	Module               string   `json:"module"`
	ProductFlavour       string   `json:"product_flavour"`
	BuildType            string   `json:"build_type"`
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
	if env.AppService == nil {
		return errors.New("No App Service defined for handler")
	}

	filterParams := map[string]interface{}{}
	if platformFilter := r.URL.Query().Get("platform"); platformFilter != "" {
		filterParams["platform"] = platformFilter
	}

	app, err := env.AppService.Find(&models.App{Record: models.Record{ID: authorizedAppID}})
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	if env.BitriseAPI == nil {
		return errors.New("No Bitrise API Service defined for handler")
	}

	response, err := newAppVersionsGetResponse(app, env)
	if err != nil {
		return errors.WithStack(err)
	}

	return httpresponse.RespondWithSuccess(w, AppVersionsGetResponse{
		Data: response,
	})
}

func newAppVersionsGetResponse(app *models.App, env *env.AppEnv) ([]AppVersionsGetResponseElement, error) {
	elements := []AppVersionsGetResponseElement{}

	appDetails, err := env.BitriseAPI.GetAppDetails(app.BitriseAPIToken, app.AppSlug)
	if err != nil {
		return nil, err
	}
	appData := AppData{
		Title:       appDetails.Title,
		AppIconURL:  appDetails.AvatarURL,
		ProjectType: appDetails.ProjectType,
	}

	for _, appVersion := range app.AppVersions {
		artifactInfo, err := appVersion.ArtifactInfo()
		if err != nil {
			return nil, err
		}
		elements = append(elements, AppVersionsGetResponseElement{
			AppInfo:              appData,
			AppVersion:           appVersion,
			IPAExportMethod:      artifactInfo.IPAExportMethod,
			Version:              artifactInfo.Version,
			MinimumOS:            artifactInfo.MinimumOS,
			MinimumSDK:           artifactInfo.MinimumSDK,
			Size:                 artifactInfo.Size,
			SupportedDeviceTypes: artifactInfo.SupportedDeviceTypes,
			BundleID:             artifactInfo.BundleID,
			PackageName:          artifactInfo.PackageName,
			VersionCode:          artifactInfo.VersionCode,
			Module:               artifactInfo.Module,
			BuildType:            artifactInfo.BuildType,
			ProductFlavour:       appVersion.ProductFlavour,
		})
	}
	return elements, nil
}

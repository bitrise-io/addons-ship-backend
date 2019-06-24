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

// IosSettingsData ...
type IosSettingsData struct {
	models.IosSettings
	AvailableProvisioningProfiles  []bitrise.ProvisioningProfile `json:"available_provisioning_profiles"`
	AvailableCodeSigningIdentities []bitrise.CodeSigningIdentity `json:"available_code_signing_identities"`
}

// AndroidSettingsData ...
type AndroidSettingsData struct {
	models.AndroidSettings
	AvailableKeystoreFiles       []bitrise.AndroidKeystoreFile `json:"available_keystore_files"`
	AvailableServiceAccountFiles []bitrise.GenericProjectFile  `json:"available_service_account_files"`
}

// AppSettingsGetResponseData ...
type AppSettingsGetResponseData struct {
	*models.AppSettings
	IosSettings     IosSettingsData     `json:"ios_settings"`
	AndroidSettings AndroidSettingsData `json:"android_settings"`
}

// AppSettingsGetResponse ...
type AppSettingsGetResponse struct {
	Data AppSettingsGetResponseData `json:"data"`
}

// AppSettingsGetHandler ...
func AppSettingsGetHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppID, err := GetAuthorizedAppIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if env.AppSettingsService == nil {
		return errors.New("No App Settings Service defined for handler")
	}

	appSettings, err := env.AppSettingsService.Find(&models.AppSettings{AppID: authorizedAppID})
	switch {
	case errors.Cause(err) == gorm.ErrRecordNotFound:
		return httpresponse.RespondWithNotFoundError(w)
	case err != nil:
		return errors.Wrap(err, "SQL Error")
	}

	if env.BitriseAPI == nil {
		return errors.New("No Bitrise API Service defined for handler")
	}

	provisioningProfiles, err := env.BitriseAPI.GetProvisioningProfiles(appSettings.App.BitriseAPIToken, appSettings.App.AppSlug)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch provisioning profiles")
	}

	codeSigningIdentities, err := env.BitriseAPI.GetCodeSigningIdentities(appSettings.App.BitriseAPIToken, appSettings.App.AppSlug)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch code signing identities")
	}

	androidKeyStoreFiles, err := env.BitriseAPI.GetAndroidKeystoreFiles(appSettings.App.BitriseAPIToken, appSettings.App.AppSlug)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch android keystore files")
	}

	serviceAccountfiles, err := env.BitriseAPI.GetServiceAccountFiles(appSettings.App.BitriseAPIToken, appSettings.App.AppSlug)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch service account files")
	}

	iosSettings, err := appSettings.IosSettings()
	if err != nil {
		return errors.WithStack(err)
	}

	androidSettings, err := appSettings.AndroidSettings()
	if err != nil {
		return errors.WithStack(err)
	}

	return httpresponse.RespondWithSuccess(w, AppSettingsGetResponse{
		Data: AppSettingsGetResponseData{
			AppSettings: appSettings,
			IosSettings: IosSettingsData{
				IosSettings:                    iosSettings,
				AvailableProvisioningProfiles:  provisioningProfiles,
				AvailableCodeSigningIdentities: codeSigningIdentities,
			},
			AndroidSettings: AndroidSettingsData{
				AndroidSettings:              androidSettings,
				AvailableKeystoreFiles:       androidKeyStoreFiles,
				AvailableServiceAccountFiles: serviceAccountfiles,
			},
		},
	})
}

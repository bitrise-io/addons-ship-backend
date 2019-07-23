package services_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/bitrise-io/api-utils/providers"
	"github.com/satori/go.uuid"
)

func Test_AppVersionConfigGetHandler(t *testing.T) {
	httpMethod := "GET"
	url := "/apps/{app-slug}/versions/{version-id}/config"
	handler := services.AppVersionConfigGetHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler,
		[]string{"AppVersionService", "AppSettingsService", "FeatureGraphicService", "AWS", "BitriseAPI", "ScreenshotService"},
		ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppVersionService:     &testAppVersionService{},
				FeatureGraphicService: &testFeatureGraphicService{},
				AWS:                &providers.AWSMock{},
				BitriseAPI:         &testBitriseAPI{},
				AppSettingsService: &testAppSettingsService{},
				ScreenshotService:  &testScreenshotService{},
			},
		},
	)

	behavesAsContextCravingHandler(t, httpMethod, url, handler,
		[]ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppVersionID},
		ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
		},
	)

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						appVersion.ArtifactInfoData = json.RawMessage(`{}`)
						appVersion.AppStoreInfoData = json.RawMessage(`{}`)
						return appVersion, nil
					},
				},
				FeatureGraphicService: &testFeatureGraphicService{
					findFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, error) {
						featureGraphic.AppVersion = models.AppVersion{App: models.App{}}
						return featureGraphic, nil
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
					getServiceAccountFilesFn: func(apiToken, appSlug string) ([]bitrise.GenericProjectFile, error) {
						return []bitrise.GenericProjectFile{bitrise.GenericProjectFile{Slug: "service-account-slug"}}, nil
					},
					getAndroidKeystoreFilesFn: func(apiToken, appSlug string) ([]bitrise.AndroidKeystoreFile, error) {
						return []bitrise.AndroidKeystoreFile{
							bitrise.AndroidKeystoreFile{Slug: "android-keystore-slug", UserEnvKey: "ANDROID_KEYSTORE"},
						}, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.AndroidSettingsData = json.RawMessage(`{"selected_service_account":"service-account-slug","selected_keystore_file":"android-keystore-slug"}`)
						return appSettings, nil
					},
				},
				ScreenshotService: &testScreenshotService{},
			},
		},
		)
	})
}

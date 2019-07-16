package services_test

import (
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/c2fo/testify/require"
	uuid "github.com/satori/go.uuid"
)

func Test_BuildWebhookHandler(t *testing.T) {
	httpMethod := "POST"
	url := "/webhook"
	handler := services.BuildWebhookHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppService", "AppSettingsService"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
		requestHeaders: map[string]string{"Bitrise-Event-Type": "build/finished"},
		env: &env.AppEnv{
			AppService:         &testAppService{},
			AppVersionService:  &testAppVersionService{},
			AppSettingsService: &testAppSettingsService{},
		},
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
		requestHeaders: map[string]string{"Bitrise-Event-Type": "build/finished"},
		env: &env.AppEnv{
			AppService:         &testAppService{},
			AppVersionService:  &testAppVersionService{},
			AppSettingsService: &testAppSettingsService{},
		},
	})

	t.Run("when build event type is started", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			requestHeaders: map[string]string{"Bitrise-Event-Type": "build/started"},
		})
	})

	t.Run("when build event type is finished", func(t *testing.T) {
		t.Run("ok - minimal", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppID: uuid.NewV4(),
				},
				requestHeaders: map[string]string{"Bitrise-Event-Type": "build/finished"},
				env: &env.AppEnv{
					AppService: &testAppService{
						findFn: func(app *models.App) (*models.App, error) {
							return app, nil
						},
					},
					AppSettingsService: &testAppSettingsService{
						findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
							return appSettings, nil
						},
					},
					AppVersionService: &testAppVersionService{},
				},
			})
		})

		t.Run("ok - more complex - when ios workflow whitelist is 'all'", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppID: uuid.NewV4(),
				},
				requestHeaders: map[string]string{"Bitrise-Event-Type": "build/finished"},
				env: &env.AppEnv{
					AppService: &testAppService{
						findFn: func(app *models.App) (*models.App, error) {
							return app, nil
						},
					},
					AppSettingsService: &testAppSettingsService{
						findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
							return &models.AppSettings{IosWorkflow: "all"}, nil
						},
					},
					AppVersionService: &testAppVersionService{
						createFn: func(appVersion *models.AppVersion) (*models.AppVersion, []error, error) {
							require.Equal(t, "ios", appVersion.Platform)
							require.Equal(t, "test-build-slug", appVersion.BuildSlug)
							return appVersion, nil, nil
						},
					},
				},
				requestBody: `{"build_slug":"test-build-slug"}`,
			})
		})

		t.Run("ok - more complex - when android workflow whitelist is 'all'", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppID: uuid.NewV4(),
				},
				requestHeaders: map[string]string{"Bitrise-Event-Type": "build/finished"},
				env: &env.AppEnv{
					AppService: &testAppService{
						findFn: func(app *models.App) (*models.App, error) {
							return app, nil
						},
					},
					AppSettingsService: &testAppSettingsService{
						findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
							return &models.AppSettings{AndroidWorkflow: "all"}, nil
						},
					},
					AppVersionService: &testAppVersionService{
						createFn: func(appVersion *models.AppVersion) (*models.AppVersion, []error, error) {
							require.Equal(t, "android", appVersion.Platform)
							require.Equal(t, "test-build-slug", appVersion.BuildSlug)
							return appVersion, nil, nil
						},
					},
				},
				requestBody: `{"build_slug":"test-build-slug"}`,
			})
		})
	})
}

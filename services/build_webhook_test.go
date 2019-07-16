package services_test

import (
	"net/http"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
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
			requestHeaders:     map[string]string{"Bitrise-Event-Type": "build/started"},
			expectedStatusCode: http.StatusOK,
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
				requestBody:        `{}`,
				expectedStatusCode: http.StatusOK,
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
				requestBody:        `{"build_slug":"test-build-slug"}`,
				expectedStatusCode: http.StatusOK,
			})
		})

		t.Run("ok - more complex - when triggered workflow is whitelisted for iOS", func(t *testing.T) {
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
							return &models.AppSettings{IosWorkflow: "ios-wf,ios-wf2"}, nil
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
				requestBody:        `{"build_slug":"test-build-slug","build_triggered_workflow":"ios-wf"}`,
				expectedStatusCode: http.StatusOK,
			})
		})

		t.Run("when request body contains invalid JSON", func(t *testing.T) {
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
				requestBody:        `invalid JSON`,
				expectedStatusCode: http.StatusBadRequest,
				expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Invalid request body, JSON decode failed"},
			})
		})

		t.Run("when db error happens at finding app", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppID: uuid.NewV4(),
				},
				requestHeaders: map[string]string{"Bitrise-Event-Type": "build/finished"},
				env: &env.AppEnv{
					AppService: &testAppService{
						findFn: func(app *models.App) (*models.App, error) {
							return nil, errors.New("SOME-SQL-ERROR")
						},
					},
					AppSettingsService: &testAppSettingsService{
						findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
							return appSettings, nil
						},
					},
					AppVersionService: &testAppVersionService{},
				},
				requestBody:         `{}`,
				expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
			})
		})

		t.Run("when app settings not found in database", func(t *testing.T) {
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
							return nil, gorm.ErrRecordNotFound
						},
					},
					AppVersionService: &testAppVersionService{},
				},
				requestBody:        `{}`,
				expectedStatusCode: http.StatusNotFound,
				expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Not Found"},
			})
		})

		t.Run("when error happens at finding app settings in database", func(t *testing.T) {
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
							return nil, errors.New("SOME-SQL-ERROR")
						},
					},
					AppVersionService: &testAppVersionService{},
				},
				requestBody:         `{}`,
				expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
			})
		})

		t.Run("when validation error is retrieved when creating new ios version", func(t *testing.T) {
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
							return nil, []error{errors.New("SOME-VALIDATION-ERROR")}, nil
						},
					},
				},
				requestBody:        `{"build_slug":"test-build-slug"}`,
				expectedStatusCode: http.StatusUnprocessableEntity,
				expectedResponse: httpresponse.ValidationErrorRespModel{
					Message: "Unprocessable Entity",
					Errors:  []string{"SOME-VALIDATION-ERROR"},
				},
			})
		})

		t.Run("when db error is retrieved when creating new ios version", func(t *testing.T) {
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
							return nil, nil, errors.New("SOME-SQL-ERROR")
						},
					},
				},
				requestBody:         `{"build_slug":"test-build-slug"}`,
				expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
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
				requestBody:        `{"build_slug":"test-build-slug"}`,
				expectedStatusCode: http.StatusOK,
			})
		})

		t.Run("ok - more complex - when triggered workflow is whitelisted for Android", func(t *testing.T) {
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
							return &models.AppSettings{AndroidWorkflow: "android-wf,android-wf2"}, nil
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
				requestBody:        `{"build_slug":"test-build-slug","build_triggered_workflow":"android-wf"}`,
				expectedStatusCode: http.StatusOK,
			})
		})

		t.Run("when validation error is retrieved when creating new android version", func(t *testing.T) {
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
							return nil, []error{errors.New("SOME-VALIDATION-ERROR")}, nil
						},
					},
				},
				requestBody:        `{"build_slug":"test-build-slug"}`,
				expectedStatusCode: http.StatusUnprocessableEntity,
				expectedResponse: httpresponse.ValidationErrorRespModel{
					Message: "Unprocessable Entity",
					Errors:  []string{"SOME-VALIDATION-ERROR"},
				},
			})
		})

		t.Run("when db error is retrieved when creating new android version", func(t *testing.T) {
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
							return nil, nil, errors.New("SOME-SQL-ERROR")
						},
					},
				},
				requestBody:         `{"build_slug":"test-build-slug"}`,
				expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
			})
		})
	})

	t.Run("when build event type is invalid", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			requestHeaders:      map[string]string{"Bitrise-Event-Type": "invalid build event type"},
			expectedInternalErr: "Invalid build event",
		})
	})
}

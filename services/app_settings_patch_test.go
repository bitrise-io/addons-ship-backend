package services_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
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

func Test_AppSettingsPatchHandler(t *testing.T) {
	httpMethod := "PATCH"
	url := "/apps/{app-slug}/settings"
	handler := services.AppSettingsPatchHandler

	testAppID := uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf")

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppSettingsService"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppSettingsService: &testAppSettingsService{
				findFn: func(*models.AppSettings) (*models.AppSettings, error) {
					return nil, nil
				},
			},
		},
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppSettingsService: &testAppSettingsService{},
		},
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						return &models.AppSettings{
							App:                 &models.App{},
							IosSettingsData:     json.RawMessage(`{}`),
							AndroidSettingsData: json.RawMessage(`{}`),
						}, nil
					},
					updateFn: func(appSettings *models.AppSettings, whitelist []string) ([]error, error) {
						return nil, nil
					},
				},
			},
			requestBody:        `{}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppSettingsPatchResponse{
				Data: services.AppSettingsPatchResponseData{AppSettings: &models.AppSettings{}},
			},
		})
	})

	t.Run("ok - more complex", func(t *testing.T) {
		expectedIosSettingsModel := models.IosSettings{AppSKU: "2019061"}
		expectedIosSettings, err := json.Marshal(expectedIosSettingsModel)
		require.NoError(t, err)
		expectedAndroidSettingsModel := models.AndroidSettings{Track: "2019062"}
		expectedAndroidSettings, err := json.Marshal(expectedAndroidSettingsModel)
		require.NoError(t, err)

		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						require.Equal(t, appSettings.AppID, testAppID)
						appSettings.IosSettingsData = json.RawMessage(`{}`)
						appSettings.AndroidSettingsData = json.RawMessage(`{}`)
						return appSettings, nil
					},
					updateFn: func(appSettings *models.AppSettings, whitelist []string) ([]error, error) {
						require.Equal(t, testAppID, appSettings.AppID)
						require.Equal(t, expectedIosSettings, appSettings.IosSettingsData)
						require.Equal(t, expectedAndroidSettings, appSettings.AndroidSettingsData)
						require.Equal(t, "ios-deploy", appSettings.IosWorkflow)
						require.Equal(t, "android-deploy", appSettings.AndroidWorkflow)
						return nil, nil
					},
				},
			},
			requestBody:        `{"ios_settings":{"app_sku":"2019061"},"android_settings":{"track":"2019062"},"ios_workflow":"ios-deploy","android_workflow":"android-deploy"}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppSettingsPatchResponse{
				Data: services.AppSettingsPatchResponseData{
					AppSettings: &models.AppSettings{
						AppID:           testAppID,
						IosWorkflow:     "ios-deploy",
						AndroidWorkflow: "android-deploy",
					},
					IosSettings:     expectedIosSettingsModel,
					AndroidSettings: expectedAndroidSettingsModel,
				},
			},
		})
	})

	t.Run("ok - when prov profile slug list contains not existing", func(t *testing.T) {
		expectedIosSettingsModel := models.IosSettings{AppSKU: "2019061", SelectedAppStoreProvisioningProfiles: []string{"prov-1-slug", "prov-3-slug"}}
		expectedIosSettings, err := json.Marshal(expectedIosSettingsModel)
		require.NoError(t, err)
		expectedAndroidSettingsModel := models.AndroidSettings{Track: "2019062"}
		expectedAndroidSettings, err := json.Marshal(expectedAndroidSettingsModel)
		require.NoError(t, err)

		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						require.Equal(t, appSettings.AppID, testAppID)
						appSettings.IosSettingsData = json.RawMessage(`{}`)
						appSettings.AndroidSettingsData = json.RawMessage(`{}`)
						appSettings.App = &models.App{BitriseAPIToken: "token", AppSlug: "test-slug"}
						return appSettings, nil
					},
					updateFn: func(appSettings *models.AppSettings, whitelist []string) ([]error, error) {
						require.Equal(t, testAppID, appSettings.AppID)
						require.Equal(t, expectedIosSettings, appSettings.IosSettingsData)
						require.Equal(t, expectedAndroidSettings, appSettings.AndroidSettingsData)
						require.Equal(t, "ios-deploy", appSettings.IosWorkflow)
						require.Equal(t, "android-deploy", appSettings.AndroidWorkflow)
						return nil, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getProvisioningProfilesFn: func(token, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						require.Equal(t, token, "token")
						require.Equal(t, appSlug, "test-slug")
						return []bitrise.ProvisioningProfile{{Slug: "prov-1-slug"}, {Slug: "prov-3-slug"}}, nil
					},
				},
			},
			requestBody: `{` +
				`"ios_settings":{"app_sku":"2019061","selected_app_store_provisioning_profiles":["prov-1-slug", "prov-2-slug", "prov-3-slug"]},` +
				`"android_settings":{"track":"2019062"},` +
				`"ios_workflow":"ios-deploy","android_workflow":"android-deploy"` +
				`}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppSettingsPatchResponse{
				Data: services.AppSettingsPatchResponseData{
					AppSettings: &models.AppSettings{
						AppID:           testAppID,
						IosWorkflow:     "ios-deploy",
						AndroidWorkflow: "android-deploy",
					},
					IosSettings:     expectedIosSettingsModel,
					AndroidSettings: expectedAndroidSettingsModel,
				},
			},
		})
	})

	t.Run("when request body is not a valid JSON", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						return &models.AppSettings{
							App:                 &models.App{},
							IosSettingsData:     json.RawMessage(`{}`),
							AndroidSettingsData: json.RawMessage(`{}`),
						}, nil
					},
					updateFn: func(appSettings *models.AppSettings, whitelist []string) ([]error, error) {
						return nil, nil
					},
				},
			},
			requestBody:        `invalid-request-body`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: httpresponse.StandardErrorRespModel{
				Message: "Invalid request body, JSON decode failed",
			},
		})
	})

	t.Run("when app settings to update not found", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						return nil, gorm.ErrRecordNotFound
					},
					updateFn: func(appSettings *models.AppSettings, whitelist []string) ([]error, error) {
						return nil, nil
					},
				},
			},
			requestBody:        `{}`,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Not Found"},
		})
	})

	t.Run("when db error happens at finding the app settings to update", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
					updateFn: func(appSettings *models.AppSettings, whitelist []string) ([]error, error) {
						return nil, nil
					},
				},
			},
			requestBody:         `{}`,
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when validation error happens at update", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						return &models.AppSettings{
							App:                 &models.App{},
							IosSettingsData:     json.RawMessage(`{}`),
							AndroidSettingsData: json.RawMessage(`{}`),
						}, nil
					},
					updateFn: func(appSettings *models.AppSettings, whitelist []string) ([]error, error) {
						return []error{errors.New("SOME-VALIDATION-ERROR")}, nil
					},
				},
			},
			requestBody:        `{}`,
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse: httpresponse.ValidationErrorRespModel{
				Message: "Unprocessable Entity",
				Errors:  []string{"SOME-VALIDATION-ERROR"},
			},
		})
	})

	t.Run("when unexpected db error happens at update", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						return &models.AppSettings{
							App:                 &models.App{},
							IosSettingsData:     json.RawMessage(`{}`),
							AndroidSettingsData: json.RawMessage(`{}`),
						}, nil
					},
					updateFn: func(appSettings *models.AppSettings, whitelist []string) ([]error, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
			},
			requestBody:         `{}`,
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when ios settings data contains an invalid JSON", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						return &models.AppSettings{
							App:                 &models.App{},
							IosSettingsData:     json.RawMessage(`{}`),
							AndroidSettingsData: json.RawMessage(`{}`),
						}, nil
					},
					updateFn: func(appSettings *models.AppSettings, whitelist []string) ([]error, error) {
						appSettings.IosSettingsData = json.RawMessage(`invalid json`)
						return nil, nil
					},
				},
			},
			requestBody:         `{}`,
			expectedInternalErr: "invalid character 'i' looking for beginning of value",
		})
	})

	t.Run("when android settings data contains an invalid JSON", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						return &models.AppSettings{
							App:                 &models.App{},
							IosSettingsData:     json.RawMessage(`{}`),
							AndroidSettingsData: json.RawMessage(`{}`),
						}, nil
					},
					updateFn: func(appSettings *models.AppSettings, whitelist []string) ([]error, error) {
						appSettings.AndroidSettingsData = json.RawMessage(`invalid json`)
						return nil, nil
					},
				},
			},
			requestBody:         `{}`,
			expectedInternalErr: "invalid character 'i' looking for beginning of value",
		})
	})
}

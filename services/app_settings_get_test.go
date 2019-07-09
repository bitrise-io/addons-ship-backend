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

func Test_AppSettingsGetHandler(t *testing.T) {
	httpMethod := "GET"
	url := "/apps/{app-slug}/settings"
	handler := services.AppSettingsGetHandler

	testAppID := uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf")
	testAppSlug := "test-app-slug"
	testAppApiToken := "test-addon-api-token"

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
						require.Equal(t, appSettings.AppID, testAppID)
						return &models.AppSettings{
							App:                 &models.App{AppSlug: testAppSlug, BitriseAPIToken: testAppApiToken},
							IosSettingsData:     json.RawMessage(`{}`),
							AndroidSettingsData: json.RawMessage(`{}`),
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return nil, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return nil, nil
					},
					getAndroidKeystoreFilesFn: func(apiToken, appSlug string) ([]bitrise.AndroidKeystoreFile, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return nil, nil
					},
					getServiceAccountFilesFn: func(apiToken, appSlug string) ([]bitrise.GenericProjectFile, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return nil, nil
					},
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return &bitrise.AppDetails{
							Title:       "Two Brothers",
							ProjectType: "ios",
						}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppSettingsGetResponse{
				Data: services.AppSettingsGetResponseData{
					AppSettings: &models.AppSettings{
						App: &models.App{AppSlug: testAppSlug, BitriseAPIToken: testAppApiToken},
					},
					ProjectType: "ios",
				},
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
						return &models.AppSettings{
							App:                 &models.App{AppSlug: testAppSlug, BitriseAPIToken: testAppApiToken},
							IosSettingsData:     expectedIosSettings,
							AndroidSettingsData: expectedAndroidSettings,
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return []bitrise.ProvisioningProfile{
							bitrise.ProvisioningProfile{Filename: "provision-profile.provisionprofile", Slug: "prov-profile-slug"},
						}, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return []bitrise.CodeSigningIdentity{
							bitrise.CodeSigningIdentity{Filename: "code-signing-id.cert", Slug: "code-signing-slug"},
						}, nil
					},
					getAndroidKeystoreFilesFn: func(apiToken, appSlug string) ([]bitrise.AndroidKeystoreFile, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return []bitrise.AndroidKeystoreFile{
							bitrise.AndroidKeystoreFile{Filename: "android.keystore", Slug: "android-keystore-slug"},
						}, nil
					},
					getServiceAccountFilesFn: func(apiToken, appSlug string) ([]bitrise.GenericProjectFile, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return []bitrise.GenericProjectFile{
							bitrise.GenericProjectFile{Filename: "service-account.json", Slug: "service-account-slug"},
						}, nil
					},
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return &bitrise.AppDetails{
							Title:       "Two Brothers",
							ProjectType: "ios",
						}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppSettingsGetResponse{
				Data: services.AppSettingsGetResponseData{
					AppSettings: &models.AppSettings{
						AppID: testAppID,
					},
					ProjectType: "ios",
					IosSettings: services.IosSettingsData{
						IosSettings: expectedIosSettingsModel,
						AvailableProvisioningProfiles: []bitrise.ProvisioningProfile{
							bitrise.ProvisioningProfile{Filename: "provision-profile.provisionprofile", Slug: "prov-profile-slug"},
						},
						AvailableCodeSigningIdentities: []bitrise.CodeSigningIdentity{
							bitrise.CodeSigningIdentity{Filename: "code-signing-id.cert", Slug: "code-signing-slug"},
						},
					},
					AndroidSettings: services.AndroidSettingsData{
						AndroidSettings: expectedAndroidSettingsModel,
						AvailableKeystoreFiles: []bitrise.AndroidKeystoreFile{
							bitrise.AndroidKeystoreFile{Filename: "android.keystore", Slug: "android-keystore-slug"},
						},
						AvailableServiceAccountFiles: []bitrise.GenericProjectFile{
							bitrise.GenericProjectFile{Filename: "service-account.json", Slug: "service-account-slug"},
						},
					},
				},
			},
		})
	})

	t.Run("when app settings not found", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						return nil, gorm.ErrRecordNotFound
					},
				},
				BitriseAPI: &testBitriseAPI{
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return nil, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return nil, nil
					},
					getAndroidKeystoreFilesFn: func(apiToken, appSlug string) ([]bitrise.AndroidKeystoreFile, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return nil, nil
					},
					getServiceAccountFilesFn: func(apiToken, appSlug string) ([]bitrise.GenericProjectFile, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return nil, nil
					},
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return nil, nil
					},
				},
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Not Found"},
		})
	})

	t.Run("when db error happens at finding the app settings", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
				BitriseAPI: &testBitriseAPI{
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return nil, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return nil, nil
					},
					getAndroidKeystoreFilesFn: func(apiToken, appSlug string) ([]bitrise.AndroidKeystoreFile, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return nil, nil
					},
					getServiceAccountFilesFn: func(apiToken, appSlug string) ([]bitrise.GenericProjectFile, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return nil, nil
					},
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						require.Equal(t, testAppSlug, appSlug)
						require.Equal(t, testAppApiToken, apiToken)
						return nil, nil
					},
				},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when failed to fetch prov profiles", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						require.Equal(t, appSettings.AppID, testAppID)
						return &models.AppSettings{
							App:                 &models.App{AppSlug: testAppSlug, BitriseAPIToken: testAppApiToken},
							IosSettingsData:     json.RawMessage(`{}`),
							AndroidSettingsData: json.RawMessage(`{}`),
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return nil, errors.New("BITRISE-API-ERROR")
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return nil, nil
					},
					getAndroidKeystoreFilesFn: func(apiToken, appSlug string) ([]bitrise.AndroidKeystoreFile, error) {
						return nil, nil
					},
					getServiceAccountFilesFn: func(apiToken, appSlug string) ([]bitrise.GenericProjectFile, error) {
						return nil, nil
					},
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return nil, nil
					},
				},
			},
			expectedInternalErr: "Failed to fetch provisioning profiles: BITRISE-API-ERROR",
		})
	})

	t.Run("when failed to fetch code signing identities", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						require.Equal(t, appSettings.AppID, testAppID)
						return &models.AppSettings{
							App:                 &models.App{AppSlug: testAppSlug, BitriseAPIToken: testAppApiToken},
							IosSettingsData:     json.RawMessage(`{}`),
							AndroidSettingsData: json.RawMessage(`{}`),
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return nil, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return nil, errors.New("BITRISE-API-ERROR")
					},
					getAndroidKeystoreFilesFn: func(apiToken, appSlug string) ([]bitrise.AndroidKeystoreFile, error) {
						return nil, nil
					},
					getServiceAccountFilesFn: func(apiToken, appSlug string) ([]bitrise.GenericProjectFile, error) {
						return nil, nil
					},
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return nil, nil
					},
				},
			},
			expectedInternalErr: "Failed to fetch code signing identities: BITRISE-API-ERROR",
		})
	})

	t.Run("when failed to fetch android keystore files", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						require.Equal(t, appSettings.AppID, testAppID)
						return &models.AppSettings{
							App:                 &models.App{AppSlug: testAppSlug, BitriseAPIToken: testAppApiToken},
							IosSettingsData:     json.RawMessage(`{}`),
							AndroidSettingsData: json.RawMessage(`{}`),
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return nil, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return nil, nil
					},
					getAndroidKeystoreFilesFn: func(apiToken, appSlug string) ([]bitrise.AndroidKeystoreFile, error) {
						return nil, errors.New("BITRISE-API-ERROR")
					},
					getServiceAccountFilesFn: func(apiToken, appSlug string) ([]bitrise.GenericProjectFile, error) {
						return nil, nil
					},
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return nil, nil
					},
				},
			},
			expectedInternalErr: "Failed to fetch android keystore files: BITRISE-API-ERROR",
		})
	})

	t.Run("when failed to fetch service account files", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						require.Equal(t, appSettings.AppID, testAppID)
						return &models.AppSettings{
							App:                 &models.App{AppSlug: testAppSlug, BitriseAPIToken: testAppApiToken},
							IosSettingsData:     json.RawMessage(`{}`),
							AndroidSettingsData: json.RawMessage(`{}`),
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return nil, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return nil, nil
					},
					getAndroidKeystoreFilesFn: func(apiToken, appSlug string) ([]bitrise.AndroidKeystoreFile, error) {
						return nil, nil
					},
					getServiceAccountFilesFn: func(apiToken, appSlug string) ([]bitrise.GenericProjectFile, error) {
						return nil, errors.New("BITRISE-API-ERROR")
					},
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return nil, nil
					},
				},
			},
			expectedInternalErr: "Failed to fetch service account files: BITRISE-API-ERROR",
		})
	})

	t.Run("when failed to fetch app details", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						require.Equal(t, appSettings.AppID, testAppID)
						return &models.AppSettings{
							App:                 &models.App{AppSlug: testAppSlug, BitriseAPIToken: testAppApiToken},
							IosSettingsData:     json.RawMessage(`{}`),
							AndroidSettingsData: json.RawMessage(`{}`),
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return nil, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return nil, nil
					},
					getAndroidKeystoreFilesFn: func(apiToken, appSlug string) ([]bitrise.AndroidKeystoreFile, error) {
						return nil, nil
					},
					getServiceAccountFilesFn: func(apiToken, appSlug string) ([]bitrise.GenericProjectFile, error) {
						return nil, nil
					},
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return nil, errors.New("BITRISE-API-ERROR")
					},
				},
			},
			expectedInternalErr: "Failed to fetch app details: BITRISE-API-ERROR",
		})
	})

	t.Run("when ios settings contains an invalid json in database", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						require.Equal(t, appSettings.AppID, testAppID)
						return &models.AppSettings{
							App:                 &models.App{AppSlug: testAppSlug, BitriseAPIToken: testAppApiToken},
							IosSettingsData:     json.RawMessage(`invalid json`),
							AndroidSettingsData: json.RawMessage(`{}`),
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return nil, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return nil, nil
					},
					getAndroidKeystoreFilesFn: func(apiToken, appSlug string) ([]bitrise.AndroidKeystoreFile, error) {
						return nil, nil
					},
					getServiceAccountFilesFn: func(apiToken, appSlug string) ([]bitrise.GenericProjectFile, error) {
						return nil, nil
					},
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return nil, nil
					},
				},
			},
			expectedInternalErr: "invalid character 'i' looking for beginning of value",
		})
	})

	t.Run("when android settings contains an invalid json in database", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						require.Equal(t, appSettings.AppID, testAppID)
						return &models.AppSettings{
							App:                 &models.App{AppSlug: testAppSlug, BitriseAPIToken: testAppApiToken},
							IosSettingsData:     json.RawMessage(`{}`),
							AndroidSettingsData: json.RawMessage(`invalid json`),
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return nil, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return nil, nil
					},
					getAndroidKeystoreFilesFn: func(apiToken, appSlug string) ([]bitrise.AndroidKeystoreFile, error) {
						return nil, nil
					},
					getServiceAccountFilesFn: func(apiToken, appSlug string) ([]bitrise.GenericProjectFile, error) {
						return nil, nil
					},
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return nil, nil
					},
				},
			},
			expectedInternalErr: "invalid character 'i' looking for beginning of value",
		})
	})
}

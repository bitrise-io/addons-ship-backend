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
	"github.com/bitrise-io/go-utils/pointers"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_AppVersionGetHandler(t *testing.T) {
	httpMethod := "GET"
	url := "/apps/{app-slug}/version{version-id}"
	handler := services.AppVersionGetHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppVersionService", "BitriseAPI"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppVersionService: &testAppVersionService{
				findFn: func(*models.AppVersion) (*models.AppVersion, error) {
					return nil, nil
				},
			},
			BitriseAPI: &testBitriseAPI{},
		},
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppVersionID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppVersionService: &testAppVersionService{},
			BitriseAPI:        &testBitriseAPI{},
		},
	})

	t.Run("when platform is ios", func(t *testing.T) {
		t.Run("ok - minimal", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
				},
				env: &env.AppEnv{
					AppVersionService: &testAppVersionService{
						findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
							require.Equal(t, appVersion.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
							return &models.AppVersion{
								App:              models.App{},
								AppStoreInfoData: json.RawMessage(`{}`),
								Platform:         "ios",
							}, nil
						},
					},
					BitriseAPI: &testBitriseAPI{
						getArtifactsFn: func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error) {
							return []bitrise.ArtifactListElementResponseModel{
								bitrise.ArtifactListElementResponseModel{
									Title: "my-awesome-app.ipa",
									ArtifactMeta: &bitrise.ArtifactMeta{
										ProvisioningInfo: bitrise.ProvisioningInfo{
											DistributionType: "development",
										},
										AppInfo: bitrise.AppInfo{},
									},
								},
							}, nil
						},
						getArtifactPublicPageURLFn: func(string, string, string, string) (string, error) {
							return "", nil
						},
						getAppDetailsFn: func(string, string) (*bitrise.AppDetails, error) {
							return &bitrise.AppDetails{}, nil
						},
					},
				},
				expectedStatusCode: http.StatusOK,
				expectedResponse: services.AppVersionGetResponse{
					Data: services.AppVersionGetResponseData{
						AppVersion:       &models.AppVersion{Platform: "ios"},
						DistributionType: "development",
					},
				},
			})
		})

		t.Run("ok - more complex - when there's an artifact which is IPA and has app store distribution type", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
				},
				env: &env.AppEnv{
					AppVersionService: &testAppVersionService{
						findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
							return &models.AppVersion{
								Version:          "v1.0",
								Platform:         "ios",
								AppStoreInfoData: json.RawMessage(`{"short_description":"Some shorter description"}`),
								App: models.App{
									BitriseAPIToken: "test-api-token",
								},
							}, nil
						},
					},
					BitriseAPI: &testBitriseAPI{
						getArtifactsFn: func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error) {
							return []bitrise.ArtifactListElementResponseModel{
								bitrise.ArtifactListElementResponseModel{
									Title: "my-awesome-app.ipa",
									ArtifactMeta: &bitrise.ArtifactMeta{
										ProvisioningInfo: bitrise.ProvisioningInfo{
											DistributionType: "app-store",
										},
										AppInfo: bitrise.AppInfo{MinimumOS: "11.1"},
									},
								},
							}, nil
						},
						getArtifactPublicPageURLFn: func(string, string, string, string) (string, error) {
							return "http://don.t.go.there", nil
						},
						getAppDetailsFn: func(string, string) (*bitrise.AppDetails, error) {
							return &bitrise.AppDetails{
								Title:       "The Adventures of Stealy",
								AvatarURL:   pointers.NewStringPtr("https://bit.ly/1LixVJu"),
								ProjectType: "ios",
							}, nil
						},
					},
				},
				expectedStatusCode: http.StatusOK,
				expectedResponse: services.AppVersionGetResponse{
					Data: services.AppVersionGetResponseData{
						AppVersion: &models.AppVersion{
							Version:  "v1.0",
							Platform: "ios",
						},
						MinimumOS: "11.1",
						AppInfo: services.AppData{
							Title:       "The Adventures of Stealy",
							AppIconURL:  pointers.NewStringPtr("https://bit.ly/1LixVJu"),
							ProjectType: "ios",
						},
						DistributionType: "app-store",
						AppStoreInfo: models.AppStoreInfo{
							ShortDescription: "Some shorter description",
						},
						PublishEnabled: true,
					},
				},
			})
		})

		t.Run("ok - more complex - when there's an artifact which is XcodeArchive.zip", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
				},
				env: &env.AppEnv{
					AppVersionService: &testAppVersionService{
						findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
							return &models.AppVersion{
								Version:          "v1.0",
								Platform:         "ios",
								AppStoreInfoData: json.RawMessage(`{"short_description":"Some shorter description"}`),
								App: models.App{
									BitriseAPIToken: "test-api-token",
								},
							}, nil
						},
					},
					BitriseAPI: &testBitriseAPI{
						getArtifactsFn: func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error) {
							return []bitrise.ArtifactListElementResponseModel{
								bitrise.ArtifactListElementResponseModel{
									Title: "my-awesome-app.ipa",
									ArtifactMeta: &bitrise.ArtifactMeta{
										ProvisioningInfo: bitrise.ProvisioningInfo{
											DistributionType: "app-store",
										},
										AppInfo: bitrise.AppInfo{MinimumOS: "11.1"},
									},
								},
							}, nil
						},
						getArtifactPublicPageURLFn: func(string, string, string, string) (string, error) {
							return "http://don.t.go.there", nil
						},
						getAppDetailsFn: func(string, string) (*bitrise.AppDetails, error) {
							return &bitrise.AppDetails{
								Title:       "The Adventures of Stealy",
								AvatarURL:   pointers.NewStringPtr("https://bit.ly/1LixVJu"),
								ProjectType: "ios",
							}, nil
						},
					},
				},
				expectedStatusCode: http.StatusOK,
				expectedResponse: services.AppVersionGetResponse{
					Data: services.AppVersionGetResponseData{
						AppVersion: &models.AppVersion{
							Version:  "v1.0",
							Platform: "ios",
						},
						MinimumOS: "11.1",
						AppInfo: services.AppData{
							Title:       "The Adventures of Stealy",
							AppIconURL:  pointers.NewStringPtr("https://bit.ly/1LixVJu"),
							ProjectType: "ios",
						},
						DistributionType: "app-store",
						AppStoreInfo: models.AppStoreInfo{
							ShortDescription: "Some shorter description",
						},
						PublishEnabled: true,
					},
				},
			})
		})

		t.Run("ok - more complex - when there's a development IPA and public install page is enabled", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
				},
				env: &env.AppEnv{
					AppVersionService: &testAppVersionService{
						findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
							return &models.AppVersion{
								Version:          "v1.0",
								Platform:         "ios",
								AppStoreInfoData: json.RawMessage(`{"short_description":"Some shorter description"}`),
								App: models.App{
									BitriseAPIToken: "test-api-token",
								},
							}, nil
						},
					},
					BitriseAPI: &testBitriseAPI{
						getArtifactsFn: func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error) {
							return []bitrise.ArtifactListElementResponseModel{
								bitrise.ArtifactListElementResponseModel{
									Title: "my-awesome-dev-app.ipa",
									ArtifactMeta: &bitrise.ArtifactMeta{
										ProvisioningInfo: bitrise.ProvisioningInfo{
											DistributionType: "development",
										},
										AppInfo: bitrise.AppInfo{MinimumOS: "10.1"},
									},
								},
							}, nil
						},
						getArtifactPublicPageURLFn: func(string, string, string, string) (string, error) {
							return "http://don.t.go.there", nil
						},
						getAppDetailsFn: func(string, string) (*bitrise.AppDetails, error) {
							return &bitrise.AppDetails{
								Title:       "The Adventures of Stealy",
								AvatarURL:   pointers.NewStringPtr("https://bit.ly/1LixVJu"),
								ProjectType: "ios",
							}, nil
						},
					},
				},
				expectedStatusCode: http.StatusOK,
				expectedResponse: services.AppVersionGetResponse{
					Data: services.AppVersionGetResponseData{
						AppVersion: &models.AppVersion{
							Version:  "v1.0",
							Platform: "ios",
						},
						MinimumOS:            "10.1",
						PublicInstallPageURL: "http://don.t.go.there",
						AppInfo: services.AppData{
							Title:       "The Adventures of Stealy",
							AppIconURL:  pointers.NewStringPtr("https://bit.ly/1LixVJu"),
							ProjectType: "ios",
						},
						DistributionType: "development",
						AppStoreInfo: models.AppStoreInfo{
							ShortDescription: "Some shorter description",
						},
					},
				},
			})
		})

		t.Run("ok - more complex - device family", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
				},
				env: &env.AppEnv{
					AppVersionService: &testAppVersionService{
						findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
							return &models.AppVersion{
								Version:          "v1.0",
								Platform:         "ios",
								AppStoreInfoData: json.RawMessage(`{}`),
								App: models.App{
									BitriseAPIToken: "test-api-token",
								},
							}, nil
						},
					},
					BitriseAPI: &testBitriseAPI{
						getArtifactsFn: func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error) {
							return []bitrise.ArtifactListElementResponseModel{
								bitrise.ArtifactListElementResponseModel{
									Title: "my-awesome-dev-app.ipa",
									ArtifactMeta: &bitrise.ArtifactMeta{
										ProvisioningInfo: bitrise.ProvisioningInfo{
											DistributionType: "development",
										},
										AppInfo: bitrise.AppInfo{MinimumOS: "10.1", DeviceFamilyList: []int{1, 2}},
									},
								},
							}, nil
						},
						getArtifactPublicPageURLFn: func(string, string, string, string) (string, error) {
							return "http://don.t.go.there", nil
						},
						getAppDetailsFn: func(string, string) (*bitrise.AppDetails, error) {
							return &bitrise.AppDetails{
								Title:       "The Adventures of Stealy",
								AvatarURL:   pointers.NewStringPtr("https://bit.ly/1LixVJu"),
								ProjectType: "ios",
							}, nil
						},
					},
				},
				expectedStatusCode: http.StatusOK,
				expectedResponse: services.AppVersionGetResponse{
					Data: services.AppVersionGetResponseData{
						AppVersion: &models.AppVersion{
							Version:  "v1.0",
							Platform: "ios",
						},
						MinimumOS:            "10.1",
						SupportedDeviceTypes: []string{"iPhone", "iPod Touch", "iPad"},
						PublicInstallPageURL: "http://don.t.go.there",
						AppInfo: services.AppData{
							Title:       "The Adventures of Stealy",
							AppIconURL:  pointers.NewStringPtr("https://bit.ly/1LixVJu"),
							ProjectType: "ios",
						},
						DistributionType: "development",
					},
				},
			})
		})

		t.Run("ok - more complex - unknown device family", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
				},
				env: &env.AppEnv{
					AppVersionService: &testAppVersionService{
						findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
							return &models.AppVersion{
								Version:          "v1.0",
								Platform:         "ios",
								AppStoreInfoData: json.RawMessage(`{}`),
								App: models.App{
									BitriseAPIToken: "test-api-token",
								},
							}, nil
						},
					},
					BitriseAPI: &testBitriseAPI{
						getArtifactsFn: func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error) {
							return []bitrise.ArtifactListElementResponseModel{
								bitrise.ArtifactListElementResponseModel{
									Title: "my-awesome-dev-app.ipa",
									ArtifactMeta: &bitrise.ArtifactMeta{
										ProvisioningInfo: bitrise.ProvisioningInfo{
											DistributionType: "development",
										},
										AppInfo: bitrise.AppInfo{MinimumOS: "10.1", DeviceFamilyList: []int{12}},
									},
								},
							}, nil
						},
						getArtifactPublicPageURLFn: func(string, string, string, string) (string, error) {
							return "http://don.t.go.there", nil
						},
						getAppDetailsFn: func(string, string) (*bitrise.AppDetails, error) {
							return &bitrise.AppDetails{
								Title:       "The Adventures of Stealy",
								AvatarURL:   pointers.NewStringPtr("https://bit.ly/1LixVJu"),
								ProjectType: "ios",
							}, nil
						},
					},
				},
				expectedStatusCode: http.StatusOK,
				expectedResponse: services.AppVersionGetResponse{
					Data: services.AppVersionGetResponseData{
						AppVersion: &models.AppVersion{
							Version:  "v1.0",
							Platform: "ios",
						},
						MinimumOS:            "10.1",
						SupportedDeviceTypes: []string{"Unknown"},
						PublicInstallPageURL: "http://don.t.go.there",
						AppInfo: services.AppData{
							Title:       "The Adventures of Stealy",
							AppIconURL:  pointers.NewStringPtr("https://bit.ly/1LixVJu"),
							ProjectType: "ios",
						},
						DistributionType: "development",
					},
				},
			})
		})

		t.Run("when error happens at reading app store info", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
				},
				env: &env.AppEnv{
					AppVersionService: &testAppVersionService{
						findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
							require.Equal(t, appVersion.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
							return &models.AppVersion{App: models.App{}, Platform: "ios"}, nil
						},
					},
					BitriseAPI: &testBitriseAPI{
						getArtifactsFn: func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error) {
							return []bitrise.ArtifactListElementResponseModel{
								bitrise.ArtifactListElementResponseModel{
									Title: "my-awesome-dev-app.ipa",
									ArtifactMeta: &bitrise.ArtifactMeta{
										ProvisioningInfo: bitrise.ProvisioningInfo{
											DistributionType: "development",
										},
										AppInfo: bitrise.AppInfo{MinimumOS: "10.1", DeviceFamilyList: []int{1, 2}},
									},
								},
							}, nil
						},
						getArtifactPublicPageURLFn: func(string, string, string, string) (string, error) {
							return "", nil
						},
						getAppDetailsFn: func(string, string) (*bitrise.AppDetails, error) {
							return &bitrise.AppDetails{}, nil
						},
					},
				},
				expectedInternalErr: "unexpected end of JSON input",
			})
		})
	})

	t.Run("when app version platform is invalid", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						require.Equal(t, appVersion.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
						return &models.AppVersion{
							App:              models.App{},
							AppStoreInfoData: json.RawMessage(`{}`),
							Platform:         "invalid platform",
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{
							bitrise.ArtifactListElementResponseModel{
								Title: "my-awesome-app.ipa",
								ArtifactMeta: &bitrise.ArtifactMeta{
									ProvisioningInfo: bitrise.ProvisioningInfo{
										DistributionType: "development",
									},
									AppInfo: bitrise.AppInfo{},
								},
							},
						}, nil
					},
					getArtifactPublicPageURLFn: func(string, string, string, string) (string, error) {
						return "", nil
					},
					getAppDetailsFn: func(string, string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
				},
			},
			expectedInternalErr: "Invalid platform type of app version: invalid platform",
		})
	})

	t.Run("when error happens at fetching artifact data from Bitrise API", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						require.Equal(t, appVersion.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
						return &models.AppVersion{App: models.App{}, Platform: "ios"}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return nil, errors.New("SOME-BITRISE-API-ERROR")
					},
					getArtifactPublicPageURLFn: func(string, string, string, string) (string, error) {
						return "", nil
					},
					getAppDetailsFn: func(string, string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
				},
			},
			expectedInternalErr: "SOME-BITRISE-API-ERROR",
		})
	})

	t.Run("when error happens at fetching artifact public install page from Bitrise API", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						require.Equal(t, appVersion.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
						return &models.AppVersion{App: models.App{}, Platform: "ios"}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{
							bitrise.ArtifactListElementResponseModel{
								Title: "my-awesome-dev-app.ipa",
								ArtifactMeta: &bitrise.ArtifactMeta{
									ProvisioningInfo: bitrise.ProvisioningInfo{
										DistributionType: "development",
									},
									AppInfo: bitrise.AppInfo{MinimumOS: "10.1", DeviceFamilyList: []int{1, 2}},
								},
							},
						}, nil
					},
					getArtifactPublicPageURLFn: func(string, string, string, string) (string, error) {
						return "", errors.New("SOME-BITRISE-API-ERROR")
					},
					getAppDetailsFn: func(string, string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
				},
			},
			expectedInternalErr: "SOME-BITRISE-API-ERROR",
		})
	})

	t.Run("when error happens at fetching app data from Bitrise API", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						require.Equal(t, appVersion.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
						return &models.AppVersion{App: models.App{}, Platform: "ios"}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{
							bitrise.ArtifactListElementResponseModel{
								Title: "my-awesome-dev-app.ipa",
								ArtifactMeta: &bitrise.ArtifactMeta{
									ProvisioningInfo: bitrise.ProvisioningInfo{
										DistributionType: "development",
									},
									AppInfo: bitrise.AppInfo{MinimumOS: "10.1", DeviceFamilyList: []int{1, 2}},
								},
							},
						}, nil
					},
					getArtifactPublicPageURLFn: func(string, string, string, string) (string, error) {
						return "", nil
					},
					getAppDetailsFn: func(string, string) (*bitrise.AppDetails, error) {
						return nil, errors.New("SOME-BITRISE-API-ERROR")
					},
				},
			},
			expectedInternalErr: "SOME-BITRISE-API-ERROR",
		})
	})

	t.Run("error - not found in database", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						return nil, gorm.ErrRecordNotFound
					},
				},
				BitriseAPI: &testBitriseAPI{},
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: httpresponse.StandardErrorRespModel{
				Message: "Not Found",
			},
		})
	})

	t.Run("error - unexpected error in database", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
				BitriseAPI: &testBitriseAPI{},
			},
			expectedStatusCode:  http.StatusNotFound,
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})
}

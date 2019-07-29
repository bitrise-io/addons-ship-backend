package services_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/api-utils/providers"
	"github.com/bitrise-io/go-utils/pointers"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_AppVersionIosConfigGetHandler(t *testing.T) {
	httpMethod := "GET"
	url := "/apps/{app-slug}/versions/{version-id}/ios-config"
	handler := services.AppVersionIosConfigGetHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler,
		[]string{"AppVersionService", "AppSettingsService", "AWS", "BitriseAPI", "ScreenshotService"},
		ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppVersionService:  &testAppVersionService{},
				AWS:                &providers.AWSMock{},
				BitriseAPI:         &testBitriseAPI{},
				AppSettingsService: &testAppSettingsService{},
				ScreenshotService:  &testScreenshotService{},
			},
		},
	)

	testAppID := uuid.FromStringOrNil("e3338a14-938a-4e5a-b0fe-e943ed3fb6d0")
	testAppVersionID := uuid.FromStringOrNil("1ca9503a-6230-4140-9fca-3867b6640ce3")

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
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return []bitrise.ProvisioningProfile{bitrise.ProvisioningProfile{Slug: "prov-profile-slug"}}, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return []bitrise.CodeSigningIdentity{bitrise.CodeSigningIdentity{Slug: "code-signing-slug"}}, nil
					},
					getArtifactsFn: func(apiToken, appSlug, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{}, nil
					},
					getArtifactFn: func(apiToken, appSlug, buildSlug, artifactSlug string) (*bitrise.ArtifactShowResponseItemModel, error) {
						return nil, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.IosSettingsData = json.RawMessage(`{"selected_app_store_provisioning_profile":"prov-profile-slug","selected_code_signing_identity":"code-signing-slug"}`)
						return appSettings, nil
					},
				},
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionIosConfigGetResponse{
				MetaData: services.IosConfigMetaData{
					ListingInfoMap: map[string]services.IosListingInfo{
						"en-US": services.IosListingInfo{Screenshots: map[string][]string{}},
					},
				},
			},
		})
	})

	t.Run("ok - more complex", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						appVersion.ArtifactInfoData = json.RawMessage(`{"package_name":"myPackage"}`)
						appVersion.AppStoreInfoData = json.RawMessage(`{` +
							`"full_description":"A bit longer description","promotional_text":"This is an awesome app, you should download it"` +
							`,"support_url":"http://we-will-help.you","marketing_url":"http://purchase-the.app"` +
							`,"keywords":"awesome,awesomeapp,awesomeness"` +
							`}`)
						appVersion.App = models.App{AppSlug: "test-app-slug", BitriseAPIToken: "test-api-token"}
						appVersion.AppID = testAppID
						return appVersion, nil
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return fmt.Sprintf("http://presigned.url/%s", path), nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						require.Equal(t, "test-app-slug", appSlug)
						require.Equal(t, "test-api-token", apiToken)
						return &bitrise.AppDetails{Title: "my-awesome-app"}, nil
					},
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						require.Equal(t, "test-app-slug", appSlug)
						require.Equal(t, "test-api-token", apiToken)
						return []bitrise.ProvisioningProfile{
							bitrise.ProvisioningProfile{Slug: "prov-profile-slug", DownloadURL: "http://provisioning-profile.url"},
						}, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						require.Equal(t, "test-app-slug", appSlug)
						require.Equal(t, "test-api-token", apiToken)
						return []bitrise.CodeSigningIdentity{
							bitrise.CodeSigningIdentity{
								Slug:                "code-signing-slug",
								DownloadURL:         "http:/code-signing.url",
								CertificatePassword: "my-super-password",
							},
						}, nil
					},
					getArtifactsFn: func(apiToken, appSlug, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{
							bitrise.ArtifactListElementResponseModel{Title: "app.xcarchive.zip", Slug: "test-artifact-slug"},
						}, nil
					},
					getArtifactFn: func(apiToken, appSlug, buildSlug, artifactSlug string) (*bitrise.ArtifactShowResponseItemModel, error) {
						require.Equal(t, "test-app-slug", appSlug)
						require.Equal(t, "test-api-token", apiToken)
						require.Equal(t, "test-artifact-slug", artifactSlug)
						return &bitrise.ArtifactShowResponseItemModel{DownloadPath: pointers.NewStringPtr("http://the-url-for-artifact.io")}, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						require.Equal(t, testAppID, appSettings.AppID)
						appSettings.IosSettingsData = json.RawMessage(`{` +
							`"selected_app_store_provisioning_profile":"prov-profile-slug","selected_code_signing_identity":"code-signing-slug"` +
							`,"include_bit_code":true,"app_sku":"some-string","apple_developer_account_email":"my.apple@email.com"` +
							`,"app_specific_password":"my-super-secret-pass"` +
							`}`)
						return appSettings, nil
					},
				},
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						testAppVersion := models.AppVersion{
							Record: models.Record{ID: testAppVersionID},
							App:    models.App{AppSlug: "test-app-slug"},
						}
						return []models.Screenshot{
							models.Screenshot{Record: models.Record{ID: uuid.FromStringOrNil("17ec78c9-e3a8-41ee-b3bd-2df9b4117aa2")}, UploadableObject: models.UploadableObject{Filename: "iPhone XS Max.png"}, ScreenSize: "6.5 inch", DeviceType: "iPhone XS Max", AppVersion: testAppVersion},
							models.Screenshot{Record: models.Record{ID: uuid.FromStringOrNil("d5c8564f-eef4-490a-a7fd-8d3050893320")}, UploadableObject: models.UploadableObject{Filename: "iPad Pro.png"}, ScreenSize: "12.9 inch", DeviceType: "iPad Pro", AppVersion: testAppVersion},
							models.Screenshot{Record: models.Record{ID: uuid.FromStringOrNil("e4d64d18-e414-4fa3-8583-f94a06b4f9a9")}, UploadableObject: models.UploadableObject{Filename: "iPhone XS Max 2.png"}, ScreenSize: "6.5 inch", DeviceType: "iPhone XS Max", AppVersion: testAppVersion},
							models.Screenshot{Record: models.Record{ID: uuid.FromStringOrNil("4faa287f-afee-46aa-bd6b-553ab11a959c")}, UploadableObject: models.UploadableObject{Filename: "iPhone XS.png"}, ScreenSize: "5.8 inch", DeviceType: "iPhone XS", AppVersion: testAppVersion},
							models.Screenshot{Record: models.Record{ID: uuid.FromStringOrNil("27cee0a1-1afd-4280-8d9f-f22526dc3d16")}, UploadableObject: models.UploadableObject{Filename: "iPad Pro 2.png"}, ScreenSize: "12.9 inch", DeviceType: "iPad Pro", AppVersion: testAppVersion},
						}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionIosConfigGetResponse{
				MetaData: services.IosConfigMetaData{
					ListingInfoMap: map[string]services.IosListingInfo{
						"en-US": services.IosListingInfo{
							Screenshots: map[string][]string{
								"12.9 inch": []string{
									"http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/iPad Pro (12.9 inch)/d5c8564f-eef4-490a-a7fd-8d3050893320.png",
									"http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/iPad Pro (12.9 inch)/27cee0a1-1afd-4280-8d9f-f22526dc3d16.png",
								},
								"6.5 inch": []string{
									"http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/iPhone XS Max (6.5 inch)/17ec78c9-e3a8-41ee-b3bd-2df9b4117aa2.png",
									"http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/iPhone XS Max (6.5 inch)/e4d64d18-e414-4fa3-8583-f94a06b4f9a9.png",
								},
								"5.8 inch": []string{"http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/iPhone XS (5.8 inch)/4faa287f-afee-46aa-bd6b-553ab11a959c.png"},
							},
							Description:     "A bit longer description",
							PromotionalText: "This is an awesome app, you should download it",
							SupportURL:      "http://we-will-help.you",
							SoftwareURL:     "http://purchase-the.app",
							Keywords:        []string{"awesome", "awesomeapp", "awesomeness"},
						},
					},
					Signing: services.Signing{
						DistributionCertificateURL:        "http:/code-signing.url",
						DistributionCertificatePasshprase: "my-super-password",
						AppStoreProfileURL:                "http://provisioning-profile.url",
					},
					ExportOptions:            services.ExportOptions{IncludeBitcode: true},
					SKU:                      "some-string",
					AppleUser:                "my.apple@email.com",
					AppleAppSpecificPassword: "my-super-secret-pass",
				},
				Artifacts: []string{"http://the-url-for-artifact.io"},
			},
		})
	})

	t.Run("when it's failed to find app version", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						return appVersion, gorm.ErrRecordNotFound
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
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return []bitrise.ProvisioningProfile{bitrise.ProvisioningProfile{Slug: "prov-profile-slug"}}, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return []bitrise.CodeSigningIdentity{bitrise.CodeSigningIdentity{Slug: "code-signing-slug"}}, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.IosSettingsData = json.RawMessage(`{"selected_app_store_provisioning_profile":"prov-profile-slug","selected_code_signing_identity":"code-signing-slug"}`)
						return appSettings, nil
					},
				},
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedInternalErr: "SQL Error: record not found",
		})
	})

	t.Run("when failed to get store info from app version", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						appVersion.ArtifactInfoData = json.RawMessage(`{}`)
						appVersion.AppStoreInfoData = json.RawMessage(`invalid JSON`)
						return appVersion, nil
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
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return []bitrise.ProvisioningProfile{bitrise.ProvisioningProfile{Slug: "prov-profile-slug"}}, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return []bitrise.CodeSigningIdentity{bitrise.CodeSigningIdentity{Slug: "code-signing-slug"}}, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.IosSettingsData = json.RawMessage(`{"selected_app_store_provisioning_profile":"prov-profile-slug","selected_code_signing_identity":"code-signing-slug"}`)
						return appSettings, nil
					},
				},
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedInternalErr: "invalid character 'i' looking for beginning of value",
		})
	})

	t.Run("when failed to find screenshots", func(t *testing.T) {
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
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return []bitrise.ProvisioningProfile{bitrise.ProvisioningProfile{Slug: "prov-profile-slug"}}, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return []bitrise.CodeSigningIdentity{bitrise.CodeSigningIdentity{Slug: "code-signing-slug"}}, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.IosSettingsData = json.RawMessage(`{"selected_app_store_provisioning_profile":"prov-profile-slug","selected_code_signing_identity":"code-signing-slug"}`)
						return appSettings, nil
					},
				},
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, errors.New("SOME-SQL-ERROR")
					},
				},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when failed to generate presigned URL for screenshot", func(t *testing.T) {
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
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", errors.New("SOME-AWS-ERROR")
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return []bitrise.ProvisioningProfile{bitrise.ProvisioningProfile{Slug: "prov-profile-slug"}}, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return []bitrise.CodeSigningIdentity{bitrise.CodeSigningIdentity{Slug: "code-signing-slug"}}, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.IosSettingsData = json.RawMessage(`{"selected_app_store_provisioning_profile":"prov-profile-slug","selected_code_signing_identity":"code-signing-slug"}`)
						return appSettings, nil
					},
				},
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{models.Screenshot{DeviceType: "Apple Watch"}}, nil
					},
				},
			},
			expectedInternalErr: "SOME-AWS-ERROR",
		})
	})

	t.Run("when it's failed to find app setting", func(t *testing.T) {
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
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return []bitrise.ProvisioningProfile{bitrise.ProvisioningProfile{Slug: "prov-profile-slug"}}, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return []bitrise.CodeSigningIdentity{bitrise.CodeSigningIdentity{Slug: "code-signing-slug"}}, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						return appSettings, gorm.ErrRecordNotFound
					},
				},
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedInternalErr: "SQL Error: record not found",
		})
	})

	t.Run("when failed to get ios settings from app settings", func(t *testing.T) {
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
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return []bitrise.ProvisioningProfile{bitrise.ProvisioningProfile{Slug: "prov-profile-slug"}}, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return []bitrise.CodeSigningIdentity{bitrise.CodeSigningIdentity{Slug: "code-signing-slug"}}, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.IosSettingsData = json.RawMessage(`invalid JSON`)
						return appSettings, nil
					},
				},
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedInternalErr: "invalid character 'i' looking for beginning of value",
		})
	})

	t.Run("when failed to fetch provisioning profiles from API", func(t *testing.T) {
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
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return []bitrise.ProvisioningProfile{}, errors.New("SOME-BITRISE-API-ERROR")
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return []bitrise.CodeSigningIdentity{bitrise.CodeSigningIdentity{Slug: "code-signing-slug"}}, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.IosSettingsData = json.RawMessage(`{"selected_app_store_provisioning_profile":"prov-profile-slug","selected_code_signing_identity":"code-signing-slug"}`)
						return appSettings, nil
					},
				},
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedInternalErr: "SOME-BITRISE-API-ERROR",
		})
	})

	t.Run("when no matching provisioning profile found", func(t *testing.T) {
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
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return []bitrise.ProvisioningProfile{bitrise.ProvisioningProfile{Slug: "not-matching-slug"}}, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return []bitrise.CodeSigningIdentity{bitrise.CodeSigningIdentity{Slug: "code-signing-slug"}}, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.IosSettingsData = json.RawMessage(`{"selected_app_store_provisioning_profile":"prov-profile-slug","selected_code_signing_identity":"code-signing-slug"}`)
						return appSettings, nil
					},
				},
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Not Found"},
		})
	})

	t.Run("when failed to fetch code signing identities from API", func(t *testing.T) {
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
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return []bitrise.ProvisioningProfile{bitrise.ProvisioningProfile{Slug: "prov-profile-slug"}}, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return []bitrise.CodeSigningIdentity{}, errors.New("SOME-BITRISE-API-ERROR")
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.IosSettingsData = json.RawMessage(`{"selected_app_store_provisioning_profile":"prov-profile-slug","selected_code_signing_identity":"code-signing-slug"}`)
						return appSettings, nil
					},
				},
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedInternalErr: "SOME-BITRISE-API-ERROR",
		})
	})

	t.Run("when no matching code signing identity found", func(t *testing.T) {
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
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return []bitrise.ProvisioningProfile{bitrise.ProvisioningProfile{Slug: "prov-profile-slug"}}, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return []bitrise.CodeSigningIdentity{bitrise.CodeSigningIdentity{Slug: "not-matching-slug"}}, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.IosSettingsData = json.RawMessage(`{"selected_app_store_provisioning_profile":"prov-profile-slug","selected_code_signing_identity":"code-signing-slug"}`)
						return appSettings, nil
					},
				},
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Not Found"},
		})
	})

	t.Run("when failed to fetch artifacts from API", func(t *testing.T) {
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
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return []bitrise.ProvisioningProfile{bitrise.ProvisioningProfile{Slug: "prov-profile-slug"}}, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return []bitrise.CodeSigningIdentity{bitrise.CodeSigningIdentity{Slug: "code-signing-slug"}}, nil
					},
					getArtifactsFn: func(apiToken, appSlug, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{}, errors.New("SOME-BITRISE-API-ERROR")
					},
					getArtifactFn: func(apiToken, appSlug, buildSlug, artifactSlug string) (*bitrise.ArtifactShowResponseItemModel, error) {
						return nil, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.IosSettingsData = json.RawMessage(`{"selected_app_store_provisioning_profile":"prov-profile-slug","selected_code_signing_identity":"code-signing-slug"}`)
						return appSettings, nil
					},
				},
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedInternalErr: "SOME-BITRISE-API-ERROR",
		})
	})

	t.Run("when failed to fetch artifact details from API", func(t *testing.T) {
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
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return []bitrise.ProvisioningProfile{bitrise.ProvisioningProfile{Slug: "prov-profile-slug"}}, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return []bitrise.CodeSigningIdentity{bitrise.CodeSigningIdentity{Slug: "code-signing-slug"}}, nil
					},
					getArtifactsFn: func(apiToken, appSlug, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{bitrise.ArtifactListElementResponseModel{Title: "app.xcarchive.zip", Slug: "test-artifact-slug"}}, nil
					},
					getArtifactFn: func(apiToken, appSlug, buildSlug, artifactSlug string) (*bitrise.ArtifactShowResponseItemModel, error) {
						return nil, errors.New("SOME-BITRISE-API-ERROR")
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.IosSettingsData = json.RawMessage(`{"selected_app_store_provisioning_profile":"prov-profile-slug","selected_code_signing_identity":"code-signing-slug"}`)
						return appSettings, nil
					},
				},
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedInternalErr: "SOME-BITRISE-API-ERROR",
		})
	})

	t.Run("when failed to fetched artifact details doesn't contain download URL", func(t *testing.T) {
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
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
					getProvisioningProfilesFn: func(apiToken, appSlug string) ([]bitrise.ProvisioningProfile, error) {
						return []bitrise.ProvisioningProfile{bitrise.ProvisioningProfile{Slug: "prov-profile-slug"}}, nil
					},
					getCodeSigningIdentitiesFn: func(apiToken, appSlug string) ([]bitrise.CodeSigningIdentity, error) {
						return []bitrise.CodeSigningIdentity{bitrise.CodeSigningIdentity{Slug: "code-signing-slug"}}, nil
					},
					getArtifactsFn: func(apiToken, appSlug, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{bitrise.ArtifactListElementResponseModel{Title: "app.xcarchive.zip", Slug: "test-artifact-slug"}}, nil
					},
					getArtifactFn: func(apiToken, appSlug, buildSlug, artifactSlug string) (*bitrise.ArtifactShowResponseItemModel, error) {
						return &bitrise.ArtifactShowResponseItemModel{Title: pointers.NewStringPtr("URLless artifact")}, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.IosSettingsData = json.RawMessage(`{"selected_app_store_provisioning_profile":"prov-profile-slug","selected_code_signing_identity":"code-signing-slug"}`)
						return appSettings, nil
					},
				},
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedInternalErr: "Failed to get download URL for artifact",
		})
	})
}

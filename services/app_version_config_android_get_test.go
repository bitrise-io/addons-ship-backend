package services_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	"github.com/satori/go.uuid"
)

func Test_AppVersionAndroidConfigGetHandler(t *testing.T) {
	httpMethod := "GET"
	url := "/apps/{app-slug}/versions/{version-id}/android-config"
	handler := services.AppVersionAndroidConfigGetHandler

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

	testAppVersionID := uuid.FromStringOrNil("1ca9503a-6230-4140-9fca-3867b6640ce3")
	testFeatureGraphicID := uuid.FromStringOrNil("6154234a-9146-4a20-b43f-f0292d98017a")

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
					getArtifactsFn: func(apiToken, appSlug, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{}, nil
					},
					getArtifactFn: func(apiToken, appSlug, buildSlug, artifactSlug string) (*bitrise.ArtifactShowResponseItemModel, error) {
						return nil, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.AndroidSettingsData = json.RawMessage(`{"selected_service_account":"service-account-slug","selected_keystore_file":"android-keystore-slug"}`)
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
			expectedResponse:   services.AppVersionAndroidConfigGetResponse{Artifacts: []string{}},
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
						appVersion.AppStoreInfoData = json.RawMessage(`{"short_description":"Description","full_description":"A bit longer description","whats_new":"This is what is new"}`)
						appVersion.App = models.App{AppSlug: "test-app-slug", BitriseAPIToken: "test-api-token"}
						return appVersion, nil
					},
				},
				FeatureGraphicService: &testFeatureGraphicService{
					findFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, error) {
						featureGraphic.ID = testFeatureGraphicID
						featureGraphic.Filename = "feature_graphic.png"
						featureGraphic.AppVersion = models.AppVersion{
							Record: models.Record{ID: testAppVersionID},
							App:    models.App{AppSlug: "test-app-slug"},
						}
						return featureGraphic, nil
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
					getServiceAccountFilesFn: func(apiToken, appSlug string) ([]bitrise.GenericProjectFile, error) {
						require.Equal(t, "test-app-slug", appSlug)
						require.Equal(t, "test-api-token", apiToken)
						return []bitrise.GenericProjectFile{bitrise.GenericProjectFile{
							Slug:        "service-account-slug",
							DownloadURL: "http://service-account-json.url",
						}}, nil
					},
					getAndroidKeystoreFilesFn: func(apiToken, appSlug string) ([]bitrise.AndroidKeystoreFile, error) {
						require.Equal(t, "test-app-slug", appSlug)
						require.Equal(t, "test-api-token", apiToken)
						return []bitrise.AndroidKeystoreFile{
							bitrise.AndroidKeystoreFile{
								Slug:        "android-keystore-slug",
								UserEnvKey:  "ANDROID_KEYSTORE",
								DownloadURL: "http://android.keystore.url",
								ExposedMetadataStore: bitrise.ExposedMetadataStore{
									Password:           "my-secret-password",
									Alias:              "AnDrOID-KeySTore",
									PrivateKeyPassword: "my-private-key-pass",
								},
							},
						}, nil
					},
					getArtifactsFn: func(apiToken, appSlug, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						require.Equal(t, "test-app-slug", appSlug)
						require.Equal(t, "test-api-token", apiToken)
						return []bitrise.ArtifactListElementResponseModel{
							bitrise.ArtifactListElementResponseModel{Title: "app.aab", Slug: "test-artifact-slug"},
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
						appSettings.AndroidSettingsData = json.RawMessage(`{"selected_service_account":"service-account-slug","selected_keystore_file":"android-keystore-slug"}`)
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
							models.Screenshot{Record: models.Record{ID: uuid.FromStringOrNil("17ec78c9-e3a8-41ee-b3bd-2df9b4117aa2")}, UploadableObject: models.UploadableObject{Filename: "tv.png"}, ScreenSize: "tv", DeviceType: "TV", AppVersion: testAppVersion},
							models.Screenshot{Record: models.Record{ID: uuid.FromStringOrNil("d5c8564f-eef4-490a-a7fd-8d3050893320")}, UploadableObject: models.UploadableObject{Filename: "wear.png"}, ScreenSize: "wear", DeviceType: "Watch", AppVersion: testAppVersion},
							models.Screenshot{Record: models.Record{ID: uuid.FromStringOrNil("e4d64d18-e414-4fa3-8583-f94a06b4f9a9")}, UploadableObject: models.UploadableObject{Filename: "phone.png"}, ScreenSize: "phone", DeviceType: "Phone", AppVersion: testAppVersion},
							models.Screenshot{Record: models.Record{ID: uuid.FromStringOrNil("4faa287f-afee-46aa-bd6b-553ab11a959c")}, UploadableObject: models.UploadableObject{Filename: "ten_inch.png"}, ScreenSize: "ten_inch", DeviceType: "Tablet", AppVersion: testAppVersion},
							models.Screenshot{Record: models.Record{ID: uuid.FromStringOrNil("27cee0a1-1afd-4280-8d9f-f22526dc3d16")}, UploadableObject: models.UploadableObject{Filename: "seven_inch.png"}, ScreenSize: "seven_inch", DeviceType: "Tablet", AppVersion: testAppVersion},
						}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionAndroidConfigGetResponse{
				MetaData: services.MetaData{
					ListingInfo: services.ListingInfo{
						ShortDescription: "Description",
						FullDescription:  "A bit longer description",
						WhatsNew:         "This is what is new",
						Title:            "my-awesome-app",
						FeatureGraphic:   "http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/6154234a-9146-4a20-b43f-f0292d98017a.png",
						Screenshots: services.Screenshots{
							Tv:        []string{"http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/TV (tv)/17ec78c9-e3a8-41ee-b3bd-2df9b4117aa2.png"},
							Wear:      []string{"http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/Watch (wear)/d5c8564f-eef4-490a-a7fd-8d3050893320.png"},
							Phone:     []string{"http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/Phone (phone)/e4d64d18-e414-4fa3-8583-f94a06b4f9a9.png"},
							TenInch:   []string{"http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/Tablet (ten_inch)/4faa287f-afee-46aa-bd6b-553ab11a959c.png"},
							SevenInch: []string{"http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/Tablet (seven_inch)/27cee0a1-1afd-4280-8d9f-f22526dc3d16.png"},
						},
					},
					PackageName:        "myPackage",
					ServiceAccountJSON: "http://service-account-json.url",
					Keystore: services.Keystore{
						URL:         "http://android.keystore.url",
						Password:    "my-secret-password",
						Alias:       "AnDrOID-KeySTore",
						KeyPassword: "my-private-key-pass",
					},
				},
				Artifacts: []string{"http://the-url-for-artifact.io"},
			},
		})
	})

	t.Run("ok - more complex - split APK", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						appVersion.ArtifactInfoData = json.RawMessage(`{"package_name":"myPackage"}`)
						appVersion.AppStoreInfoData = json.RawMessage(`{"short_description":"Description","full_description":"A bit longer description","whats_new":"This is what is new"}`)
						appVersion.App = models.App{AppSlug: "test-app-slug", BitriseAPIToken: "test-api-token"}
						return appVersion, nil
					},
				},
				FeatureGraphicService: &testFeatureGraphicService{
					findFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, error) {
						featureGraphic.ID = testFeatureGraphicID
						featureGraphic.Filename = "feature_graphic.png"
						featureGraphic.AppVersion = models.AppVersion{
							Record: models.Record{ID: testAppVersionID},
							App:    models.App{AppSlug: "test-app-slug"},
						}
						return featureGraphic, nil
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
					getServiceAccountFilesFn: func(apiToken, appSlug string) ([]bitrise.GenericProjectFile, error) {
						require.Equal(t, "test-app-slug", appSlug)
						require.Equal(t, "test-api-token", apiToken)
						return []bitrise.GenericProjectFile{bitrise.GenericProjectFile{
							Slug:        "service-account-slug",
							DownloadURL: "http://service-account-json.url",
						}}, nil
					},
					getAndroidKeystoreFilesFn: func(apiToken, appSlug string) ([]bitrise.AndroidKeystoreFile, error) {
						require.Equal(t, "test-app-slug", appSlug)
						require.Equal(t, "test-api-token", apiToken)
						return []bitrise.AndroidKeystoreFile{
							bitrise.AndroidKeystoreFile{
								Slug:        "android-keystore-slug",
								UserEnvKey:  "ANDROID_KEYSTORE",
								DownloadURL: "http://android.keystore.url",
								ExposedMetadataStore: bitrise.ExposedMetadataStore{
									Password:           "my-secret-password",
									Alias:              "AnDrOID-KeySTore",
									PrivateKeyPassword: "my-private-key-pass",
								},
							},
						}, nil
					},
					getArtifactsFn: func(apiToken, appSlug, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						require.Equal(t, "test-app-slug", appSlug)
						require.Equal(t, "test-api-token", apiToken)
						return []bitrise.ArtifactListElementResponseModel{
							bitrise.ArtifactListElementResponseModel{
								Title: "app-armeabi-my-awesome-app.apk", Slug: "test-artifact-slug-1",
								ArtifactMeta: &bitrise.ArtifactMeta{
									AppInfo: bitrise.AppInfo{AppName: "My Awesome app", VersionName: "1.0", PackageName: "my.package"},
								},
							},
							bitrise.ArtifactListElementResponseModel{
								Title: "app-x86-my-awesome-app.apk", Slug: "test-artifact-slug-2",
								ArtifactMeta: &bitrise.ArtifactMeta{
									AppInfo: bitrise.AppInfo{AppName: "My Awesome app", VersionName: "1.0", PackageName: "my.package"},
								},
							},
						}, nil
					},
					getArtifactFn: func(apiToken, appSlug, buildSlug, artifactSlug string) (*bitrise.ArtifactShowResponseItemModel, error) {
						require.Equal(t, "test-app-slug", appSlug)
						require.Equal(t, "test-api-token", apiToken)
						require.Contains(t, []string{"test-artifact-slug-1", "test-artifact-slug-2"}, artifactSlug)
						artifactURL := "http://the-url-for-artifact.io/1"
						if artifactSlug == "test-artifact-slug-2" {
							artifactURL = "http://the-url-for-artifact.io/2"
						}
						return &bitrise.ArtifactShowResponseItemModel{DownloadPath: pointers.NewStringPtr(artifactURL)}, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.AndroidSettingsData = json.RawMessage(`{"selected_service_account":"service-account-slug","selected_keystore_file":"android-keystore-slug"}`)
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
							models.Screenshot{Record: models.Record{ID: uuid.FromStringOrNil("17ec78c9-e3a8-41ee-b3bd-2df9b4117aa2")}, UploadableObject: models.UploadableObject{Filename: "tv.png"}, ScreenSize: "tv", DeviceType: "TV", AppVersion: testAppVersion},
							models.Screenshot{Record: models.Record{ID: uuid.FromStringOrNil("d5c8564f-eef4-490a-a7fd-8d3050893320")}, UploadableObject: models.UploadableObject{Filename: "wear.png"}, ScreenSize: "wear", DeviceType: "Watch", AppVersion: testAppVersion},
							models.Screenshot{Record: models.Record{ID: uuid.FromStringOrNil("e4d64d18-e414-4fa3-8583-f94a06b4f9a9")}, UploadableObject: models.UploadableObject{Filename: "phone.png"}, ScreenSize: "phone", DeviceType: "Phone", AppVersion: testAppVersion},
							models.Screenshot{Record: models.Record{ID: uuid.FromStringOrNil("4faa287f-afee-46aa-bd6b-553ab11a959c")}, UploadableObject: models.UploadableObject{Filename: "ten_inch.png"}, ScreenSize: "ten_inch", DeviceType: "Tablet", AppVersion: testAppVersion},
							models.Screenshot{Record: models.Record{ID: uuid.FromStringOrNil("27cee0a1-1afd-4280-8d9f-f22526dc3d16")}, UploadableObject: models.UploadableObject{Filename: "seven_inch.png"}, ScreenSize: "seven_inch", DeviceType: "Tablet", AppVersion: testAppVersion},
						}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionAndroidConfigGetResponse{
				MetaData: services.MetaData{
					ListingInfo: services.ListingInfo{
						ShortDescription: "Description",
						FullDescription:  "A bit longer description",
						WhatsNew:         "This is what is new",
						Title:            "my-awesome-app",
						FeatureGraphic:   "http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/6154234a-9146-4a20-b43f-f0292d98017a.png",
						Screenshots: services.Screenshots{
							Tv:        []string{"http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/TV (tv)/17ec78c9-e3a8-41ee-b3bd-2df9b4117aa2.png"},
							Wear:      []string{"http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/Watch (wear)/d5c8564f-eef4-490a-a7fd-8d3050893320.png"},
							Phone:     []string{"http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/Phone (phone)/e4d64d18-e414-4fa3-8583-f94a06b4f9a9.png"},
							TenInch:   []string{"http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/Tablet (ten_inch)/4faa287f-afee-46aa-bd6b-553ab11a959c.png"},
							SevenInch: []string{"http://presigned.url/test-app-slug/1ca9503a-6230-4140-9fca-3867b6640ce3/Tablet (seven_inch)/27cee0a1-1afd-4280-8d9f-f22526dc3d16.png"},
						},
					},
					PackageName:        "myPackage",
					ServiceAccountJSON: "http://service-account-json.url",
					Keystore: services.Keystore{
						URL:         "http://android.keystore.url",
						Password:    "my-secret-password",
						Alias:       "AnDrOID-KeySTore",
						KeyPassword: "my-private-key-pass",
					},
				},
				Artifacts: []string{"http://the-url-for-artifact.io/1", "http://the-url-for-artifact.io/2"},
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
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedInternalErr: "SQL Error: record not found",
		})
	})

	t.Run("when failed to get artifact info", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						appVersion.ArtifactInfoData = json.RawMessage(`invalid JSON`)
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
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedInternalErr: "invalid character 'i' looking for beginning of value",
		})
	})

	t.Run("when no feature graphic found", func(t *testing.T) {
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
						return featureGraphic, gorm.ErrRecordNotFound
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

	t.Run("when error happens at finding feature graphic", func(t *testing.T) {
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
						return featureGraphic, errors.New("SOME-SQL-ERROR")
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
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when failed to get AWS presigned URL for feature graphic", func(t *testing.T) {
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
						return "", errors.New("SOME-AWS-ERROR")
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
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedInternalErr: "SOME-AWS-ERROR",
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
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedInternalErr: "invalid character 'i' looking for beginning of value",
		})
	})

	t.Run("when failed to get app data from API", func(t *testing.T) {
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
						return &bitrise.AppDetails{}, errors.New("SOME-BITRISE-API-EROR")
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
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedInternalErr: "SOME-BITRISE-API-EROR",
		})
	})

	t.Run("when it's failed to find app settings", func(t *testing.T) {
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

	t.Run("when failed to get android settings from app settings", func(t *testing.T) {
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
						appSettings.AndroidSettingsData = json.RawMessage(`invalid JSON`)
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

	t.Run("when failed to fetch service account file from API", func(t *testing.T) {
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
						return []bitrise.GenericProjectFile{}, errors.New("SOME-BITRISE-API-ERROR")
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
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
				},
			},
			expectedInternalErr: "SOME-BITRISE-API-ERROR",
		})
	})

	t.Run("when no matching service account file found", func(t *testing.T) {
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
						return []bitrise.GenericProjectFile{bitrise.GenericProjectFile{Slug: "not-matching-slug"}}, nil
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

	t.Run("when failed to fetch android keystore file from API", func(t *testing.T) {
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
						return []bitrise.AndroidKeystoreFile{}, errors.New("SOME-BITRISE-API-ERROR")
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.AndroidSettingsData = json.RawMessage(`{"selected_service_account":"service-account-slug","selected_keystore_file":"android-keystore-slug"}`)
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

	t.Run("when no mathcing android keystore file found", func(t *testing.T) {
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
							bitrise.AndroidKeystoreFile{Slug: "not-matching-slug", UserEnvKey: "ANDROID_KEYSTORE"},
						}, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.AndroidSettingsData = json.RawMessage(`{"selected_service_account":"service-account-slug","selected_keystore_file":"android-keystore-slug"}`)
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
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, errors.New("SOME-SQL-ERROR")
					},
				},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when failed to get presigned URL for screenshot", func(t *testing.T) {
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
						if strings.Contains(path, "Apple Watch") {
							return "", errors.New("SOME-AWS-ERROR")
						}
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
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{models.Screenshot{DeviceType: "Apple Watch"}}, nil
					},
				},
			},
			expectedInternalErr: "SOME-AWS-ERROR",
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
					getArtifactsFn: func(apiToken, appSlug, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{}, errors.New("SOME-BITRISE-API-ERROR")
					},
					getArtifactFn: func(apiToken, appSlug, buildSlug, artifactSlug string) (*bitrise.ArtifactShowResponseItemModel, error) {
						return nil, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.AndroidSettingsData = json.RawMessage(`{"selected_service_account":"service-account-slug","selected_keystore_file":"android-keystore-slug"}`)
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
					getArtifactsFn: func(apiToken, appSlug, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{
							bitrise.ArtifactListElementResponseModel{Title: "app.aab", Slug: "test-artifact-slug"},
						}, nil
					},
					getArtifactFn: func(apiToken, appSlug, buildSlug, artifactSlug string) (*bitrise.ArtifactShowResponseItemModel, error) {
						return nil, errors.New("SOME-BITRISE-API-ERROR")
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.AndroidSettingsData = json.RawMessage(`{"selected_service_account":"service-account-slug","selected_keystore_file":"android-keystore-slug"}`)
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

	t.Run("when artifact doesn't have download URL", func(t *testing.T) {
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
					getArtifactsFn: func(apiToken, appSlug, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{
							bitrise.ArtifactListElementResponseModel{Title: "app.aab", Slug: "test-artifact-slug"},
						}, nil
					},
					getArtifactFn: func(apiToken, appSlug, buildSlug, artifactSlug string) (*bitrise.ArtifactShowResponseItemModel, error) {
						return &bitrise.ArtifactShowResponseItemModel{Title: pointers.NewStringPtr("app.aab")}, nil
					},
				},
				AppSettingsService: &testAppSettingsService{
					findFn: func(appSettings *models.AppSettings) (*models.AppSettings, error) {
						appSettings.AndroidSettingsData = json.RawMessage(`{"selected_service_account":"service-account-slug","selected_keystore_file":"android-keystore-slug"}`)
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

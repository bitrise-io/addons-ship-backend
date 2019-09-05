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
	"github.com/bitrise-io/go-utils/envutil"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_AppVersionPublishPostHandler(t *testing.T) {
	httpMethod := "POST"
	url := "/apps/{app-slug}/versions/{version-id}/publish"
	handler := services.AppVersionPublishPostHandler

	testAppVersionID := uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879")
	testTaskIdentifier := uuid.FromStringOrNil("13a94c5d-4609-404e-ae69-c625e93b8b71")

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppVersionService", "PublishTaskService", "BitriseAPI"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppVersionService: &testAppVersionService{
				findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
					return &models.AppVersion{}, nil
				},
			},
			PublishTaskService: &testPublishTaskService{},
			BitriseAPI: &testBitriseAPI{
				getArtifactsFn: func(apiToken string, appSlug string, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
					return []bitrise.ArtifactListElementResponseModel{}, nil
				},
			},
		},
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppVersionID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppVersionService:  &testAppVersionService{},
			PublishTaskService: &testPublishTaskService{},
			BitriseAPI: &testBitriseAPI{
				getArtifactsFn: func(apiToken string, appSlug string, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
					return []bitrise.ArtifactListElementResponseModel{}, nil
				},
			},
		},
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						require.Equal(t, appVersion.ID, testAppVersionID)
						return &models.AppVersion{App: models.App{}, AppStoreInfoData: json.RawMessage(`{}`)}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(apiToken string, appSlug string, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{}, nil
					},
					triggerDENTaskFn: func(bitrise.TaskParams) (*bitrise.TriggerResponse, error) {
						return &bitrise.TriggerResponse{}, nil
					},
				},
				PublishTaskService: &testPublishTaskService{
					createFn: func(*models.PublishTask) (*models.PublishTask, error) {
						return nil, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionPublishResponse{
				Data: &bitrise.TriggerResponse{},
			},
		})
	})

	t.Run("ok - more complex - ios", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AddonHostURL: "http://ship.addon.url",
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						require.Equal(t, appVersion.ID, testAppVersionID)
						return &models.AppVersion{
							App: models.App{
								AppSlug:         "test-app-slug",
								BitriseAPIToken: "bitrise-api-addon-token",
								APIToken:        "addon-access-token",
							},
							Platform:         "ios",
							AppStoreInfoData: json.RawMessage(`{}`),
							BuildSlug:        "test-build-slug",
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(apiToken string, appSlug string, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						require.Equal(t, "bitrise-api-addon-token", apiToken)
						require.Equal(t, "test-app-slug", appSlug)
						require.Equal(t, "test-build-slug", buildSlug)
						return []bitrise.ArtifactListElementResponseModel{
							bitrise.ArtifactListElementResponseModel{
								Slug:  "test-artifact-slug",
								Title: "my-awesome-app.xcarchive.zip",
							},
						}, nil
					},
					triggerDENTaskFn: func(params bitrise.TaskParams) (*bitrise.TriggerResponse, error) {
						require.Equal(t, map[string]string{
							"BITRISE_APP_SLUG":      "test-app-slug",
							"BITRISE_ARTIFACT_SLUG": "test-artifact-slug",
							"BITRISE_BUILD_SLUG":    "test-build-slug",
							"CONFIG_JSON_URL":       "http://ship.addon.url/apps/test-app-slug/versions/de438ddc-98e5-4226-a5f4-fd2d53474879/ios-config",
						}, params.InlineEnvs)
						require.Equal(t, map[string]interface{}{"envs": []bitrise.TaskSecret{
							bitrise.TaskSecret{"BITRISE_ACCESS_TOKEN": "bitrise-api-addon-token"},
							bitrise.TaskSecret{"ADDON_SHIP_APP_ACCESS_TOKEN": "addon-access-token"},
							bitrise.TaskSecret{"SSH_RSA_PRIVATE_KEY": ""},
						}}, params.Secrets)
						require.Equal(t, "http://ship.addon.url/task-webhook", params.WebhookURL)
						require.Equal(t, "resign_archive_app_store", params.Workflow)
						return &bitrise.TriggerResponse{TaskIdentifier: testTaskIdentifier}, nil
					},
				},
				PublishTaskService: &testPublishTaskService{
					createFn: func(publishTask *models.PublishTask) (*models.PublishTask, error) {
						require.Equal(t, testTaskIdentifier, publishTask.TaskID)
						return publishTask, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionPublishResponse{
				Data: &bitrise.TriggerResponse{TaskIdentifier: testTaskIdentifier},
			},
		})
	})

	t.Run("ok - more complex - android", func(t *testing.T) {
		revokeGitUserFn, err := envutil.RevokableSetenv("ANDROID_PUBLISH_WF_GIT_CLONE_USER", "git_user")
		require.NoError(t, err)
		revokeGitPwdFn, err := envutil.RevokableSetenv("ANDROID_PUBLISH_WF_GIT_CLONE_PWD", "git_pwd")
		require.NoError(t, err)

		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AddonHostURL:     "http://ship.addon.url",
				AddonAccessToken: "super-secret-token",
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						require.Equal(t, appVersion.ID, testAppVersionID)
						return &models.AppVersion{
							App: models.App{
								AppSlug:         "test-app-slug",
								BitriseAPIToken: "bitrise-api-addon-token",
								APIToken:        "addon-access-token",
							},
							Platform:         "android",
							AppStoreInfoData: json.RawMessage(`{}`),
							BuildSlug:        "test-build-slug",
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(apiToken string, appSlug string, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						require.Equal(t, "bitrise-api-addon-token", apiToken)
						require.Equal(t, "test-app-slug", appSlug)
						require.Equal(t, "test-build-slug", buildSlug)
						return []bitrise.ArtifactListElementResponseModel{bitrise.ArtifactListElementResponseModel{Slug: "test-artifact-slug"}}, nil
					},
					triggerDENTaskFn: func(params bitrise.TaskParams) (*bitrise.TriggerResponse, error) {
						require.Equal(t, map[string]string{
							"CONFIG_JSON_URL":    "http://ship.addon.url/apps/test-app-slug/versions/de438ddc-98e5-4226-a5f4-fd2d53474879/android-config",
							"GIT_REPOSITORY_URL": "https://git_user:git_pwd@github.com/bitrise-io/addons-ship-bg-worker-task-android",
						}, params.InlineEnvs)
						require.Equal(t, "http://ship.addon.url/task-webhook", params.WebhookURL)
						require.Equal(t, "resign_android", params.Workflow)
						require.Equal(t, map[string]interface{}{"envs": []bitrise.TaskSecret{
							bitrise.TaskSecret{"ADDON_SHIP_ACCESS_TOKEN": "super-secret-token"},
							bitrise.TaskSecret{"ADDON_SHIP_APP_ACCESS_TOKEN": "addon-access-token"},
							bitrise.TaskSecret{"BITRISE_ACCESS_TOKEN": "bitrise-api-addon-token"},
						}}, params.Secrets)
						return &bitrise.TriggerResponse{TaskIdentifier: testTaskIdentifier}, nil
					},
				},
				PublishTaskService: &testPublishTaskService{
					createFn: func(publishTask *models.PublishTask) (*models.PublishTask, error) {
						require.Equal(t, testTaskIdentifier, publishTask.TaskID)
						return publishTask, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionPublishResponse{
				Data: &bitrise.TriggerResponse{TaskIdentifier: testTaskIdentifier},
			},
		})
		require.NoError(t, revokeGitUserFn())
		require.NoError(t, revokeGitPwdFn())
	})

	t.Run("when app version not found in database", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						require.Equal(t, appVersion.ID, testAppVersionID)
						return nil, gorm.ErrRecordNotFound
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(apiToken string, appSlug string, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{}, nil
					},
					triggerDENTaskFn: func(bitrise.TaskParams) (*bitrise.TriggerResponse, error) {
						return nil, nil
					},
				},
				PublishTaskService: &testPublishTaskService{},
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Not Found"},
		})
	})

	t.Run("when error happens at finding app version", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						require.Equal(t, appVersion.ID, testAppVersionID)
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(apiToken string, appSlug string, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{}, nil
					},
					triggerDENTaskFn: func(bitrise.TaskParams) (*bitrise.TriggerResponse, error) {
						return nil, nil
					},
				},
				PublishTaskService: &testPublishTaskService{},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when error happens at getting artifact data from API", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						require.Equal(t, appVersion.ID, testAppVersionID)
						return &models.AppVersion{App: models.App{}, AppStoreInfoData: json.RawMessage(`{}`)}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(apiToken string, appSlug string, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{}, errors.New("SOME-BITRISE-API-ERROR")
					},
					triggerDENTaskFn: func(bitrise.TaskParams) (*bitrise.TriggerResponse, error) {
						return nil, nil
					},
				},
				PublishTaskService: &testPublishTaskService{},
			},
			expectedInternalErr: "SOME-BITRISE-API-ERROR",
		})
	})

	t.Run("when error happens at triggering DEN task", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						require.Equal(t, appVersion.ID, testAppVersionID)
						return &models.AppVersion{App: models.App{}, AppStoreInfoData: json.RawMessage(`{}`)}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(apiToken string, appSlug string, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{}, nil
					},
					triggerDENTaskFn: func(bitrise.TaskParams) (*bitrise.TriggerResponse, error) {
						return nil, errors.New("SOME-BITRISE-API-ERROR")
					},
				},
				PublishTaskService: &testPublishTaskService{},
			},
			expectedInternalErr: "SOME-BITRISE-API-ERROR",
		})
	})

	t.Run("when error happens creating publish task object", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						require.Equal(t, appVersion.ID, testAppVersionID)
						return &models.AppVersion{App: models.App{}, AppStoreInfoData: json.RawMessage(`{}`)}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(apiToken string, appSlug string, buildSlug string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{}, nil
					},
					triggerDENTaskFn: func(bitrise.TaskParams) (*bitrise.TriggerResponse, error) {
						return &bitrise.TriggerResponse{}, nil
					},
				},
				PublishTaskService: &testPublishTaskService{
					createFn: func(*models.PublishTask) (*models.PublishTask, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})
}

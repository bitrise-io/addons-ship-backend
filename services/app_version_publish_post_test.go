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
	"github.com/c2fo/testify/require"
	uuid "github.com/satori/go.uuid"
)

func Test_AppVersionPublishPostHandler(t *testing.T) {
	httpMethod := "POST"
	url := "/apps/{app-slug}/versions/{version-id}/publish"
	handler := services.AppVersionPublishPostHandler

	testAppVersionID := uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879")

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppVersionService", "BitriseAPI"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppVersionService: &testAppVersionService{
				findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
					return &models.AppVersion{}, nil
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
					getArtifactDataFn: func(string, string, string) (*bitrise.ArtifactData, error) {
						return &bitrise.ArtifactData{}, nil
					},
					triggerDENTaskFn: func(bitrise.TaskParams) (*bitrise.TriggerResponse, error) {
						return nil, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   services.AppVersionPublishResponse{},
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
						require.Equal(t, appVersion.ID, testAppVersionID)
						return &models.AppVersion{
							App: models.App{
								AppSlug:         "test-app-slug",
								BitriseAPIToken: "bitrise-api-addon-token",
							},
							AppStoreInfoData: json.RawMessage(`{}`),
							BuildSlug:        "test-build-slug",
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactDataFn: func(string, string, string) (*bitrise.ArtifactData, error) {
						return &bitrise.ArtifactData{Slug: "test-artifact-slug"}, nil
					},
					triggerDENTaskFn: func(params bitrise.TaskParams) (*bitrise.TriggerResponse, error) {
						require.Equal(t, `{"BITRISE_APP_SLUG":"test-app-slug","BITRISE_ARTIFACT_SLUG":"test-artifact-slug","BITRISE_BUILD_SLUG":"test-build-slug"}`, params.InlineEnvs)
						require.Equal(t, `{"BITRISE_ACCESS_TOKEN":"bitrise-api-addon-token"}`, params.Secrets)
						return &bitrise.TriggerResponse{}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionPublishResponse{
				Data: &bitrise.TriggerResponse{},
			},
		})
	})
}

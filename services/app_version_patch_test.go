package services_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/c2fo/testify/require"
	uuid "github.com/satori/go.uuid"
)

func Test_AppVersionsPatchHandler(t *testing.T) {
	httpMethod := "PATCH"
	url := "/apps/{app-slug}/version{version-id}"
	handler := services.AppVersionPatchHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppVersionService"}, ControllerTestCase{
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

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						require.Equal(t, appVersion.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
						return &models.AppVersion{App: models.App{}, AppStoreInfoData: json.RawMessage(`{}`)}, nil
					},
				},
			},
			requestBody:        `{}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionPatchResponse{
				Data: services.AppVersionPatchResponseData{AppVersion: &models.AppVersion{}},
			},
		})
	})
}

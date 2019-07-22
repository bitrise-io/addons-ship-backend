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
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_AppVersionsGetHandler(t *testing.T) {
	httpMethod := "GET"
	url := "/apps/{app-slug}/app-versions"
	handler := services.AppVersionsGetHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppVersionService"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppVersionService: &testAppVersionService{},
		},
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppVersionService: &testAppVersionService{},
		},
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findAllFn: func(app *models.App, filterParams map[string]interface{}) ([]models.AppVersion, error) {
						require.Equal(t, app.ID.String(), "211afc15-127a-40f9-8cbe-1dadc1f86cdf")
						require.Equal(t, filterParams, map[string]interface{}{})
						return []models.AppVersion{}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionsGetResponse{
				Data: []services.AppVersionsGetResponseElement{},
			},
		})
	})

	t.Run("ok - more complex", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findAllFn: func(app *models.App, filterParams map[string]interface{}) ([]models.AppVersion, error) {
						return []models.AppVersion{
							models.AppVersion{
								Platform:         "ios",
								ArtifactInfoData: json.RawMessage(`{"version":"v1.0"}`),
							},
							models.AppVersion{
								Platform:         "android",
								ArtifactInfoData: json.RawMessage(`{"version":"v1.12"}`),
							},
						}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionsGetResponse{
				Data: []services.AppVersionsGetResponseElement{
					services.AppVersionsGetResponseElement{
						AppVersion: models.AppVersion{
							Platform: "ios",
						},
						Version: "v1.0",
					},
					services.AppVersionsGetResponseElement{
						AppVersion: models.AppVersion{
							Platform: "android",
						},
						Version: "v1.12",
					},
				},
			},
		})
	})

	t.Run("ok - with platform filter", func(t *testing.T) {
		urlWithFilter := url + "?platform=ios"
		performControllerTest(t, httpMethod, urlWithFilter, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findAllFn: func(app *models.App, filterParams map[string]interface{}) ([]models.AppVersion, error) {
						require.Equal(t, app.ID.String(), "211afc15-127a-40f9-8cbe-1dadc1f86cdf")
						require.Equal(t, filterParams, map[string]interface{}{
							"platform": "ios",
						})
						return []models.AppVersion{
							models.AppVersion{
								ArtifactInfoData: json.RawMessage(`{"version":"v1.0"}`),
								Platform:         "ios",
							},
						}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionsGetResponse{
				Data: []services.AppVersionsGetResponseElement{
					services.AppVersionsGetResponseElement{
						AppVersion: models.AppVersion{
							Platform: "ios",
						},
						Version: "v1.0",
					},
				},
			},
		})
	})

	t.Run("error - unexpected error in database", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findAllFn: func(app *models.App, filterParams map[string]interface{}) ([]models.AppVersion, error) {
						return []models.AppVersion{}, errors.New("SOME-SQL-ERROR")
					},
				},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})
}

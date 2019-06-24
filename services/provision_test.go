package services_test

import (
	"net/http"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func Test_ProvisionHandler(t *testing.T) {
	httpMethod := "POST"
	url := "/provision"
	handler := services.ProvisionHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppService"}, ControllerTestCase{
		env: &env.AppEnv{
			AppService: &testAppService{},
		},
	})

	t.Run("ok", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return app, nil
					},
				},
			},
			requestBody:        `{}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.ProvisionPostResponse{
				Envs: []services.Env{
					services.Env{Key: "ADDON_SHIP_API_URL"},
					services.Env{Key: "ADDON_SHIP_API_TOKEN"},
				},
			},
		})
	})

	t.Run("ok when app exists", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						require.Equal(t, "test-app-slug", app.AppSlug)
						return app, nil
					},
				},
			},
			requestBody:        `{"app_slug":"test-app-slug","bitrise_api_token":"test-bitrise-api-token","plan":"free"}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.ProvisionPostResponse{
				Envs: []services.Env{
					services.Env{Key: "ADDON_SHIP_API_URL"},
					services.Env{Key: "ADDON_SHIP_API_TOKEN"},
				},
			},
		})
	})

	t.Run("ok when app not exists", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			env: &env.AppEnv{
				AddonHostURL: "http://ship.addon.url",
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						require.Equal(t, "test-app-slug", app.AppSlug)
						return nil, gorm.ErrRecordNotFound
					},
					createFn: func(app *models.App) (*models.App, error) {
						require.Equal(t, "test-app-slug", app.AppSlug)
						require.Equal(t, "test-bitrise-api-token", app.BitriseAPIToken)
						require.Equal(t, "free", app.Plan)
						require.NotEmpty(t, app.APIToken)
						app.APIToken = "test-api-token"
						return app, nil
					},
				},
			},
			requestBody:        `{"app_slug":"test-app-slug","bitrise_api_token":"test-bitrise-api-token","plan":"free"}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.ProvisionPostResponse{
				Envs: []services.Env{
					services.Env{Key: "ADDON_SHIP_API_URL", Value: "http://ship.addon.url"},
					services.Env{Key: "ADDON_SHIP_API_TOKEN", Value: "test-api-token"},
				},
			},
		})
	})

	t.Run("when request body is invalid", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return app, nil
					},
				},
			},
			requestBody:        `invalid JSON`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Invalid request body, JSON decode failed"},
		})
	})

	t.Run("when database error happest at find", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
			},
			requestBody:         `{}`,
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when database error happens at create", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return nil, gorm.ErrRecordNotFound
					},
					createFn: func(app *models.App) (*models.App, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
			},
			requestBody:         `{}`,
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})
}

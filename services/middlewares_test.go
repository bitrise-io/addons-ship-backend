package services_test

import (
	"net/http"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	"github.com/bitrise-io/api-utils/middleware"
	"github.com/bitrise-io/api-utils/providers"
)

func Test_AuthorizedAppMiddleware(t *testing.T) {
	middleware.PerformTest(t, "GET", "/...", middleware.TestCase{
		RequestHeaders: map[string]string{
			"Authorization": "token ADDON_AUTH_TOKEN",
		},
		ExpectedStatus: http.StatusOK,
		ExpectedResponse: map[string]interface{}{
			"message": "Success",
		},
		Middleware: services.AuthorizedAppMiddleware(&env.AppEnv{
			RequestParams: &providers.RequestParamsProviderMock{
				Params: map[string]string{
					"app-slug": "test_app_slug",
				},
			},
			AppService: &testAppService{
				findFn: func(app *models.App) (*models.App, error) {
					return app, nil
				},
			},
		}),
	})
}

func Test_AuthorizationAppVersionMiddleware(t *testing.T) {
	middleware.PerformTest(t, "GET", "/...", middleware.TestCase{
		RequestHeaders: map[string]string{
			"Authorization": "token ADDON_AUTH_TOKEN",
		},
		ExpectedStatus: http.StatusOK,
		ExpectedResponse: map[string]interface{}{
			"message": "Success",
		},
		Middleware: services.AuthorizationAppVersionMiddleware(&env.AppEnv{
			RequestParams: &providers.RequestParamsProviderMock{
				Params: map[string]string{
					"app-slug":   "test_app_slug",
					"version-id": "de438ddc-98e5-4226-a5f4-fd2d53474879",
				},
			},
			AppService: &testAppService{
				findFn: func(app *models.App) (*models.App, error) {
					return app, nil
				},
			},
			AppVersionService: &testAppVersionService{
				findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
					return appVersion, nil
				},
			},
		}),
	})
}

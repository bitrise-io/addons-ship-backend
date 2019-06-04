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
			RequestParams: &providers.RequestParamsMock{
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

func Test_AuthorizedAppVersionMiddleware(t *testing.T) {
	middleware.PerformTest(t, "GET", "/...", middleware.TestCase{
		RequestHeaders: map[string]string{
			"Authorization": "token ADDON_AUTH_TOKEN",
		},
		ExpectedStatus: http.StatusOK,
		ExpectedResponse: map[string]interface{}{
			"message": "Success",
		},
		Middleware: services.AuthorizedAppVersionMiddleware(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
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

func Test_AuthorizedAppVersionScreenshotMiddleware(t *testing.T) {
	middleware.PerformTest(t, "GET", "/...", middleware.TestCase{
		RequestHeaders: map[string]string{
			"Authorization": "token ADDON_AUTH_TOKEN",
		},
		ExpectedStatus: http.StatusOK,
		ExpectedResponse: map[string]interface{}{
			"message": "Success",
		},
		Middleware: services.AuthorizedAppVersionScreenshotMiddleware(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
				Params: map[string]string{
					"app-slug":      "test_app_slug",
					"version-id":    "de438ddc-98e5-4226-a5f4-fd2d53474879",
					"screenshot-id": "abcd1234-5678-ef12-9012-fd2d53474123",
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
			ScreenshotService: &testScreenshotService{
				findFn: func(screenshot *models.Screenshot) (*models.Screenshot, error) {
					return screenshot, nil
				},
			},
		}),
	})
}

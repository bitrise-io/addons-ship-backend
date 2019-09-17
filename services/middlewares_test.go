package services_test

import (
	"net/http"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	"github.com/bitrise-io/api-utils/middleware"
	"github.com/bitrise-io/api-utils/providers"
	"github.com/bitrise-io/api-utils/security"
	"github.com/bitrise-io/go-crypto/crypto"
	"github.com/bitrise-io/go-utils/envutil"
	"github.com/c2fo/testify/require"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

func Test_AuthenticateForProvisioning(t *testing.T) {
	middleware.PerformTest(t, "GET", "/...", middleware.TestCase{
		RequestHeaders: map[string]string{
			"Authentication": "ADDON_AUTH_TOKEN",
		},
		ExpectedStatus: http.StatusOK,
		ExpectedResponse: map[string]interface{}{
			"message": "Success",
		},
		Middleware: services.AuthenticateForProvisioning(&env.AppEnv{
			AddonAccessToken: "ADDON_AUTH_TOKEN",
		}),
	})
}

func Test_AuthenticateForDeprovisioning(t *testing.T) {
	middleware.PerformTest(t, "GET", "/...", middleware.TestCase{
		RequestHeaders: map[string]string{
			"Authentication": "ADDON_AUTH_TOKEN",
		},
		ExpectedStatus: http.StatusOK,
		ExpectedResponse: map[string]interface{}{
			"message": "Success",
		},
		Middleware: services.AuthenticateForDeprovisioning(&env.AppEnv{
			AddonAccessToken: "ADDON_AUTH_TOKEN",
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
			JWTService: &security.JWTMock{
				VerifyFn: func(token string) (bool, error) {
					return true, nil
				},
				GetTokenFn: func(token string) (interface{}, error) {
					return "auth-token-from-jwt", nil
				},
			},
		}),
	})
}

func Test_AuthorizedAppResourceMiddleware(t *testing.T) {
	middleware.PerformTest(t, "GET", "/...", middleware.TestCase{
		RequestHeaders: map[string]string{
			"Authorization": "token ADDON_AUTH_TOKEN",
		},
		ExpectedStatus: http.StatusOK,
		ExpectedResponse: map[string]interface{}{
			"message": "Success",
		},
		Middleware: services.AuthorizedAppResourceMiddleware(&env.AppEnv{
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
			JWTService: &security.JWTMock{
				VerifyFn: func(token string) (bool, error) {
					return true, nil
				},
				GetTokenFn: func(token string) (interface{}, error) {
					return "auth-token-from-jwt", nil
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
			JWTService: &security.JWTMock{
				VerifyFn: func(token string) (bool, error) {
					return true, nil
				},
				GetTokenFn: func(token string) (interface{}, error) {
					return "auth-token-from-jwt", nil
				},
			},
		}),
	})
}

func Test_AuthorizeForWebhookHandling(t *testing.T) {
	revokeFn, err := envutil.RevokableSetenv("BITRISE_DEN_WEBHOOK_SECRET", "secret-token")
	require.NoError(t, err)

	middleware.PerformTest(t, "GET", "/...", middleware.TestCase{
		RequestHeaders: map[string]string{
			"Bitrise-Den-Webhook-Secret": "secret-token",
		},
		RequestBody: services.WebhookPayload{
			TaskID: uuid.FromStringOrNil("cb8ddaf5-e6f9-470f-b84e-8bc9a0cbf78a"),
		},
		ExpectedStatus: http.StatusOK,
		ExpectedResponse: map[string]interface{}{
			"message": "Success",
		},
		Middleware: services.AuthorizeForWebhookHandling(&env.AppEnv{
			PublishTaskService: &testPublishTaskService{
				findFn: func(publishTask *models.PublishTask) (*models.PublishTask, error) {
					require.Equal(t, uuid.FromStringOrNil("cb8ddaf5-e6f9-470f-b84e-8bc9a0cbf78a"), publishTask.TaskID)
					return publishTask, nil
				},
			},
		}),
	})
	require.NoError(t, revokeFn())
}

func Test_AuthorizeForAppContactEmailConfirmationHandling(t *testing.T) {
	middleware.PerformTest(t, "PATCH", "/...", middleware.TestCase{
		RequestBody:    map[string]string{"confirmation_token": "5om3-r4nd0m-5tr1ng"},
		ExpectedStatus: http.StatusOK,
		ExpectedResponse: map[string]interface{}{
			"message": "Success",
		},
		Middleware: services.AuthorizeForAppContactEmailConfirmationHandling(&env.AppEnv{
			AppContactService: &testAppContactService{
				findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
					require.NotNil(t, appContact.ConfirmationToken)
					require.Equal(t, "5om3-r4nd0m-5tr1ng", *appContact.ConfirmationToken)
					appContact.ID = uuid.FromStringOrNil("8a230385-0113-4cf3-a9c6-469a313e587a")
					return appContact, nil
				},
			},
		}),
	})
}

func Test_AuthorizedAppContactMiddleware(t *testing.T) {
	middleware.PerformTest(t, "GET", "/...", middleware.TestCase{
		RequestHeaders: map[string]string{
			"Authorization": "token ADDON_AUTH_TOKEN",
		},
		ExpectedStatus: http.StatusOK,
		ExpectedResponse: map[string]interface{}{
			"message": "Success",
		},
		Middleware: services.AuthorizedAppContactMiddleware(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
				Params: map[string]string{
					"app-slug":   "test_app_slug",
					"contact-id": "de438ddc-98e5-4226-a5f4-fd2d53474879",
				},
			},
			AppService: &testAppService{
				findFn: func(app *models.App) (*models.App, error) {
					return app, nil
				},
			},
			AppContactService: &testAppContactService{
				findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
					return appContact, nil
				},
			},
			JWTService: &security.JWTMock{
				VerifyFn: func(token string) (bool, error) {
					return true, nil
				},
				GetTokenFn: func(token string) (interface{}, error) {
					return "auth-token-from-jwt", nil
				},
			},
		}),
	})
}

func Test_AuthorizedBuildWebhookMiddleware(t *testing.T) {
	revokeFn, err := envutil.RevokableSetenv("APP_WEBHOOK_SECRET_ENCRYPT_KEY", "06042e86a7bd421c642c8c3e4ab13840")
	require.NoError(t, err)
	middleware.PerformTest(t, "POST", "/...", middleware.TestCase{
		RequestHeaders: map[string]string{
			"Bitrise-App-Id":         "test_app_slug",
			"Bitrise-Hook-Signature": "sha256=0d86929661b1c7b216ca6a7ef4abe740ee6dc07d4afc2f21d78c888235d88713",
		},
		ExpectedStatus: http.StatusOK,
		ExpectedResponse: map[string]interface{}{
			"message": "Success",
		},
		Middleware: services.AuthorizedBuildWebhookMiddleware(&env.AppEnv{
			AppService: &testAppService{
				findFn: func(app *models.App) (*models.App, error) {
					app.ID = uuid.NewV4()
					iv, err := crypto.GenerateIV()
					require.NoError(t, err)
					encryptedSecret, err := crypto.AES256GCMCipher("my-super-secret", iv, "06042e86a7bd421c642c8c3e4ab13840")
					require.NoError(t, err)

					app.EncryptedSecret = encryptedSecret
					app.EncryptedSecretIV = iv
					return app, nil
				},
			},
		}),
		RequestBody: map[string]string{"app_slug": "test-app-slug"},
	})
	require.NoError(t, revokeFn())
}

func Test_AuthenticatedForLoginMiddleware(t *testing.T) {
	testLogger, err := zap.NewDevelopment()
	require.NoError(t, err)
	middleware.PerformTest(t, "POST", "/...", middleware.TestCase{
		Middleware: services.AuthenticatedForLoginMiddleware(&env.AppEnv{
			Logger: testLogger,
			AppService: &testAppService{
				findFn: func(app *models.App) (*models.App, error) {
					app.ID = uuid.NewV4()
					return app, nil
				},
			},
			SsoTokenVerifier: &security.SsoTokenVerifierMock{
				VerifyFn: func(timestamp, ssoToken, appSlug string) (bool, error) {
					return true, nil
				},
			},
		}),
		RequestBody:    map[string]string{"app_slug": "test-app-slug"},
		ExpectedStatus: http.StatusOK,
		ExpectedResponse: map[string]interface{}{
			"message": "Success",
		},
	})
}

package services_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	"github.com/bitrise-io/api-utils/security"
	"github.com/bitrise-io/go-utils/envutil"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Test_AuthenticateWithAddonAccessTokenHandlerFunc(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		performAuthenticationTest(t, "GET", "...", AuthenticationTestCase{
			requestHeaders: map[string]string{
				"Authentication": "test-auth-token",
			},
			env:                &env.AppEnv{AddonHostURL: "http://ship.addon.url", AddonAccessToken: "test-auth-token"},
			authHandlerFunc:    services.AuthenticateWithAddonAccessTokenHandlerFunc,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"ok"}` + "\n",
		})
	})

	t.Run("when no Authentication header is provided", func(t *testing.T) {
		performAuthenticationTest(t, "GET", "...", AuthenticationTestCase{
			requestHeaders:     map[string]string{},
			env:                &env.AppEnv{AddonHostURL: "http://ship.addon.url", AddonAccessToken: "test-auth-token"},
			authHandlerFunc:    services.AuthenticateWithAddonAccessTokenHandlerFunc,
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"message":"Unauthorized"}` + "\n",
		})
	})

	t.Run("when Authentication header has empty value", func(t *testing.T) {
		performAuthenticationTest(t, "GET", "...", AuthenticationTestCase{
			requestHeaders: map[string]string{
				"Authentication": "",
			},
			env:                &env.AppEnv{AddonHostURL: "http://ship.addon.url", AddonAccessToken: "test-auth-token"},
			authHandlerFunc:    services.AuthenticateWithAddonAccessTokenHandlerFunc,
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"message":"Unauthorized"}` + "\n",
		})
	})

	t.Run("when no addon token is set application level", func(t *testing.T) {
		performAuthenticationTest(t, "GET", "...", AuthenticationTestCase{
			requestHeaders: map[string]string{
				"Authentication": "test-auth-token",
			},
			env:                &env.AppEnv{AddonHostURL: "http://ship.addon.url"},
			authHandlerFunc:    services.AuthenticateWithAddonAccessTokenHandlerFunc,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"message":"Internal Server Error"}` + "\n",
		})
	})

	t.Run("when provided auth token does not match", func(t *testing.T) {
		performAuthenticationTest(t, "GET", "...", AuthenticationTestCase{
			requestHeaders: map[string]string{
				"Authentication": "invalid-token",
			},
			env:                &env.AppEnv{AddonHostURL: "http://ship.addon.url", AddonAccessToken: "test-auth-token"},
			authHandlerFunc:    services.AuthenticateWithAddonAccessTokenHandlerFunc,
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"message":"Unauthorized"}` + "\n",
		})
	})
}

func Test_AuthenticateWithDENSecretHandlerFunc(t *testing.T) {
	for _, tc := range []struct {
		testName           string
		denWebhookSecret   string
		requestHeaders     map[string]string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			testName:           "ok",
			denWebhookSecret:   "test-auth-key",
			requestHeaders:     map[string]string{"Bitrise-Den-Webhook-Secret": "test-auth-key"},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"ok"}` + "\n",
		},
		{
			testName:           "when value in request header is empty",
			denWebhookSecret:   "test-auth-key",
			requestHeaders:     map[string]string{"Bitrise-Den-Webhook-Secret": ""},
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"message":"Unauthorized"}` + "\n",
		},
		{
			testName:           "when no den webhook secret is set in envs",
			denWebhookSecret:   "",
			requestHeaders:     map[string]string{"Bitrise-Den-Webhook-Secret": "test-auth-key"},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"message":"Internal Server Error"}` + "\n",
		},
		{
			testName:           "invalid secret provided in request",
			denWebhookSecret:   "test-auth-key",
			requestHeaders:     map[string]string{"Bitrise-Den-Webhook-Secret": "some-totally-different-key"},
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"message":"Unauthorized"}` + "\n",
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			revokeFn, err := envutil.RevokableSetenv("BITRISE_DEN_WEBHOOK_SECRET", tc.denWebhookSecret)
			require.NoError(t, err)

			performAuthenticationTest(t, "GET", "...", AuthenticationTestCase{
				requestHeaders:     tc.requestHeaders,
				authHandlerFunc:    services.AuthenticateWithDENSecretHandlerFunc,
				expectedStatusCode: tc.expectedStatusCode,
				expectedBody:       tc.expectedBody,
			})

			require.NoError(t, revokeFn())
		})
	}
}

func Test_AuthenticateWithSSOTokenHandlerFunc(t *testing.T) {
	reqTimestamp := fmt.Sprintf("%d", time.Now().Unix())

	t.Run("ok", func(t *testing.T) {
		performAuthenticationTest(t, "POST", "...", AuthenticationTestCase{
			requestFormValues: map[string]string{
				"timestamp": reqTimestamp,
				"token":     "sha256-request-sso-token",
				"app_slug":  "test-app-slug",
			},
			env: &env.AppEnv{
				Logger: zap.NewNop(),
				SsoTokenVerifier: &security.SsoTokenVerifierMock{
					VerifyFn: func(timestamp, ssoToken, appSlug string) (bool, error) {
						require.Equal(t, reqTimestamp, timestamp)
						require.Equal(t, "sha256-request-sso-token", ssoToken)
						require.Equal(t, "test-app-slug", appSlug)
						return true, nil
					},
				},
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						require.Equal(t, "test-app-slug", app.AppSlug)
						return app, nil
					},
				},
			},
			authHandlerFunc:    services.AuthenticateWithSSOTokenHandlerFunc,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"ok"}` + "\n",
		})
	})

	t.Run("when app service is not provided", func(t *testing.T) {
		performAuthenticationTest(t, "POST", "...", AuthenticationTestCase{
			requestFormValues: map[string]string{
				"timestamp": reqTimestamp,
				"token":     "sha256-request-sso-token",
				"app_slug":  "test-app-slug",
			},
			env: &env.AppEnv{
				Logger: zap.NewNop(),
				SsoTokenVerifier: &security.SsoTokenVerifierMock{
					VerifyFn: func(timestamp, ssoToken, appSlug string) (bool, error) {
						return true, nil
					},
				},
				AppService: nil,
			},
			authHandlerFunc:    services.AuthenticateWithSSOTokenHandlerFunc,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"message":"Internal Server Error"}` + "\n",
		})
	})

	t.Run("when sso verifies is not defined", func(t *testing.T) {
		performAuthenticationTest(t, "POST", "...", AuthenticationTestCase{
			requestFormValues: map[string]string{
				"timestamp": reqTimestamp,
				"token":     "sha256-request-sso-token",
				"app_slug":  "test-app-slug",
			},
			env: &env.AppEnv{
				Logger:           zap.NewNop(),
				SsoTokenVerifier: nil,
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return app, nil
					},
				},
			},
			authHandlerFunc:    services.AuthenticateWithSSOTokenHandlerFunc,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"message":"Internal Server Error"}` + "\n",
		})
	})

	t.Run("when error happened at verification", func(t *testing.T) {
		performAuthenticationTest(t, "POST", "...", AuthenticationTestCase{
			requestFormValues: map[string]string{
				"timestamp": reqTimestamp,
				"token":     "sha256-request-sso-token",
				"app_slug":  "test-app-slug",
			},
			env: &env.AppEnv{
				Logger: zap.NewNop(),
				SsoTokenVerifier: &security.SsoTokenVerifierMock{
					VerifyFn: func(timestamp, ssoToken, appSlug string) (bool, error) {
						return false, errors.New("SOME-VERIFICATION-ERROR")
					},
				},
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return app, nil
					},
				},
			},
			authHandlerFunc:    services.AuthenticateWithSSOTokenHandlerFunc,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"message":"Internal Server Error"}` + "\n",
		})
	})

	t.Run("when verification returns false", func(t *testing.T) {
		performAuthenticationTest(t, "POST", "...", AuthenticationTestCase{
			requestFormValues: map[string]string{
				"timestamp": reqTimestamp,
				"token":     "sha256-request-sso-token",
				"app_slug":  "test-app-slug",
			},
			env: &env.AppEnv{
				Logger: zap.NewNop(),
				SsoTokenVerifier: &security.SsoTokenVerifierMock{
					VerifyFn: func(timestamp, ssoToken, appSlug string) (bool, error) {
						return false, nil
					},
				},
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						require.Equal(t, "test-app-slug", app.AppSlug)
						return app, nil
					},
				},
			},
			authHandlerFunc:    services.AuthenticateWithSSOTokenHandlerFunc,
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       `{"message":"Not Found"}` + "\n",
		})
	})

	t.Run("when app cannot be found", func(t *testing.T) {
		performAuthenticationTest(t, "POST", "...", AuthenticationTestCase{
			requestFormValues: map[string]string{
				"timestamp": reqTimestamp,
				"token":     "sha256-request-sso-token",
				"app_slug":  "test-app-slug",
			},
			env: &env.AppEnv{
				Logger: zap.NewNop(),
				SsoTokenVerifier: &security.SsoTokenVerifierMock{
					VerifyFn: func(timestamp, ssoToken, appSlug string) (bool, error) {
						require.Equal(t, reqTimestamp, timestamp)
						require.Equal(t, "sha256-request-sso-token", ssoToken)
						require.Equal(t, "test-app-slug", appSlug)
						return true, nil
					},
				},
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						require.Equal(t, "test-app-slug", app.AppSlug)
						return app, gorm.ErrRecordNotFound
					},
				},
			},
			authHandlerFunc:    services.AuthenticateWithSSOTokenHandlerFunc,
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       `{"message":"Not Found"}` + "\n",
		})
	})

	t.Run("when error happens at finding app", func(t *testing.T) {
		performAuthenticationTest(t, "POST", "...", AuthenticationTestCase{
			requestFormValues: map[string]string{
				"timestamp": reqTimestamp,
				"token":     "sha256-request-sso-token",
				"app_slug":  "test-app-slug",
			},
			env: &env.AppEnv{
				Logger: zap.NewNop(),
				SsoTokenVerifier: &security.SsoTokenVerifierMock{
					VerifyFn: func(timestamp, ssoToken, appSlug string) (bool, error) {
						require.Equal(t, reqTimestamp, timestamp)
						require.Equal(t, "sha256-request-sso-token", ssoToken)
						require.Equal(t, "test-app-slug", appSlug)
						return true, nil
					},
				},
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						require.Equal(t, "test-app-slug", app.AppSlug)
						return app, errors.New("SOME-SQL-ERROR")
					},
				},
			},
			authHandlerFunc:    services.AuthenticateWithSSOTokenHandlerFunc,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"message":"Internal Server Error"}` + "\n",
		})
	})
}

package services_test

import (
	"net/http"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/services"
	"github.com/bitrise-io/go-utils/envutil"
	"github.com/c2fo/testify/require"
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

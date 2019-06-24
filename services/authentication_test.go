package services_test

import (
	"net/http"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/services"
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

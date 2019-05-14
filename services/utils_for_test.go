package services_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/c2fo/testify/require"
)

type ControllerTestCase struct {
	requestBody         string
	requestHeaders      map[string]string
	expectedStatusCode  int
	expectedResponse    interface{}
	expectedInternalErr string

	contextElements map[ctxpkg.RequestContextKey]interface{}

	env *env.AppEnv
}

func performControllerTest(t *testing.T,
	httpMethod, url string,
	handler func(*env.AppEnv, http.ResponseWriter, *http.Request) error,
	tc ControllerTestCase,
) {
	t.Helper() // This call silences this function in error reports. See: https://blog.golang.org/go1.9

	r, err := http.NewRequest(httpMethod, url, bytes.NewBuffer([]byte(tc.requestBody)))
	require.NoError(t, err)

	if len(tc.requestHeaders) > 0 {
		for headerKey, headerValue := range tc.requestHeaders {
			r.Header.Set(headerKey, headerValue)
		}
	}

	for key, val := range tc.contextElements {
		r = r.WithContext(context.WithValue(r.Context(), key, val))
	}

	rr := httptest.NewRecorder()
	internalServerError := handler(tc.env, rr, r)

	var expectedBytes []byte
	if tc.expectedResponse != nil {
		expectedBytes, err := json.Marshal(tc.expectedResponse)
		require.NoError(t, err)
		require.Equal(t, string(expectedBytes), strings.Trim(rr.Body.String(), "\n"))
	}

	if tc.expectedInternalErr != "" {
		require.EqualError(t, internalServerError, tc.expectedInternalErr,
			"Expected internal err: %s | Request Body: %s | Response Code: %d, Expected Response Body: %s | Got Body: %s", tc.expectedInternalErr, tc.requestBody, rr.Code, string(expectedBytes), rr.Body.String())
	} else {
		require.NoError(t, internalServerError)
		if tc.expectedStatusCode != 0 {
			require.Equal(t, tc.expectedStatusCode, rr.Code,
				"Expected body: %s | Got body: %s", string(expectedBytes), rr.Body.String())
		}
	}
}

func behavesAsServiceCravingHandler(t *testing.T, method, url string, handler func(*env.AppEnv, http.ResponseWriter, *http.Request) error, serviceNames []string, baseCT ControllerTestCase) {
	t.Run("behaves as service craving handler", func(t *testing.T) {
		for _, sn := range serviceNames {
			controllerTestCase := baseCT
			if sn == "AppService" {
				controllerTestCase.env.AppService = nil
				controllerTestCase.expectedInternalErr = "No App Service defined for handler"
			} else if sn == "AppVersionService" {
				controllerTestCase.env.AppVersionService = nil
				controllerTestCase.expectedInternalErr = "No App Version Service defined for handler"
			} else if sn == "RequestParams" {
				controllerTestCase.env.RequestParams = nil
				controllerTestCase.expectedInternalErr = "No RequestParams defined for handler"
			} else {

				t.Fatalf("Invalid context element name defined: %s", sn)
			}
			performControllerTest(t, method, url, handler, controllerTestCase)
		}
	})
}

func behavesAsContextCravingHandler(t *testing.T, method, url string, handler func(*env.AppEnv, http.ResponseWriter, *http.Request) error, contextKeys []ctxpkg.RequestContextKey, baseCT ControllerTestCase) {
	t.Run("behaves as context craving handler", func(t *testing.T) {
		for _, ck := range contextKeys {
			controllerTestCase := baseCT
			if ck == services.ContextKeyAuthorizedAppID {
				controllerTestCase.contextElements[ck] = nil
				controllerTestCase.expectedInternalErr = "Authorized App ID not found in Context"
			} else {

				t.Fatalf("Invalid context element name defined: %s", ck)
			}
			performControllerTest(t, method, url, handler, controllerTestCase)
		}
	})
}

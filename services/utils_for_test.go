package services_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/c2fo/testify/require"
)

type ControllerTestCase struct {
	requestBody              string
	requestHeaders           map[string]string
	expectedStatusCode       int
	expectedResponse         interface{}
	expectedResponseLocation string
	expectedSetCookie        string
	expectedInternalErr      string

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

	if tc.expectedInternalErr != "" {
		require.EqualError(t, internalServerError, tc.expectedInternalErr,
			"Expected internal err: %s | Request Body: %s | Response Code: %d, Expected Response Body: %#v | Got Body: %s", tc.expectedInternalErr, tc.requestBody, rr.Code, tc.expectedResponse, rr.Body.String())
	} else {
		require.NoError(t, internalServerError)
		if tc.expectedStatusCode != 0 {
			require.Equal(t, tc.expectedStatusCode, rr.Code,
				"Expected body: %#v | Got body: %s", tc.expectedResponse, rr.Body.String())
		}
	}

	if tc.expectedResponseLocation != "" {
		redirectLocationURL, err := rr.Result().Location()
		require.NoError(t, err)
		require.Equal(t, tc.expectedResponseLocation, redirectLocationURL.String())
	}

	if tc.expectedSetCookie != "" {
		setCookie := rr.Header().Get("Set-Cookie")
		require.Equal(t, tc.expectedSetCookie, setCookie)
	}

	if tc.expectedResponse != nil {
		expectedBytes, err := json.Marshal(tc.expectedResponse)
		require.NoError(t, err)
		require.Equal(t, string(expectedBytes), strings.Trim(rr.Body.String(), "\n"))
	}
}

func behavesAsServiceCravingHandler(t *testing.T, method, url string, handler func(*env.AppEnv, http.ResponseWriter, *http.Request) error, serviceNames []string, baseCT ControllerTestCase) {
	t.Run("behaves as service craving handler", func(t *testing.T) {
		for _, sn := range serviceNames {
			baseEnv := *baseCT.env
			controllerTestCase := baseCT
			controllerTestCase.env = &baseEnv
			if sn == "AppService" {
				controllerTestCase.env.AppService = nil
				controllerTestCase.expectedInternalErr = "No App Service defined for handler"
			} else if sn == "AppVersionService" {
				controllerTestCase.env.AppVersionService = nil
				controllerTestCase.expectedInternalErr = "No App Version Service defined for handler"
			} else if sn == "ScreenshotService" {
				controllerTestCase.env.ScreenshotService = nil
				controllerTestCase.expectedInternalErr = "No Screenshot Service defined for handler"
			} else if sn == "FeatureGraphicService" {
				controllerTestCase.env.FeatureGraphicService = nil
				controllerTestCase.expectedInternalErr = "No Feature Graphic Service defined for handler"
			} else if sn == "AppSettingsService" {
				controllerTestCase.env.AppSettingsService = nil
				controllerTestCase.expectedInternalErr = "No App Settings Service defined for handler"
			} else if sn == "AppVersionEventService" {
				controllerTestCase.env.AppVersionEventService = nil
				controllerTestCase.expectedInternalErr = "No App Version Event Service defined for handler"
			} else if sn == "PublishTaskService" {
				controllerTestCase.env.PublishTaskService = nil
				controllerTestCase.expectedInternalErr = "No Publish Task Service defined for handler"
			} else if sn == "AppContactService" {
				controllerTestCase.env.AppContactService = nil
				controllerTestCase.expectedInternalErr = "No App Contact Service defined for handler"
			} else if sn == "RequestParams" {
				controllerTestCase.env.RequestParams = nil
				controllerTestCase.expectedInternalErr = "No RequestParams defined for handler"
			} else if sn == "AWS" {
				controllerTestCase.env.AWS = nil
				controllerTestCase.expectedInternalErr = "No AWS Provider defined for handler"
			} else if sn == "BitriseAPI" {
				controllerTestCase.env.BitriseAPI = nil
				controllerTestCase.expectedInternalErr = "No Bitrise API Service defined for handler"
			} else if sn == "WorkerService" {
				controllerTestCase.env.WorkerService = nil
				controllerTestCase.expectedInternalErr = "No Worker Service defined for handler"
			} else if sn == "Mailer" {
				controllerTestCase.env.Mailer = nil
				controllerTestCase.expectedInternalErr = "No Mailer defined for handler"
			} else {
				t.Fatalf("Invalid service element name defined: %s", sn)
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
			} else if ck == services.ContextKeyAuthorizedAppVersionID {
				controllerTestCase.contextElements[ck] = nil
				controllerTestCase.expectedInternalErr = "Authorized App Version ID not found in Context"
			} else if ck == services.ContextKeyAuthorizedAppContactID {
				controllerTestCase.contextElements[ck] = nil
				controllerTestCase.expectedInternalErr = "Authorized App Contact ID not found in Context"
			} else if ck == services.ContextKeyAuthorizedScreenshotID {
				controllerTestCase.contextElements[ck] = nil
				controllerTestCase.expectedInternalErr = "Authorized App Version Screenshot ID not found in Context"
			} else {

				t.Fatalf("Invalid context element name defined: %s", ck)
			}
			performControllerTest(t, method, url, handler, controllerTestCase)
		}
	})
}

// -----------
// Authentication
// -----------

type testAuthHandler struct{}

func (h *testAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpresponse.RespondWithSuccessNoErr(w, map[string]string{"message": "ok"})
}

type AuthenticationTestCase struct {
	requestHeaders    map[string]string
	requestFormValues map[string]string
	env               *env.AppEnv
	authHandlerFunc   func(*env.AppEnv, http.Handler) http.Handler

	expectedStatusCode int
	expectedBody       string
}

func performAuthenticationTest(t *testing.T,
	httpMethod, testURL string,
	tc AuthenticationTestCase,
) {
	t.Helper()

	req, err := http.NewRequest(httpMethod, testURL, nil)
	require.NoError(t, err)

	for key, value := range tc.requestHeaders {
		req.Header.Set(key, value)
	}
	if tc.requestFormValues != nil {
		for key, value := range tc.requestFormValues {
			if req.URL.RawQuery == "" {
				req.URL.RawQuery = fmt.Sprintf("%s=%s", key, value)
			} else {
				req.URL.RawQuery = fmt.Sprintf("%s&%s=%s", req.URL.RawQuery, key, value)
			}
		}
	}

	rr := httptest.NewRecorder()
	handler := tc.authHandlerFunc(tc.env, &testAuthHandler{})
	handler.ServeHTTP(rr, req)

	require.Equal(t, tc.expectedStatusCode, rr.Code)
	require.Equal(t, tc.expectedBody, rr.Body.String())
}

// -----------
// Authorization
// -----------

// AuthorizationTestCase ...
type AuthorizationTestCase struct {
	requestHeaders     map[string]string
	requestPayload     interface{}
	expectedStatusCode int
	expectedResponse   interface{}

	contextElements map[ctxpkg.RequestContextKey]interface{}
}

func performAuthorizationTest(t *testing.T,
	httpMethod, url string,
	handler http.Handler,
	tc AuthorizationTestCase,
) {
	t.Helper()

	requestPayloadBytes, err := json.Marshal(tc.requestPayload)
	require.NoError(t, err)

	r, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(requestPayloadBytes))
	require.NoError(t, err)

	for headerKey, headerValue := range tc.requestHeaders {
		r.Header.Set(headerKey, headerValue)
	}

	for key, val := range tc.contextElements {
		r = r.WithContext(context.WithValue(r.Context(), key, val))
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, r)

	if tc.expectedResponse != nil {
		expectedBytes, err := json.Marshal(tc.expectedResponse)
		require.NoError(t, err)
		require.Equal(t, string(expectedBytes), strings.Trim(rr.Body.String(), "\n"))
	}

	require.Equal(t, tc.expectedStatusCode, rr.Code)
}

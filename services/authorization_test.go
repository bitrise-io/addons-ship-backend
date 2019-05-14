package services_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/bitrise-io/api-utils/handlers"
	"github.com/bitrise-io/api-utils/providers"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// AuthorizationTestCase ...
type AuthorizationTestCase struct {
	requestHeaders     map[string]string
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

	r, err := http.NewRequest(httpMethod, url, nil)
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

func Test_AuthorizeForAppAccessHandlerFunc(t *testing.T) {
	authHandler := &handlers.TestAuthHandler{
		ContextElementList: map[string]ctxpkg.RequestContextKey{
			"authorizedAppID": services.ContextKeyAuthorizedAppID,
		},
	}
	httpMethod := "GET"
	url := "/apps/test_app_slug"

	t.Run("ok", func(t *testing.T) {
		handler := services.AuthorizeForAppAccessHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsProviderMock{
				Params: map[string]string{
					"app-slug": "test_app_slug",
				},
			},
			AppService: &testAppService{
				findFn: func(app *models.App) (*models.App, error) {
					require.Equal(t, app.AppSlug, "test_app_slug")
					require.Equal(t, app.APIToken, "test-auth-token")
					return &models.App{
						Record:   models.Record{ID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf")},
						AppSlug:  "test_app_slug",
						APIToken: "test-auth-token",
					}, nil
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			requestHeaders: map[string]string{
				"Authorization": "token test-auth-token",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"authorizedAppID": "211afc15-127a-40f9-8cbe-1dadc1f86cdf",
			},
		})
	})

	t.Run("when no auth token provided in header", func(t *testing.T) {
		handler := services.AuthorizeForAppAccessHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsProviderMock{
				Params: map[string]string{
					"app-slug": "test_app_slug",
				},
			},
			AppService: &testAppService{
				findFn: func(app *models.App) (*models.App, error) {
					require.Equal(t, app.AppSlug, "test_app_slug")
					require.Equal(t, app.APIToken, "test-auth-token")
					return &models.App{
						Record:   models.Record{ID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf")},
						AppSlug:  "test_app_slug",
						APIToken: "test-auth-token",
					}, nil
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: map[string]interface{}{
				"message": "Unauthorized",
			},
		})
	})

	t.Run("when no Request Params object is provided", func(t *testing.T) {
		handler := services.AuthorizeForAppAccessHandlerFunc(&env.AppEnv{
			AppService: &testAppService{
				findFn: func(app *models.App) (*models.App, error) {
					require.Equal(t, app.AppSlug, "test_app_slug")
					require.Equal(t, app.APIToken, "test-auth-token")
					return &models.App{
						Record:   models.Record{ID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf")},
						AppSlug:  "test_app_slug",
						APIToken: "test-auth-token",
					}, nil
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			requestHeaders: map[string]string{
				"Authorization": "token test-auth-token",
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"message": "Internal Server Error",
			},
		})
	})

	t.Run("when no app slug found in url params", func(t *testing.T) {
		handler := services.AuthorizeForAppAccessHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsProviderMock{
				Params: map[string]string{},
			},
			AppService: &testAppService{
				findFn: func(app *models.App) (*models.App, error) {
					require.Equal(t, app.AppSlug, "test_app_slug")
					require.Equal(t, app.APIToken, "test-auth-token")
					return &models.App{
						Record:   models.Record{ID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf")},
						AppSlug:  "test_app_slug",
						APIToken: "test-auth-token",
					}, nil
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			requestHeaders: map[string]string{
				"Authorization": "token test-auth-token",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"message": "App Slug not provided",
			},
		})
	})

	t.Run("when app no found in database", func(t *testing.T) {
		handler := services.AuthorizeForAppAccessHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsProviderMock{
				Params: map[string]string{
					"app-slug": "test_app_slug",
				},
			},
			AppService: &testAppService{
				findFn: func(app *models.App) (*models.App, error) {
					require.Equal(t, app.AppSlug, "test_app_slug")
					require.Equal(t, app.APIToken, "test-auth-token")
					return &models.App{}, gorm.ErrRecordNotFound
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			requestHeaders: map[string]string{
				"Authorization": "token test-auth-token",
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: map[string]interface{}{
				"message": "Not Found",
			},
		})
	})
	t.Run("when unexpected error happens at database query", func(t *testing.T) {
		handler := services.AuthorizeForAppAccessHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsProviderMock{
				Params: map[string]string{
					"app-slug": "test_app_slug",
				},
			},
			AppService: &testAppService{
				findFn: func(app *models.App) (*models.App, error) {
					require.Equal(t, app.AppSlug, "test_app_slug")
					require.Equal(t, app.APIToken, "test-auth-token")
					return &models.App{}, errors.New("SOME-SQL-ERROR")
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			requestHeaders: map[string]string{
				"Authorization": "token test-auth-token",
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"message": "Internal Server Error",
			},
		})
	})
}

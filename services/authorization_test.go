package services_test

import (
	"net/http"
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

func Test_AuthorizeForAppDeprovisioningHandlerFunc(t *testing.T) {
	authHandler := &handlers.TestAuthHandler{
		ContextElementList: map[string]ctxpkg.RequestContextKey{
			"authorizedAppID": services.ContextKeyAuthorizedAppID,
		},
	}
	httpMethod := "DELETE"
	url := "/provision/test_app_slug"

	t.Run("ok", func(t *testing.T) {
		handler := services.AuthorizeForAppDeprovisioningHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
				Params: map[string]string{
					"app-slug": "test_app_slug",
				},
			},
			AppService: &testAppService{
				findFn: func(app *models.App) (*models.App, error) {
					require.Equal(t, app.AppSlug, "test_app_slug")
					return &models.App{
						Record:  models.Record{ID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf")},
						AppSlug: "test_app_slug",
					}, nil
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			requestHeaders: map[string]string{
				"Authentication": "test-auth-token",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"authorizedAppID": "211afc15-127a-40f9-8cbe-1dadc1f86cdf",
			},
		})
	})

	t.Run("when no Request Params object is provided", func(t *testing.T) {
		handler := services.AuthorizeForAppDeprovisioningHandlerFunc(&env.AppEnv{
			AppService: &testAppService{
				findFn: func(app *models.App) (*models.App, error) {
					require.Equal(t, app.AppSlug, "test_app_slug")
					return &models.App{
						Record:  models.Record{ID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf")},
						AppSlug: "test_app_slug",
					}, nil
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			requestHeaders: map[string]string{
				"Authentication": "test-auth-token",
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"message": "Internal Server Error",
			},
		})
	})

	t.Run("when no app slug found in url params", func(t *testing.T) {
		handler := services.AuthorizeForAppDeprovisioningHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
				Params: map[string]string{},
			},
			AppService: &testAppService{
				findFn: func(app *models.App) (*models.App, error) {
					return &models.App{
						Record:  models.Record{ID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf")},
						AppSlug: "test_app_slug",
					}, nil
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			requestHeaders: map[string]string{
				"Authentication": "test-auth-token",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"message": "App Slug not provided",
			},
		})
	})

	t.Run("when no app service provided in app env", func(t *testing.T) {
		handler := services.AuthorizeForAppDeprovisioningHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
				Params: map[string]string{
					"app-slug": "test_app_slug",
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			requestHeaders: map[string]string{
				"Authentication": "test-auth-token",
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"message": "Internal Server Error",
			},
		})
	})

	t.Run("when app not found in database", func(t *testing.T) {
		handler := services.AuthorizeForAppDeprovisioningHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
				Params: map[string]string{
					"app-slug": "test_app_slug",
				},
			},
			AppService: &testAppService{
				findFn: func(app *models.App) (*models.App, error) {
					require.Equal(t, app.AppSlug, "test_app_slug")
					return &models.App{}, gorm.ErrRecordNotFound
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			requestHeaders: map[string]string{
				"Authentication": "test-auth-token",
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: map[string]interface{}{
				"message": "Not Found",
			},
		})
	})

	t.Run("when unexpected error happens at database query", func(t *testing.T) {
		handler := services.AuthorizeForAppDeprovisioningHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
				Params: map[string]string{
					"app-slug": "test_app_slug",
				},
			},
			AppService: &testAppService{
				findFn: func(app *models.App) (*models.App, error) {
					require.Equal(t, app.AppSlug, "test_app_slug")
					return &models.App{}, errors.New("SOME-SQL-ERROR")
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			requestHeaders: map[string]string{
				"Authentication": "test-auth-token",
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"message": "Internal Server Error",
			},
		})
	})
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
			RequestParams: &providers.RequestParamsMock{
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
			RequestParams: &providers.RequestParamsMock{
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
			RequestParams: &providers.RequestParamsMock{
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

	t.Run("when no app service provided in app env", func(t *testing.T) {
		handler := services.AuthorizeForAppAccessHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
				Params: map[string]string{
					"app-slug": "test_app_slug",
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

	t.Run("when app no found in database", func(t *testing.T) {
		handler := services.AuthorizeForAppAccessHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
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
			RequestParams: &providers.RequestParamsMock{
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

func Test_AuthorizeForAppVersionAccessHandlerFunc(t *testing.T) {
	authHandler := &handlers.TestAuthHandler{
		ContextElementList: map[string]ctxpkg.RequestContextKey{
			"authorizedAppID":        services.ContextKeyAuthorizedAppID,
			"authorizedAppVersionID": services.ContextKeyAuthorizedAppVersionID,
		},
	}
	httpMethod := "GET"
	url := "/apps/test_app_slug/versions/version_uuid"

	t.Run("ok", func(t *testing.T) {
		handler := services.AuthorizeForAppVersionAccessHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
				Params: map[string]string{
					"version-id": "de438ddc-98e5-4226-a5f4-fd2d53474879",
				},
			},
			AppVersionService: &testAppVersionService{
				findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
					require.Equal(t, appVersion.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
					return &models.AppVersion{
						Record: models.Record{ID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879")},
					}, nil
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID:        uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			requestHeaders: map[string]string{
				"Authorization": "token test-auth-token",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"authorizedAppID":        "211afc15-127a-40f9-8cbe-1dadc1f86cdf",
				"authorizedAppVersionID": "de438ddc-98e5-4226-a5f4-fd2d53474879",
			},
		})
	})

	t.Run("when no Request Params object is provided", func(t *testing.T) {
		handler := services.AuthorizeForAppVersionAccessHandlerFunc(&env.AppEnv{
			AppVersionService: &testAppVersionService{
				findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
					require.Equal(t, appVersion.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
					return &models.AppVersion{
						Record: models.Record{ID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879")},
					}, nil
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
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

	t.Run("when no app version id found in url params", func(t *testing.T) {
		handler := services.AuthorizeForAppVersionAccessHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
				Params: map[string]string{},
			},
			AppVersionService: &testAppVersionService{
				findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
					require.Equal(t, appVersion.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
					return &models.AppVersion{
						Record: models.Record{ID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879")},
					}, nil
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			requestHeaders: map[string]string{
				"Authorization": "token test-auth-token",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"message": "Failed to fetch URL param version-id",
			},
		})
	})

	t.Run("when no valid app version id found in url params", func(t *testing.T) {
		handler := services.AuthorizeForAppVersionAccessHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
				Params: map[string]string{
					"version-id": "invalid-uuid",
				},
			},
			AppVersionService: &testAppVersionService{
				findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
					require.Equal(t, appVersion.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
					return &models.AppVersion{
						Record: models.Record{ID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879")},
					}, nil
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			requestHeaders: map[string]string{
				"Authorization": "token test-auth-token",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"message": "Invalid UUID format for version-id",
			},
		})
	})

	t.Run("when no app version service is provided in app env", func(t *testing.T) {
		handler := services.AuthorizeForAppVersionAccessHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
				Params: map[string]string{
					"version-id": "de438ddc-98e5-4226-a5f4-fd2d53474879",
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
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

	t.Run("when app no found in database", func(t *testing.T) {
		handler := services.AuthorizeForAppVersionAccessHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
				Params: map[string]string{
					"version-id": "de438ddc-98e5-4226-a5f4-fd2d53474879",
				},
			},
			AppVersionService: &testAppVersionService{
				findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
					require.Equal(t, appVersion.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
					return &models.AppVersion{}, gorm.ErrRecordNotFound
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
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
		handler := services.AuthorizeForAppVersionAccessHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
				Params: map[string]string{
					"version-id": "de438ddc-98e5-4226-a5f4-fd2d53474879",
				},
			},
			AppVersionService: &testAppVersionService{
				findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
					require.Equal(t, appVersion.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
					return &models.AppVersion{}, errors.New("SOME-SQL-ERROR")
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
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

func Test_AuthorizeForAppVersionScreenshotAccessHandlerFunc(t *testing.T) {
	authHandler := &handlers.TestAuthHandler{
		ContextElementList: map[string]ctxpkg.RequestContextKey{
			"authorizedAppID":        services.ContextKeyAuthorizedAppID,
			"authorizedAppVersionID": services.ContextKeyAuthorizedAppVersionID,
			"authorizedScreenshotID": services.ContextKeyAuthorizedScreenshotID,
		},
	}
	httpMethod := "GET"
	url := "/apps/test_app_slug/versions/version_uuid/screenshots/screenshot_uuid"

	testAppID := "211afc15-127a-40f9-8cbe-1dadc1f86cdf"
	testAppVersionID := "de438ddc-98e5-4226-a5f4-fd2d53474879"
	testScreenshotID := "123afc15-127a-40f9-8cbe-1dadc1f86cdf"
	validRequestParams := &providers.RequestParamsMock{
		Params: map[string]string{
			"version-id":    testAppVersionID,
			"screenshot-id": testScreenshotID,
		},
	}

	successfulTestScreenshotService := &testScreenshotService{
		findFn: func(screenshot *models.Screenshot) (*models.Screenshot, error) {
			require.Equal(t, screenshot.AppVersionID.String(), testAppVersionID)
			require.Equal(t, screenshot.ID.String(), testScreenshotID)

			return &models.Screenshot{
				Record: models.Record{ID: uuid.FromStringOrNil(testScreenshotID)},
			}, nil
		},
	}

	testRequestHeaders := map[string]string{
		"Authorization": "token test-auth-token",
	}

	t.Run("ok", func(t *testing.T) {
		handler := services.AuthorizeForAppVersionScreenshotAccessHandlerFunc(&env.AppEnv{
			RequestParams:     validRequestParams,
			ScreenshotService: successfulTestScreenshotService,
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID:        uuid.FromStringOrNil(testAppID),
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil(testAppVersionID),
				services.ContextKeyAuthorizedScreenshotID: uuid.FromStringOrNil(testScreenshotID),
			},
			requestHeaders:     testRequestHeaders,
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"authorizedAppID":        testAppID,
				"authorizedAppVersionID": testAppVersionID,
				"authorizedScreenshotID": testScreenshotID,
			},
		})
	})

	t.Run("when no Request Params object is provided", func(t *testing.T) {
		handler := services.AuthorizeForAppVersionScreenshotAccessHandlerFunc(&env.AppEnv{
			ScreenshotService: successfulTestScreenshotService,
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil(testAppVersionID),
			},
			requestHeaders:     testRequestHeaders,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"message": "Internal Server Error",
			},
		})
	})

	t.Run("when no screenshot id found in url params", func(t *testing.T) {
		handler := services.AuthorizeForAppVersionScreenshotAccessHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
				Params: map[string]string{},
			},
			ScreenshotService: successfulTestScreenshotService,
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil(testAppVersionID),
			},
			requestHeaders:     testRequestHeaders,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"message": "Failed to fetch URL param screenshot-id",
			},
		})
	})

	t.Run("when no valid app version id found in url params", func(t *testing.T) {
		handler := services.AuthorizeForAppVersionScreenshotAccessHandlerFunc(&env.AppEnv{
			RequestParams: &providers.RequestParamsMock{
				Params: map[string]string{
					"screenshot-id": "invalid-uuid",
				},
			},
			ScreenshotService: successfulTestScreenshotService,
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil(testAppVersionID),
			},
			requestHeaders:     testRequestHeaders,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"message": "Invalid UUID format for screenshot-id",
			},
		})
	})

	t.Run("when no screenshot service is provided in app env", func(t *testing.T) {
		handler := services.AuthorizeForAppVersionScreenshotAccessHandlerFunc(&env.AppEnv{
			RequestParams: validRequestParams,
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil(testAppVersionID),
			},
			requestHeaders:     testRequestHeaders,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"message": "Internal Server Error",
			},
		})
	})

	t.Run("when app not found in database", func(t *testing.T) {
		handler := services.AuthorizeForAppVersionScreenshotAccessHandlerFunc(&env.AppEnv{
			RequestParams: validRequestParams,
			ScreenshotService: &testScreenshotService{
				findFn: func(screenshot *models.Screenshot) (*models.Screenshot, error) {
					require.Equal(t, screenshot.ID.String(), testScreenshotID)
					return &models.Screenshot{}, gorm.ErrRecordNotFound
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil(testAppVersionID),
			},
			requestHeaders:     testRequestHeaders,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: map[string]interface{}{
				"message": "Not Found",
			},
		})
	})

	t.Run("when unexpected error happens at database query", func(t *testing.T) {
		handler := services.AuthorizeForAppVersionScreenshotAccessHandlerFunc(&env.AppEnv{
			RequestParams: validRequestParams,
			ScreenshotService: &testScreenshotService{
				findFn: func(screenshot *models.Screenshot) (*models.Screenshot, error) {
					require.Equal(t, screenshot.ID.String(), testScreenshotID)
					return &models.Screenshot{}, errors.New("SOME-SQL-ERROR")
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil(testAppVersionID),
			},
			requestHeaders:     testRequestHeaders,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"message": "Internal Server Error",
			},
		})
	})
}

func Test_AuthorizeForWebhookHandlerFunc(t *testing.T) {
	authHandler := &handlers.TestAuthHandler{
		ContextElementList: map[string]ctxpkg.RequestContextKey{
			"authorizedAppVersionID": services.ContextKeyAuthorizedAppVersionID,
		},
	}
	httpMethod := "POST"
	url := "/webhook"

	t.Run("ok", func(t *testing.T) {
		handler := services.AuthorizeForWebhookHandlerFunc(&env.AppEnv{
			PublishTaskService: &testPublishTaskService{
				findFn: func(publishTask *models.PublishTask) (*models.PublishTask, error) {
					require.Equal(t, uuid.FromStringOrNil("13a94c5d-4609-404e-ae69-c625e93b8b71"), publishTask.TaskID)
					publishTask.AppVersionID = uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879")
					return publishTask, nil
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			requestPayload:     map[string]string{"task_id": "13a94c5d-4609-404e-ae69-c625e93b8b71"},
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"authorizedAppVersionID": "de438ddc-98e5-4226-a5f4-fd2d53474879",
			},
		})
	})

	t.Run("when request payload is invalid", func(t *testing.T) {
		handler := services.AuthorizeForWebhookHandlerFunc(&env.AppEnv{
			PublishTaskService: &testPublishTaskService{
				findFn: func(publishTask *models.PublishTask) (*models.PublishTask, error) {
					require.Equal(t, uuid.FromStringOrNil("13a94c5d-4609-404e-ae69-c625e93b8b71"), publishTask.TaskID)
					publishTask.AppVersionID = uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879")
					return publishTask, nil
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			requestPayload:     "invalid-json",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"message": "Invalid request body, JSON decode failed",
			},
		})
	})

	t.Run("when no publish task service is defined", func(t *testing.T) {
		handler := services.AuthorizeForWebhookHandlerFunc(&env.AppEnv{
			PublishTaskService: nil,
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			requestPayload:     map[string]string{"task_id": "13a94c5d-4609-404e-ae69-c625e93b8b71"},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"message": "Internal Server Error",
			},
		})
	})

	t.Run("when no publish task found by task id", func(t *testing.T) {
		handler := services.AuthorizeForWebhookHandlerFunc(&env.AppEnv{
			PublishTaskService: &testPublishTaskService{
				findFn: func(publishTask *models.PublishTask) (*models.PublishTask, error) {
					require.Equal(t, uuid.FromStringOrNil("13a94c5d-4609-404e-ae69-c625e93b8b71"), publishTask.TaskID)
					return nil, gorm.ErrRecordNotFound
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			requestPayload:     map[string]string{"task_id": "13a94c5d-4609-404e-ae69-c625e93b8b71"},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: map[string]interface{}{
				"message": "Not Found",
			},
		})
	})

	t.Run("when error happens at finding publish task", func(t *testing.T) {
		handler := services.AuthorizeForWebhookHandlerFunc(&env.AppEnv{
			PublishTaskService: &testPublishTaskService{
				findFn: func(publishTask *models.PublishTask) (*models.PublishTask, error) {
					require.Equal(t, uuid.FromStringOrNil("13a94c5d-4609-404e-ae69-c625e93b8b71"), publishTask.TaskID)
					return nil, errors.New("SOME-SQL-ERROR")
				},
			},
		}, authHandler)
		performAuthorizationTest(t, httpMethod, url, handler, AuthorizationTestCase{
			requestPayload:     map[string]string{"task_id": "13a94c5d-4609-404e-ae69-c625e93b8b71"},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"message": "Internal Server Error",
			},
		})
	})
}

package services_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func Test_ResourcesHandler(t *testing.T) {
	httpMethod := "GET"
	testURL := "/resources/*"
	handler := services.ResourcesHandler

	behavesAsContextCravingHandler(t, httpMethod, testURL, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
	})

	behavesAsServiceCravingHandler(t, httpMethod, testURL, handler, []string{"AppService"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppService: &testAppService{},
		},
	})

	t.Run("ok", func(t *testing.T) {
		testURL := "/resources/apps/test-app-slug"
		testAPIURL, err := url.Parse("http://api.bitrise.io")
		require.NoError(t, err)
		performControllerTest(t, httpMethod, testURL, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				BitriseAPIRootURL: testAPIURL,
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return app, nil
					},
				},
			},
			expectedStatusCode:       http.StatusPermanentRedirect,
			expectedResponseLocation: "https://api.bitrise.io/v0.1/apps/test-app-slug",
		})
	})

	t.Run("when app cannot be found in database", func(t *testing.T) {
		testURL := "/resources/apps/test-app-slug"
		testAPIURL, err := url.Parse("http://api.bitrise.io")
		require.NoError(t, err)
		performControllerTest(t, httpMethod, testURL, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				BitriseAPIRootURL: testAPIURL,
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return nil, gorm.ErrRecordNotFound
					},
				},
			},
			expectedInternalErr: "SQL Error: record not found",
		})
	})
}

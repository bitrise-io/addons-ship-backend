package services_test

import (
	"net/http"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_AppGetHandler(t *testing.T) {
	httpMethod := "GET"
	url := "/apps/{app-slug}"
	handler := services.AppGetHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppService"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppService: &testAppService{},
		},
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppService: &testAppService{},
		},
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						require.Equal(t, app.ID.String(), "211afc15-127a-40f9-8cbe-1dadc1f86cdf")
						return nil, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   services.AppGetResponse{},
		})
	})

	t.Run("ok - more complex", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return &models.App{
							AppSlug: "test_app_slug",
						}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppGetResponse{
				Data: &models.App{
					AppSlug: "test_app_slug",
				},
			},
		})
	})

	t.Run("error - not found in database", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return nil, gorm.ErrRecordNotFound
					},
				},
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: httpresponse.StandardErrorRespModel{
				Message: "Not Found",
			},
		})
	})

	t.Run("error - unexpected error in database", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
			},
			expectedStatusCode:  http.StatusNotFound,
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})
}

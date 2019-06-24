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

func Test_DeprovisionHandler(t *testing.T) {
	httpMethod := "DELETE"
	url := "/provision/{app-slug}"
	handler := services.DeprovisionHandler

	testAppID := uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf")

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
			AppVersionService: &testAppVersionService{},
		},
	})

	t.Run("ok", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						require.Equal(t, testAppID, app.ID)
						return app, nil
					},
					deleteFn: func(app *models.App) error {
						require.Equal(t, testAppID, app.ID)
						return nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   models.App{Record: models.Record{ID: testAppID}},
		})
	})

	t.Run("when app not found in database", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return nil, gorm.ErrRecordNotFound
					},
					deleteFn: func(app *models.App) error {
						require.Equal(t, testAppID, app.ID)
						return nil
					},
				},
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Not Found"},
		})
	})

	t.Run("when database error happens at find", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
					deleteFn: func(app *models.App) error {
						require.Equal(t, testAppID, app.ID)
						return nil
					},
				},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when database error happens at find", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: testAppID,
			},
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						require.Equal(t, testAppID, app.ID)
						return app, nil
					},
					deleteFn: func(app *models.App) error {
						return errors.New("SOME-SQL-ERROR")
					},
				},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})
}

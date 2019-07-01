package services_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/api-utils/providers"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_AppEventsGetHandler(t *testing.T) {
	httpMethod := "GET"
	url := "/apps/{app-slug}/events"
	handler := services.AppEventsGetHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppEventService", "AWS"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppEventService: &testAppEventService{
				findAllFn: func(app *models.App) ([]models.AppEvent, error) {
					return []models.AppEvent{}, nil
				},
			},
			AWS: &providers.AWSMock{},
		},
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppEventService: &testAppEventService{},
			AWS:             &providers.AWSMock{},
		},
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppEventService: &testAppEventService{
					findAllFn: func(app *models.App) ([]models.AppEvent, error) {
						return []models.AppEvent{}, nil
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppEventsGetResponse{
				Data: []services.AppEventData{},
			},
		})
	})

	t.Run("ok - more complex", func(t *testing.T) {
		testAppEventUUID := uuid.FromStringOrNil("b22daf1a-7a4b-482d-a6c5-f55dbd229afc")

		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			env: &env.AppEnv{
				AppEventService: &testAppEventService{
					findAllFn: func(app *models.App) ([]models.AppEvent, error) {
						require.Equal(t, app.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
						app.AppSlug = "test-app-slug"
						return []models.AppEvent{
							models.AppEvent{
								Record: models.Record{ID: testAppEventUUID},
								App:    *app,
							},
						}, nil
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return fmt.Sprintf("http://presigned.aws.url/%s", path), nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppEventsGetResponse{
				Data: []services.AppEventData{
					services.AppEventData{
						AppEvent: models.AppEvent{
							Record: models.Record{ID: testAppEventUUID},
						},
						LogDownloadURL: "http://presigned.aws.url/logs/test-app-slug/b22daf1a-7a4b-482d-a6c5-f55dbd229afc.log",
					},
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
				AppEventService: &testAppEventService{
					findAllFn: func(app *models.App) ([]models.AppEvent, error) {
						return nil, gorm.ErrRecordNotFound
					},
				},
				AWS: &providers.AWSMock{},
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Not Found"},
		})
	})

	t.Run("error - unexpected error in database", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppEventService: &testAppEventService{
					findAllFn: func(app *models.App) ([]models.AppEvent, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
				AWS: &providers.AWSMock{},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("error - when generating AWS presigned URL", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			env: &env.AppEnv{
				AppEventService: &testAppEventService{
					findAllFn: func(app *models.App) ([]models.AppEvent, error) {
						return []models.AppEvent{models.AppEvent{}}, nil
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", errors.New("SOME-AWS-ERROR")
					},
				},
			},
			expectedInternalErr: "SOME-AWS-ERROR",
		})
	})
}

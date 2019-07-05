package services_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/redis"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/c2fo/testify/require"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_WebhookPostHandler(t *testing.T) {
	httpMethod := "POST"
	url := "/webhook"
	handler := services.WebhookPostHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppVersionService", "AppVersionEventService"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppVersionService:      &testAppVersionService{},
			AppVersionEventService: &testAppVersionEventService{},
		},
		requestBody: `{}`,
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppVersionID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppVersionService:      &testAppVersionService{},
			AppVersionEventService: &testAppVersionEventService{},
		},
		requestBody: `{}`,
	})

	t.Run("when incoming webhook has 'log' type", func(t *testing.T) {
		t.Run("ok - minimal", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
				},
				env: &env.AppEnv{
					AppVersionService: &testAppVersionService{
						findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
							return &models.AppVersion{}, nil
						},
					},
					AppVersionEventService: &testAppVersionEventService{},
					WorkerService: &testWorkerService{
						enqueueStoreLogChunkToRedisFn: func(string, models.LogChunk, time.Duration) error {
							return nil
						},
					},
				},
				requestBody:        `{"type_id":"log"}`,
				expectedStatusCode: http.StatusOK,
				expectedResponse:   httpresponse.StandardErrorRespModel{Message: "ok"},
			})
		})

		t.Run("ok - more complex", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
				},
				env: &env.AppEnv{
					AppVersionService: &testAppVersionService{
						findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
							return &models.AppVersion{}, nil
						},
					},
					AppVersionEventService: &testAppVersionEventService{},
					WorkerService: &testWorkerService{
						enqueueStoreLogChunkToRedisFn: func(taskID string, chunk models.LogChunk, secondsToStartFromNow time.Duration) error {
							require.Equal(t, "96e72f92-6e4c-40d5-b829-48a1ea6440a1", taskID)
							require.Equal(t, models.LogChunk{
								TaskID:  uuid.FromStringOrNil("96e72f92-6e4c-40d5-b829-48a1ea6440a1"),
								Content: "My awesome log chunk",
								Pos:     1,
							}, chunk)
							require.Equal(t, int64(5), int64(secondsToStartFromNow))
							return nil
						},
					},
				},
				requestBody:        `{"type_id":"log","task_id":"96e72f92-6e4c-40d5-b829-48a1ea6440a1","data":{"chunk":"My awesome log chunk","position":1}}`,
				expectedStatusCode: http.StatusOK,
				expectedResponse:   httpresponse.StandardErrorRespModel{Message: "ok"},
			})
		})

		t.Run("when chunk data has wrong format", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
				},
				env: &env.AppEnv{
					AppVersionService: &testAppVersionService{
						findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
							return &models.AppVersion{}, nil
						},
					},
					AppVersionEventService: &testAppVersionEventService{},
					WorkerService: &testWorkerService{
						enqueueStoreLogChunkToRedisFn: func(string, models.LogChunk, time.Duration) error {
							return nil
						},
					},
				},
				requestBody:        `{"type_id":"log","data":"invalid JSON"}`,
				expectedStatusCode: http.StatusBadRequest,
				expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Invalid format of log data"},
			})
		})

		t.Run("when error happens at worker enqueue", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
				},
				env: &env.AppEnv{
					AppVersionService: &testAppVersionService{
						findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
							return &models.AppVersion{}, nil
						},
					},
					AppVersionEventService: &testAppVersionEventService{},
					WorkerService: &testWorkerService{
						enqueueStoreLogChunkToRedisFn: func(string, models.LogChunk, time.Duration) error {
							return errors.New("SOME-WORKER-ERROR")
						},
					},
				},
				requestBody:         `{"type_id":"log"}`,
				expectedInternalErr: "Worker error: SOME-WORKER-ERROR",
			})
		})
	})

	t.Run("when incoming webhook has 'status' type", func(t *testing.T) {
		t.Run("when status is 'started'", func(t *testing.T) {
			t.Run("ok - minimal", func(t *testing.T) {
				performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
					contextElements: map[ctxpkg.RequestContextKey]interface{}{
						services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
					},
					env: &env.AppEnv{
						AppVersionService: &testAppVersionService{
							findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
								return &models.AppVersion{}, nil
							},
						},
						AppVersionEventService: &testAppVersionEventService{
							createFn: func(*models.AppVersionEvent) (*models.AppVersionEvent, error) {
								return nil, nil
							},
						},
						WorkerService: &testWorkerService{},
						Redis: &redis.Mock{
							SetFn: func(string, interface{}, int) error {
								return nil
							},
						},
					},
					requestBody:        `{"type_id":"status","data":{"new_status":"started"}}`,
					expectedStatusCode: http.StatusOK,
					expectedResponse:   httpresponse.StandardErrorRespModel{Message: "ok"},
				})
			})

			t.Run("ok - more complex", func(t *testing.T) {
				performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
					contextElements: map[ctxpkg.RequestContextKey]interface{}{
						services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
					},
					env: &env.AppEnv{
						AppVersionService: &testAppVersionService{
							findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
								return &models.AppVersion{}, nil
							},
						},
						AppVersionEventService: &testAppVersionEventService{
							createFn: func(*models.AppVersionEvent) (*models.AppVersionEvent, error) {
								return nil, nil
							},
						},
						WorkerService:       &testWorkerService{},
						RedisExpirationTime: 10,
						Redis: &redis.Mock{
							SetFn: func(key string, value interface{}, ttl int) error {
								require.Equal(t, "96e72f92-6e4c-40d5-b829-48a1ea6440a1_chunk_count", key)
								require.Equal(t, 0, value)
								require.Equal(t, 10, ttl)
								return nil
							},
						},
					},
					requestBody:        `{"type_id":"status","task_id":"96e72f92-6e4c-40d5-b829-48a1ea6440a1","data":{"new_status":"started"}}`,
					expectedStatusCode: http.StatusOK,
					expectedResponse:   httpresponse.StandardErrorRespModel{Message: "ok"},
				})
			})

			t.Run("when error happens at creating new app version event", func(t *testing.T) {
				performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
					contextElements: map[ctxpkg.RequestContextKey]interface{}{
						services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
					},
					env: &env.AppEnv{
						AppVersionService: &testAppVersionService{
							findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
								return &models.AppVersion{}, nil
							},
						},
						AppVersionEventService: &testAppVersionEventService{
							createFn: func(*models.AppVersionEvent) (*models.AppVersionEvent, error) {
								return nil, errors.New("SOME-SQL-ERROR")
							},
						},
						WorkerService: &testWorkerService{},
						Redis: &redis.Mock{
							SetFn: func(string, interface{}, int) error {
								return nil
							},
						},
					},
					requestBody:         `{"type_id":"status","data":{"new_status":"started"}}`,
					expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
				})
			})

			t.Run("when error happens at saving chunk count to Redis", func(t *testing.T) {
				performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
					contextElements: map[ctxpkg.RequestContextKey]interface{}{
						services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
					},
					env: &env.AppEnv{
						AppVersionService: &testAppVersionService{
							findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
								return &models.AppVersion{}, nil
							},
						},
						AppVersionEventService: &testAppVersionEventService{
							createFn: func(*models.AppVersionEvent) (*models.AppVersionEvent, error) {
								return nil, nil
							},
						},
						WorkerService: &testWorkerService{},
						Redis: &redis.Mock{
							SetFn: func(string, interface{}, int) error {
								return errors.New("SOME-REDIS-ERROR")
							},
						},
					},
					requestBody:         `{"type_id":"status","data":{"new_status":"started"}}`,
					expectedInternalErr: "SOME-REDIS-ERROR",
				})
			})
		})
	})
}

package services_test

import (
	"net/http"
	"testing"

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
	url := "/task-webhook"
	handler := services.WebhookPostHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppVersionService", "AppVersionEventService", "WorkerService"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppVersionService:      &testAppVersionService{},
			AppVersionEventService: &testAppVersionEventService{},
			WorkerService:          &testWorkerService{},
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
						enqueueStoreLogChunkToRedisFn: func(string, models.LogChunk, int64) error {
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
						enqueueStoreLogChunkToRedisFn: func(taskID string, chunk models.LogChunk, secondsToStartFromNow int64) error {
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
						enqueueStoreLogChunkToRedisFn: func(string, models.LogChunk, int64) error {
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
						enqueueStoreLogChunkToRedisFn: func(string, models.LogChunk, int64) error {
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
		testAppVersionID := uuid.FromStringOrNil("e2915475-381d-4252-b5ec-c0fe511b12e8")

		t.Run("when status data has invalid format", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
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
				requestBody:        `{"type_id":"status","data":"some invalid JSON"}`,
				expectedStatusCode: http.StatusBadRequest,
				expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Invalid format of status data"},
			})
		})

		t.Run("when status in payload is not valid", func(t *testing.T) {
			performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
				contextElements: map[ctxpkg.RequestContextKey]interface{}{
					services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
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
				requestBody:         `{"type_id":"status","data":{"new_status":"some invalid status"}}`,
				expectedInternalErr: "Invalid status of incoming webhook: some invalid status",
			})
		})

		t.Run("when status is 'started'", func(t *testing.T) {
			t.Run("ok - minimal", func(t *testing.T) {
				performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
					contextElements: map[ctxpkg.RequestContextKey]interface{}{
						services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
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
						services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
					},
					env: &env.AppEnv{
						AppVersionService: &testAppVersionService{
							findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
								return &models.AppVersion{Record: models.Record{ID: testAppVersionID}}, nil
							},
						},
						AppVersionEventService: &testAppVersionEventService{
							createFn: func(event *models.AppVersionEvent) (*models.AppVersionEvent, error) {
								require.Equal(t, &models.AppVersionEvent{
									Status:       "in_progress",
									Text:         "Publishing has started",
									AppVersionID: testAppVersionID,
								}, event)
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
						services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
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
		t.Run("when status is 'finished'", func(t *testing.T) {
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
							createFn: func(event *models.AppVersionEvent) (*models.AppVersionEvent, error) {
								event.AppVersion = models.AppVersion{Record: models.Record{ID: uuid.NewV4()}, App: models.App{AppSlug: "test-app-slug"}}
								return event, nil
							},
						},
						WorkerService: &testWorkerService{
							enqueueStoreLogToAWSFn: func(uuid.UUID, int64, string, int64) error {
								return nil
							},
						},
					},
					requestBody:        `{"type_id":"status","data":{"new_status":"finished"}}`,
					expectedStatusCode: http.StatusOK,
					expectedResponse:   httpresponse.StandardErrorRespModel{Message: "ok"},
				})
			})

			t.Run("ok - more complex - success", func(t *testing.T) {
				performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
					contextElements: map[ctxpkg.RequestContextKey]interface{}{
						services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
					},
					env: &env.AppEnv{
						AppVersionService: &testAppVersionService{
							findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
								return &models.AppVersion{Record: models.Record{ID: testAppVersionID}}, nil
							},
						},
						AppVersionEventService: &testAppVersionEventService{
							createFn: func(event *models.AppVersionEvent) (*models.AppVersionEvent, error) {
								require.Equal(t, &models.AppVersionEvent{
									Status:       "success",
									Text:         "Successfully published",
									AppVersionID: testAppVersionID,
								}, event)
								event.ID = uuid.FromStringOrNil("507db32c-9f92-43b6-9a53-d8d7594736c7")
								event.AppVersion = models.AppVersion{Record: models.Record{ID: testAppVersionID}, App: models.App{AppSlug: "test-app-slug"}}
								return event, nil
							},
						},
						WorkerService: &testWorkerService{
							enqueueStoreLogToAWSFn: func(taskID uuid.UUID, logChunkCount int64, awsPath string, secondsToStartFromNow int64) error {
								require.Equal(t, "96e72f92-6e4c-40d5-b829-48a1ea6440a1", taskID.String())
								require.Equal(t, 2, logChunkCount)
								require.Equal(t, "logs/test-app-slug/e2915475-381d-4252-b5ec-c0fe511b12e8/507db32c-9f92-43b6-9a53-d8d7594736c7.log", awsPath)
								require.Equal(t, 30, secondsToStartFromNow)
								return nil
							},
						},
					},
					requestBody:        `{"type_id":"status","task_id":"96e72f92-6e4c-40d5-b829-48a1ea6440a1","data":{"new_status":"finished","exit_code":0,"generated_log_chunk_count":2}}`,
					expectedStatusCode: http.StatusOK,
					expectedResponse:   httpresponse.StandardErrorRespModel{Message: "ok"},
				})
			})

			t.Run("ok - more complex - failed", func(t *testing.T) {
				performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
					contextElements: map[ctxpkg.RequestContextKey]interface{}{
						services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
					},
					env: &env.AppEnv{
						AppVersionService: &testAppVersionService{
							findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
								return &models.AppVersion{Record: models.Record{ID: testAppVersionID}}, nil
							},
						},
						AppVersionEventService: &testAppVersionEventService{
							createFn: func(event *models.AppVersionEvent) (*models.AppVersionEvent, error) {
								require.Equal(t, &models.AppVersionEvent{
									Status:       "failed",
									Text:         "Failed to publish",
									AppVersionID: testAppVersionID,
								}, event)
								event.ID = uuid.FromStringOrNil("507db32c-9f92-43b6-9a53-d8d7594736c7")
								event.AppVersion = models.AppVersion{Record: models.Record{ID: testAppVersionID}, App: models.App{AppSlug: "test-app-slug"}}
								return event, nil
							},
						},
						WorkerService: &testWorkerService{
							enqueueStoreLogToAWSFn: func(taskID uuid.UUID, logChunkCount int64, awsPath string, secondsToStartFromNow int64) error {
								require.Equal(t, "96e72f92-6e4c-40d5-b829-48a1ea6440a1", taskID.String())
								require.Equal(t, 2, logChunkCount)
								require.Equal(t, "logs/test-app-slug/e2915475-381d-4252-b5ec-c0fe511b12e8/507db32c-9f92-43b6-9a53-d8d7594736c7.log", awsPath)
								require.Equal(t, 30, secondsToStartFromNow)
								return nil
							},
						},
					},
					requestBody:        `{"type_id":"status","task_id":"96e72f92-6e4c-40d5-b829-48a1ea6440a1","data":{"new_status":"finished","exit_code":-1,"generated_log_chunk_count":2}}`,
					expectedStatusCode: http.StatusOK,
					expectedResponse:   httpresponse.StandardErrorRespModel{Message: "ok"},
				})
			})

			t.Run("when error happens at creating new app version event", func(t *testing.T) {
				performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
					contextElements: map[ctxpkg.RequestContextKey]interface{}{
						services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
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
					},
					requestBody:         `{"type_id":"status","data":{"new_status":"finished"}}`,
					expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
				})
			})

			t.Run("when AWS path cannot be constructed", func(t *testing.T) {
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
							createFn: func(event *models.AppVersionEvent) (*models.AppVersionEvent, error) {
								event.AppVersion = models.AppVersion{Record: models.Record{ID: uuid.NewV4()}}
								return event, nil
							},
						},
						WorkerService: &testWorkerService{
							enqueueStoreLogToAWSFn: func(uuid.UUID, int64, string, int64) error {
								return nil
							},
						},
					},
					requestBody:         `{"type_id":"status","data":{"new_status":"finished"}}`,
					expectedInternalErr: "App has empty App Slug, App has to be preloaded",
				})
			})

			t.Run("when error happens at enqueuing new job", func(t *testing.T) {
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
							createFn: func(event *models.AppVersionEvent) (*models.AppVersionEvent, error) {
								event.AppVersion = models.AppVersion{Record: models.Record{ID: uuid.NewV4()}, App: models.App{AppSlug: "test-app-slug"}}
								return event, nil
							},
						},
						WorkerService: &testWorkerService{
							enqueueStoreLogToAWSFn: func(uuid.UUID, int64, string, int64) error {
								return errors.New("SOME-WORKER-ERROR")
							},
						},
					},
					requestBody:         `{"type_id":"status","data":{"new_status":"finished"}}`,
					expectedInternalErr: "Worker error: SOME-WORKER-ERROR",
				})
			})
		})
	})

	t.Run("when request body is invalid", func(t *testing.T) {
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
					enqueueStoreLogChunkToRedisFn: func(string, models.LogChunk, int64) error {
						return nil
					},
				},
			},
			requestBody:        `invalid JSON`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Invalid request body, JSON decode failed"},
		})
	})

	t.Run("when error happens at finding app version by authorized ID", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
				AppVersionEventService: &testAppVersionEventService{},
				WorkerService: &testWorkerService{
					enqueueStoreLogChunkToRedisFn: func(string, models.LogChunk, int64) error {
						return nil
					},
				},
			},
			requestBody:         `{"type_id":"log"}`,
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when webhook type is invalid", func(t *testing.T) {
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
					enqueueStoreLogChunkToRedisFn: func(string, models.LogChunk, int64) error {
						return nil
					},
				},
			},
			requestBody:         `{"type_id":"invalid hook type"}`,
			expectedInternalErr: "Invalid type of webhook: invalid hook type",
		})
	})
}

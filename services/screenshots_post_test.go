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
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_ScreenshotsPostHandler(t *testing.T) {
	httpMethod := "POST"
	url := "/apps/{app-slug}/versions/{version-id}/screenshots"
	handler := services.ScreenshotsPostHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"ScreenshotService"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			ScreenshotService: &testScreenshotService{},
			AWS:               &providers.AWSMock{},
		},
		requestBody: `{}`,
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppVersionID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			ScreenshotService: &testScreenshotService{},
			AWS:               &providers.AWSMock{},
		},
		requestBody: `{}`,
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				ScreenshotService: &testScreenshotService{
					batchCreateFn: func(screenshots []*models.Screenshot) ([]*models.Screenshot, []error, error) {
						return nil, nil, nil
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedPUTURLFn: func(path string, expiration time.Duration, size int64) (string, error) {
						return "", nil
					},
				},
			},
			requestBody:        `{}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.ScreenshotsGetResponse{
				Data: []services.ScreenshotData{},
			},
		})
	})

	t.Run("ok - more complex", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			env: &env.AppEnv{
				ScreenshotService: &testScreenshotService{
					batchCreateFn: func(screenshots []*models.Screenshot) ([]*models.Screenshot, []error, error) {
						appVersionID := uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879")

						require.Equal(t, &models.Screenshot{
							AppVersionID: appVersionID,
							Uploadable: models.Uploadable{
								Filename: "screenshot.png",
								Filesize: 1234,
							},
							DeviceType: "iPhone XS Max",
							ScreenSize: "6.5 inch",
						}, screenshots[0])
						require.Equal(t, &models.Screenshot{
							AppVersionID: appVersionID,
							Uploadable: models.Uploadable{
								Filename: "screenshot2.png",
								Filesize: 4321,
							},
							DeviceType: "iPhone XS",
							ScreenSize: "5.5 inch",
						}, screenshots[1])

						return []*models.Screenshot{
							&models.Screenshot{
								Uploadable: models.Uploadable{Filename: "screenshot.png"},
								DeviceType: "iPhone XS Max",
								ScreenSize: "6.5 inch",
								AppVersion: models.AppVersion{
									Record: models.Record{ID: appVersionID},
									App:    models.App{AppSlug: "test-app-slug"},
								},
							},
							&models.Screenshot{
								Uploadable: models.Uploadable{Filename: "screenshot2.png"},
								DeviceType: "iPhone XS",
								ScreenSize: "5.5 inch",
								AppVersion: models.AppVersion{
									Record: models.Record{ID: appVersionID},
									App:    models.App{AppSlug: "test-app-slug"},
								},
							},
						}, nil, nil
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedPUTURLFn: func(path string, expiration time.Duration, size int64) (string, error) {
						return fmt.Sprintf("http://presigned.aws.url/%s", path), nil
					},
				},
			},
			requestBody: `{"screenshots":[` +
				`{"filename":"screenshot.png","filesize":1234,"device_type":"iPhone XS Max","screen_size":"6.5 inch"},` +
				`{"filename":"screenshot2.png","filesize":4321,"device_type":"iPhone XS","screen_size":"5.5 inch"}` +
				`]}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.ScreenshotsGetResponse{
				Data: []services.ScreenshotData{
					services.ScreenshotData{
						Screenshot: models.Screenshot{
							Uploadable: models.Uploadable{Filename: "screenshot.png"},
							DeviceType: "iPhone XS Max",
							ScreenSize: "6.5 inch",
						},
						UploadURL: "http://presigned.aws.url/test-app-slug/de438ddc-98e5-4226-a5f4-fd2d53474879/iPhone XS Max (6.5 inch)/screenshot.png",
					},
					services.ScreenshotData{
						Screenshot: models.Screenshot{
							Uploadable: models.Uploadable{Filename: "screenshot2.png"},
							DeviceType: "iPhone XS",
							ScreenSize: "5.5 inch",
						},
						UploadURL: "http://presigned.aws.url/test-app-slug/de438ddc-98e5-4226-a5f4-fd2d53474879/iPhone XS (5.5 inch)/screenshot2.png",
					},
				},
			},
		})
	})

	t.Run("error - invalid request body format", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				ScreenshotService: &testScreenshotService{
					batchCreateFn: func(screenshots []*models.Screenshot) ([]*models.Screenshot, []error, error) {
						return nil, nil, nil
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedPUTURLFn: func(path string, expiration time.Duration, size int64) (string, error) {
						return "", nil
					},
				},
			},
			requestBody:        `invalid request body`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   map[string]string{"message": "Invalid request body, JSON decode failed"},
		})
	})

	t.Run("error - validation error at screenshot create", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				ScreenshotService: &testScreenshotService{
					batchCreateFn: func(screenshots []*models.Screenshot) ([]*models.Screenshot, []error, error) {
						return nil, []error{errors.New("SOME-VALIDATION-ERROR")}, nil
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedPUTURLFn: func(path string, expiration time.Duration, size int64) (string, error) {
						return "", nil
					},
				},
			},
			requestBody: `{"screenshots":[` +
				`{"filename":"screenshot.png","filesize":1234,"device_type":"iPhone XS Max","screen_size":"6.5 inch"},` +
				`{"filename":"screenshot2.png","filesize":4321,"device_type":"iPhone XS","screen_size":"5.5 inch"}` +
				`]}`,
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse: httpresponse.ValidationErrorRespModel{
				Message: "Unprocessable Entity",
				Errors:  []string{"SOME-VALIDATION-ERROR"},
			},
		})
	})

	t.Run("error - unexpected error in database", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				ScreenshotService: &testScreenshotService{
					batchCreateFn: func(screenshots []*models.Screenshot) ([]*models.Screenshot, []error, error) {
						return nil, nil, errors.New("SOME-SQL-ERROR")
					},
				},
				AWS: &providers.AWSMock{},
			},
			requestBody: `{"screenshots":[` +
				`{"filename":"screenshot.png","filesize":1234,"device_type":"iPhone XS Max","screen_size":"6.5 inch"},` +
				`{"filename":"screenshot2.png","filesize":4321,"device_type":"iPhone XS","screen_size":"5.5 inch"}` +
				`]}`,
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("error - when generating AWS presigned URL", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			env: &env.AppEnv{
				ScreenshotService: &testScreenshotService{
					batchCreateFn: func(screenshots []*models.Screenshot) ([]*models.Screenshot, []error, error) {
						appVersionID := uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879")
						return []*models.Screenshot{
							&models.Screenshot{
								Uploadable: models.Uploadable{Filename: "screenshot.png"},
								DeviceType: "iPhone XS Max",
								ScreenSize: "6.5 inch",
								AppVersion: models.AppVersion{
									Record: models.Record{ID: appVersionID},
									App:    models.App{AppSlug: "test-app-slug"},
								},
							},
							&models.Screenshot{
								Uploadable: models.Uploadable{Filename: "screenshot2.png"},
								DeviceType: "iPhone XS",
								ScreenSize: "5.5 inch",
								AppVersion: models.AppVersion{
									Record: models.Record{ID: appVersionID},
									App:    models.App{AppSlug: "test-app-slug"},
								},
							},
						}, nil, nil
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedPUTURLFn: func(path string, expiration time.Duration, size int64) (string, error) {
						return "", errors.New("SOME-AWS-ERROR")
					},
				},
			},
			requestBody: `{"screenshots":[` +
				`{"filename":"screenshot.png","filesize":1234,"device_type":"iPhone XS Max","screen_size":"6.5 inch"},` +
				`{"filename":"screenshot2.png","filesize":4321,"device_type":"iPhone XS","screen_size":"5.5 inch"}` +
				`]}`,
			expectedInternalErr: "SOME-AWS-ERROR",
		})
	})
}

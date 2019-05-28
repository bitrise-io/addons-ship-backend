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

func Test_ScreenshotsUploadedPatchHandler(t *testing.T) {
	httpMethod := "PATCH"
	url := "/apps/{app-slug}/versions/{version-id}/screenshots/uploaded"
	handler := services.ScreenshotsUploadedPatchHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"ScreenshotService", "AWS"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			ScreenshotService: &testScreenshotService{
				findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
					return []models.Screenshot{}, nil
				},
				batchUpdateFn: func(screenshots []models.Screenshot, whitelist []string) ([]error, error) {
					return nil, nil
				},
			},
			AWS: &providers.AWSMock{},
		},
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppVersionID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			ScreenshotService: &testScreenshotService{},
			AWS:               &providers.AWSMock{},
		},
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
					batchUpdateFn: func(screenshots []models.Screenshot, whitelist []string) ([]error, error) {
						return nil, nil
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.ScreenshotsUploadedPatchResponse{
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
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						require.Equal(t, appVersion.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
						return []models.Screenshot{
							models.Screenshot{
								Filename:   "screenshot.png",
								DeviceType: "iPhone XS Max",
								ScreenSize: "6.5 inch",
								AppVersion: models.AppVersion{
									Record: models.Record{ID: appVersion.ID},
									App:    models.App{AppSlug: "test-app-slug"},
								},
							},
							models.Screenshot{
								Filename:   "screenshot2.png",
								DeviceType: "iPhone XS",
								ScreenSize: "5.5 inch",
								AppVersion: models.AppVersion{
									Record: models.Record{ID: appVersion.ID},
									App:    models.App{AppSlug: "test-app-slug"},
								},
							},
						}, nil
					},
					batchUpdateFn: func(screenshots []models.Screenshot, whitelist []string) ([]error, error) {
						appVersionID := uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879")

						require.Equal(t, []string{"Uploaded"}, whitelist)
						require.Equal(t, models.Screenshot{
							Filename:   "screenshot.png",
							DeviceType: "iPhone XS Max",
							ScreenSize: "6.5 inch",
							Uploaded:   true,
							AppVersion: models.AppVersion{
								Record: models.Record{ID: appVersionID},
								App:    models.App{AppSlug: "test-app-slug"},
							},
						}, screenshots[0])
						require.Equal(t, models.Screenshot{
							Filename:   "screenshot2.png",
							DeviceType: "iPhone XS",
							ScreenSize: "5.5 inch",
							Uploaded:   true,
							AppVersion: models.AppVersion{
								Record: models.Record{ID: appVersionID},
								App:    models.App{AppSlug: "test-app-slug"},
							},
						}, screenshots[1])

						return nil, nil
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return fmt.Sprintf("http://presigned.aws.url/%s", path), nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.ScreenshotsGetResponse{
				Data: []services.ScreenshotData{
					services.ScreenshotData{
						Screenshot: models.Screenshot{
							Filename:   "screenshot.png",
							DeviceType: "iPhone XS Max",
							ScreenSize: "6.5 inch",
							Uploaded:   true,
						},
						DownloadURL: "http://presigned.aws.url/test-app-slug/de438ddc-98e5-4226-a5f4-fd2d53474879/iPhone XS Max (6.5 inch)/screenshot.png",
					},
					services.ScreenshotData{
						Screenshot: models.Screenshot{
							Filename:   "screenshot2.png",
							DeviceType: "iPhone XS",
							ScreenSize: "5.5 inch",
							Uploaded:   true,
						},
						DownloadURL: "http://presigned.aws.url/test-app-slug/de438ddc-98e5-4226-a5f4-fd2d53474879/iPhone XS (5.5 inch)/screenshot2.png",
					},
				},
			},
		})
	})

	t.Run("when unexpected error happpens at find", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, errors.New("SOME-SQL-ERROR-AT-FIND")
					},
					batchUpdateFn: func(screenshots []models.Screenshot, whitelist []string) ([]error, error) {
						return nil, nil
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR-AT-FIND",
		})
	})

	t.Run("when validation error at screenshot update", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
					batchUpdateFn: func(screenshots []models.Screenshot, whitelist []string) ([]error, error) {
						return []error{errors.New("SOME-VALIDATION-ERROR")}, nil
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse: httpresponse.ValidationErrorRespModel{
				Message: "Unprocessable Entity",
				Errors:  []string{"SOME-VALIDATION-ERROR"},
			},
		})
	})

	t.Run("when unexpected error at update", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, nil
					},
					batchUpdateFn: func(screenshots []models.Screenshot, whitelist []string) ([]error, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when error at generating AWS presigned URL", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			env: &env.AppEnv{
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{
							models.Screenshot{},
						}, nil
					},
					batchUpdateFn: func(screenshots []models.Screenshot, whitelist []string) ([]error, error) {
						return nil, nil
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

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
	"github.com/bitrise-io/api-utils/providers"
	"github.com/c2fo/testify/require"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_ScreenshotsGetHandler(t *testing.T) {
	httpMethod := "GET"
	url := "/apps/{app-slug}/versions/{version-id}/screenshots"
	handler := services.ScreenshotsGetHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"ScreenshotService"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			ScreenshotService: &testScreenshotService{},
			AWS:               &providers.AWSMock{},
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
				},
				AWS: &providers.AWSMock{
					GeneratePresignedGETURLFn: func(path string, expiration time.Duration) (string, error) {
						return "", nil
					},
				},
			},
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
						},
						DownloadURL: "http://presigned.aws.url/test-app-slug/de438ddc-98e5-4226-a5f4-fd2d53474879/iPhone XS Max (6.5 inch)/screenshot.png",
					},
					services.ScreenshotData{
						Screenshot: models.Screenshot{
							Filename:   "screenshot2.png",
							DeviceType: "iPhone XS",
							ScreenSize: "5.5 inch",
						},
						DownloadURL: "http://presigned.aws.url/test-app-slug/de438ddc-98e5-4226-a5f4-fd2d53474879/iPhone XS (5.5 inch)/screenshot2.png",
					},
				},
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
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						return []models.Screenshot{}, errors.New("SOME-SQL-ERROR")
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
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			env: &env.AppEnv{
				ScreenshotService: &testScreenshotService{
					findAllFn: func(appVersion *models.AppVersion) ([]models.Screenshot, error) {
						require.Equal(t, appVersion.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
						return []models.Screenshot{
							models.Screenshot{},
						}, nil
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

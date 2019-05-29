package services_test

import (
	"net/http"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/bitrise-io/api-utils/providers"
	"github.com/c2fo/testify/require"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_ScreenshotDeleteHandler(t *testing.T) {
	httpMethod := "DELETE"
	url := "/apps/{app-slug}/versions/{version-id}/screenshots/{screenshot-id}"
	handler := services.ScreenshotDeleteHandler

	screenshotID := uuid.NewV4()
	testScreenshot := &models.Screenshot{Record: models.Record{ID: screenshotID}}

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"ScreenshotService"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedScreenshotID: screenshotID,
		},
		env: &env.AppEnv{},
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedScreenshotID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{},
		env:             &env.AppEnv{},
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedScreenshotID: screenshotID,
			},
			env: &env.AppEnv{
				ScreenshotService: &testScreenshotService{
					deleteFn: func(*models.Screenshot) (validationErrors []error, dbError error) {
						return nil, nil
					},
					findFn: func(screemshot *models.Screenshot) (*models.Screenshot, error) {
						require.Equal(t, screemshot.ID.String(), screenshotID.String())
						return testScreenshot, nil
					},
				},
				AWS: &providers.AWSMock{
					DeleteObjectFn: func(path string) error {
						return nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.ScreenshotDeleteResponse{
				Data: testScreenshot,
			},
		})
	})

	t.Run("error - unexpected error in database", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedScreenshotID: screenshotID,
			},
			env: &env.AppEnv{
				ScreenshotService: &testScreenshotService{
					findFn: func(screemshot *models.Screenshot) (*models.Screenshot, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
				AWS: &providers.AWSMock{},
			},
			expectedInternalErr: "SOME-SQL-ERROR",
		})
	})

	t.Run("error - aws error", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedScreenshotID: screenshotID,
			},
			env: &env.AppEnv{
				ScreenshotService: &testScreenshotService{
					findFn: func(screemshot *models.Screenshot) (*models.Screenshot, error) {
						return testScreenshot, nil
					},
				},
				AWS: &providers.AWSMock{
					DeleteObjectFn: func(path string) error {
						return errors.New("An AWS error")
					},
				},
			},
			expectedInternalErr: "An AWS error",
		})
	})
}

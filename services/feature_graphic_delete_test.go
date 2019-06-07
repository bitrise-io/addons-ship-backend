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

func Test_FeatureGraphicDeleteHandler(t *testing.T) {
	httpMethod := "DELETE"
	url := "/apps/{app-slug}/versions/{version-id}/feature-graphic"
	handler := services.FeatureGraphicDeleteHandler

	testAppVersionID := uuid.NewV4()
	testFeatureGraphic := &models.FeatureGraphic{AppVersionID: testAppVersionID}

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"FeatureGraphicService", "AWS"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
		},
		env: &env.AppEnv{
			FeatureGraphicService: &testFeatureGraphicService{
				deleteFn: func(featureGraphic *models.FeatureGraphic) error {
					return nil
				},
				findFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, error) {
					return &models.FeatureGraphic{}, nil
				},
			},
			AWS: &providers.AWSMock{},
		},
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppVersionID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
		},
		env: &env.AppEnv{},
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				FeatureGraphicService: &testFeatureGraphicService{
					deleteFn: func(featureGraphic *models.FeatureGraphic) error {
						return nil
					},
					findFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, error) {
						require.Equal(t, featureGraphic.AppVersionID, testAppVersionID)
						return &models.FeatureGraphic{}, nil
					},
				},
				AWS: &providers.AWSMock{
					DeleteObjectFn: func(path string) error {
						return nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.FeatureGraphicDeleteResponse{
				Data: testFeatureGraphic,
			},
		})
	})

	t.Run("error - unexpected error in database at find", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				FeatureGraphicService: &testFeatureGraphicService{
					deleteFn: func(featureGraphic *models.FeatureGraphic) error {
						return nil
					},
					findFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
				AWS: &providers.AWSMock{},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("error - unexpected error in database at delete", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				FeatureGraphicService: &testFeatureGraphicService{
					findFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, error) {
						return &models.FeatureGraphic{}, nil
					},
					deleteFn: func(featureGraphic *models.FeatureGraphic) error {
						return errors.New("SOME-SQL-ERROR")
					},
				},
				AWS: &providers.AWSMock{
					DeleteObjectFn: func(path string) error {
						return nil
					},
				},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("error - aws error", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				FeatureGraphicService: &testFeatureGraphicService{
					deleteFn: func(featureGraphic *models.FeatureGraphic) error {
						return nil
					},
					findFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, error) {
						require.Equal(t, featureGraphic.AppVersionID, testAppVersionID)
						return &models.FeatureGraphic{}, nil
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

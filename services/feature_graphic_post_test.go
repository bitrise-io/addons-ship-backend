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

func Test_FeatureGraphicPostHandler(t *testing.T) {
	httpMethod := "POST"
	url := "/apps/{app-slug}/versions/{version-id}/feature-graphic"
	handler := services.FeatureGraphicPostHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"FeatureGraphicService", "AWS"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			FeatureGraphicService: &testFeatureGraphicService{
				createFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, []error, error) {
					return &models.FeatureGraphic{}, nil, nil
				},
			},
			AWS: &providers.AWSMock{},
		},
		requestBody: `{}`,
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppVersionID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			FeatureGraphicService: &testFeatureGraphicService{},
			AWS: &providers.AWSMock{},
		},
		requestBody: `{}`,
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				FeatureGraphicService: &testFeatureGraphicService{
					createFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, []error, error) {
						return &models.FeatureGraphic{}, nil, nil
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedPUTURLFn: func(path string, expiration time.Duration, size int64) (string, error) {
						return fmt.Sprintf("http://presigned.aws.url/%s", path), nil
					},
				},
			},
			requestBody:        `{}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.FeatureGraphicPostResponse{
				Data: services.FeatureGraphicData{
					UploadURL: "http://presigned.aws.url//00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000",
				},
			},
		})
	})

	t.Run("ok - more complex", func(t *testing.T) {
		testFeatureGraphicUUID := uuid.FromStringOrNil("33c7223f-2203-4109-b439-6026e7a374c9")

		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			env: &env.AppEnv{
				FeatureGraphicService: &testFeatureGraphicService{
					createFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, []error, error) {
						appVersionID := uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879")

						require.Equal(t, &models.FeatureGraphic{
							AppVersionID: appVersionID,
							UploadableObject: models.UploadableObject{
								Filename: "feature_graphic.png",
								Filesize: 1234,
							},
						}, featureGraphic)

						return &models.FeatureGraphic{
							Record:           models.Record{ID: testFeatureGraphicUUID},
							UploadableObject: models.UploadableObject{Filename: "feature_graphic.png"},
							AppVersion: models.AppVersion{
								Record: models.Record{ID: appVersionID},
								App:    models.App{AppSlug: "test-app-slug"},
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
			requestBody:        `{"filename":"feature_graphic.png","filesize":1234}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.FeatureGraphicPostResponse{
				Data: services.FeatureGraphicData{
					FeatureGraphic: models.FeatureGraphic{
						Record:           models.Record{ID: testFeatureGraphicUUID},
						UploadableObject: models.UploadableObject{Filename: "feature_graphic.png"},
					},
					UploadURL: "http://presigned.aws.url/test-app-slug/de438ddc-98e5-4226-a5f4-fd2d53474879/33c7223f-2203-4109-b439-6026e7a374c9.png",
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
				FeatureGraphicService: &testFeatureGraphicService{
					createFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, []error, error) {
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

	t.Run("error - validation error at feature graphic create", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				FeatureGraphicService: &testFeatureGraphicService{
					createFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, []error, error) {
						return nil, []error{errors.New("SOME-VALIDATION-ERROR")}, nil
					},
				},
				AWS: &providers.AWSMock{
					GeneratePresignedPUTURLFn: func(path string, expiration time.Duration, size int64) (string, error) {
						return "", nil
					},
				},
			},
			requestBody:        `{}`,
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
				FeatureGraphicService: &testFeatureGraphicService{
					createFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, []error, error) {
						return nil, nil, errors.New("SOME-SQL-ERROR")
					},
				},
				AWS: &providers.AWSMock{},
			},
			requestBody:         `{}`,
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("error - when generating AWS presigned URL", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			env: &env.AppEnv{
				FeatureGraphicService: &testFeatureGraphicService{
					createFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, []error, error) {
						appVersionID := uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879")
						return &models.FeatureGraphic{
							Record:           models.Record{ID: uuid.FromStringOrNil("33c7223f-2203-4109-b439-6026e7a374c9")},
							UploadableObject: models.UploadableObject{Filename: "feature_graphic.png"},
							AppVersion: models.AppVersion{
								Record: models.Record{ID: appVersionID},
								App:    models.App{AppSlug: "test-app-slug"},
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
			requestBody:         `{}`,
			expectedInternalErr: "SOME-AWS-ERROR",
		})
	})
}

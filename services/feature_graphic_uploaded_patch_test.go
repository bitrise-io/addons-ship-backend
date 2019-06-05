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

func Test_FeatureGraphicUploadedPatchHandler(t *testing.T) {
	httpMethod := "PATCH"
	url := "/apps/{app-slug}/versions/{version-id}/feature-graphic"
	handler := services.FeatureGraphicUploadedPatchHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"FeatureGraphicService", "AWS"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			FeatureGraphicService: &testFeatureGraphicService{
				findFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, error) {
					return &models.FeatureGraphic{}, nil
				},
				updateFn: func(featureGraphic models.FeatureGraphic, whitelist []string) ([]error, error) {
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
			FeatureGraphicService: &testFeatureGraphicService{},
			AWS: &providers.AWSMock{},
		},
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				FeatureGraphicService: &testFeatureGraphicService{
					findFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, error) {
						return &models.FeatureGraphic{}, nil
					},
					updateFn: func(featureGraphic models.FeatureGraphic, whitelist []string) ([]error, error) {
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
			expectedResponse: services.FeatureGraphicUploadedPatchResponse{
				Data: services.FeatureGraphicData{
					FeatureGraphic: models.FeatureGraphic{UploadableObject: models.UploadableObject{Uploaded: true}},
				},
			},
		})
	})

	t.Run("ok - more complex", func(t *testing.T) {
		testAppVersionID := uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879")
		testFeatureGraphicUUID := uuid.FromStringOrNil("33c7223f-2203-4109-b439-6026e7a374c9")

		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				FeatureGraphicService: &testFeatureGraphicService{
					findFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, error) {
						require.Equal(t, featureGraphic.AppVersionID, testAppVersionID)
						return &models.FeatureGraphic{
							Record:           models.Record{ID: testFeatureGraphicUUID},
							UploadableObject: models.UploadableObject{Filename: "feature_graphic.png"},
							AppVersion: models.AppVersion{
								Record: models.Record{ID: featureGraphic.AppVersionID},
								App:    models.App{AppSlug: "test-app-slug"},
							},
						}, nil
					},
					updateFn: func(featureGraphic models.FeatureGraphic, whitelist []string) ([]error, error) {
						require.Equal(t, []string{"Uploaded"}, whitelist)
						require.Equal(t, models.FeatureGraphic{
							Record: models.Record{ID: testFeatureGraphicUUID},
							UploadableObject: models.UploadableObject{
								Filename: "feature_graphic.png",
								Uploaded: true,
							},
							AppVersion: models.AppVersion{
								Record: models.Record{ID: testAppVersionID},
								App:    models.App{AppSlug: "test-app-slug"},
							},
						}, featureGraphic)

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
			expectedResponse: services.FeatureGraphicUploadedPatchResponse{
				Data: services.FeatureGraphicData{
					FeatureGraphic: models.FeatureGraphic{
						Record: models.Record{ID: testFeatureGraphicUUID},
						UploadableObject: models.UploadableObject{
							Filename: "feature_graphic.png",
							Uploaded: true,
						},
					},
					DownloadURL: "http://presigned.aws.url/test-app-slug/de438ddc-98e5-4226-a5f4-fd2d53474879/33c7223f-2203-4109-b439-6026e7a374c9.png",
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
				FeatureGraphicService: &testFeatureGraphicService{
					findFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, error) {
						return nil, errors.New("SOME-SQL-ERROR-AT-FIND")
					},
					updateFn: func(featureGraphic models.FeatureGraphic, whitelist []string) ([]error, error) {
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

	t.Run("when validation error at feature graphic update", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				FeatureGraphicService: &testFeatureGraphicService{
					findFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, error) {
						return &models.FeatureGraphic{}, nil
					},
					updateFn: func(featureGraphic models.FeatureGraphic, whitelist []string) ([]error, error) {
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
				FeatureGraphicService: &testFeatureGraphicService{
					findFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, error) {
						return &models.FeatureGraphic{}, nil
					},
					updateFn: func(featureGraphic models.FeatureGraphic, whitelist []string) ([]error, error) {
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
				FeatureGraphicService: &testFeatureGraphicService{
					findFn: func(featureGraphic *models.FeatureGraphic) (*models.FeatureGraphic, error) {
						return &models.FeatureGraphic{}, nil
					},
					updateFn: func(featureGraphic models.FeatureGraphic, whitelist []string) ([]error, error) {
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

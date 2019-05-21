package services_test

import (
	"net/http"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_AppVersionGetHandler(t *testing.T) {
	httpMethod := "GET"
	url := "/apps/{app-slug}/version{version-id}"
	handler := services.AppVersionGetHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppVersionService", "BitriseAPI"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppVersionService: &testAppVersionService{
				findFn: func(*models.AppVersion) (*models.AppVersion, error) {
					return nil, nil
				},
			},
			BitriseAPI: &testBitriseAPI{},
		},
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppVersionID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppVersionService: &testAppVersionService{},
			BitriseAPI:        &testBitriseAPI{},
		},
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						require.Equal(t, appVersion.ID.String(), "de438ddc-98e5-4226-a5f4-fd2d53474879")
						return &models.AppVersion{App: models.App{}}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactMetadataFn: func(string, string, string) (*bitrise.ArtifactMeta, error) {
						return &bitrise.ArtifactMeta{}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionGetResponse{
				Data: services.AppVersionData{AppVersion: &models.AppVersion{}},
			},
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
						return &models.AppVersion{
							Version:  "v1.0",
							Platform: "ios",
							App: models.App{
								BitriseAPIToken: "test-api-token",
							},
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactMetadataFn: func(string, string, string) (*bitrise.ArtifactMeta, error) {
						return &bitrise.ArtifactMeta{
							AppInfo: bitrise.AppInfo{
								MinimumOS: "11.1",
							},
						}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionGetResponse{
				Data: services.AppVersionData{
					AppVersion: &models.AppVersion{
						Version:  "v1.0",
						Platform: "ios",
					},
					MinimumOS: "11.1",
				},
			},
		})
	})

	t.Run("ok - more complex - device family", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						return &models.AppVersion{
							Version:  "v1.0",
							Platform: "ios",
							App: models.App{
								BitriseAPIToken: "test-api-token",
							},
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactMetadataFn: func(string, string, string) (*bitrise.ArtifactMeta, error) {
						return &bitrise.ArtifactMeta{
							AppInfo: bitrise.AppInfo{
								MinimumOS:        "11.1",
								DeviceFamilyList: []int{1, 2},
							},
						}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionGetResponse{
				Data: services.AppVersionData{
					AppVersion: &models.AppVersion{
						Version:  "v1.0",
						Platform: "ios",
					},
					MinimumOS:            "11.1",
					SupportedDeviceTypes: []string{"iPhone", "iPod Touch", "iPad"},
				},
			},
		})
	})

	t.Run("ok - more complex - unknown device family", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						return &models.AppVersion{
							Version:  "v1.0",
							Platform: "ios",
							App: models.App{
								BitriseAPIToken: "test-api-token",
							},
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactMetadataFn: func(string, string, string) (*bitrise.ArtifactMeta, error) {
						return &bitrise.ArtifactMeta{
							AppInfo: bitrise.AppInfo{
								MinimumOS:        "11.1",
								DeviceFamilyList: []int{12},
							},
						}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionGetResponse{
				Data: services.AppVersionData{
					AppVersion: &models.AppVersion{
						Version:  "v1.0",
						Platform: "ios",
					},
					MinimumOS:            "11.1",
					SupportedDeviceTypes: []string{"Unknown"},
				},
			},
		})
	})

	t.Run("error - not found in database", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						return nil, gorm.ErrRecordNotFound
					},
				},
				BitriseAPI: &testBitriseAPI{},
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: httpresponse.StandardErrorRespModel{
				Message: "Not Found",
			},
		})
	})

	t.Run("error - unexpected error in database", func(t *testing.T) {
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
				BitriseAPI: &testBitriseAPI{},
			},
			expectedStatusCode:  http.StatusNotFound,
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})
}

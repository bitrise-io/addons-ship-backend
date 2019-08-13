package services_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/bitrise-io/go-utils/pointers"
	"github.com/c2fo/testify/require"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_AppVersionsGetHandler(t *testing.T) {
	httpMethod := "GET"
	url := "/apps/{app-slug}/app-versions"
	handler := services.AppVersionsGetHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppVersionService", "BitriseAPI"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppVersionService: &testAppVersionService{
				findAllFn: func(app *models.App, filterParams map[string]interface{}) ([]models.AppVersion, error) {
					return []models.AppVersion{}, nil
				},
			},
			BitriseAPI: &testBitriseAPI{},
		},
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppVersionService: &testAppVersionService{},
			BitriseAPI:        &testBitriseAPI{},
		},
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findAllFn: func(app *models.App, filterParams map[string]interface{}) ([]models.AppVersion, error) {
						require.Equal(t, app.ID.String(), "211afc15-127a-40f9-8cbe-1dadc1f86cdf")
						require.Equal(t, filterParams, map[string]interface{}{})
						return []models.AppVersion{}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{}, nil
					},
					getArtifactPublicPageURLFn: func(string, string, string, string) (string, error) {
						return "", nil
					},
					getAppDetailsFn: func(string, string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionsGetResponse{
				Data: []services.AppVersionsGetResponseElement{},
			},
		})
	})

	t.Run("ok - more complex", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findAllFn: func(app *models.App, filterParams map[string]interface{}) ([]models.AppVersion, error) {
						return []models.AppVersion{
							models.AppVersion{
								Platform:         "ios",
								ArtifactInfoData: json.RawMessage(`{"version":"v1.0"}`),
							},
							models.AppVersion{
								Platform:         "android",
								ArtifactInfoData: json.RawMessage(`{"version":"v1.12"}`),
							},
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{}, nil
					},
					getArtifactPublicPageURLFn: func(string, string, string, string) (string, error) {
						return "", nil
					},
					getAppDetailsFn: func(string, string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{
							Title:       "The Adventures of Stealy",
							AvatarURL:   pointers.NewStringPtr("https://bit.ly/1LixVJu"),
							ProjectType: "other",
						}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionsGetResponse{
				Data: []services.AppVersionsGetResponseElement{
					services.AppVersionsGetResponseElement{
						AppVersion: models.AppVersion{
							Platform: "ios",
						},
						Version: "v1.0",
						AppInfo: services.AppData{
							Title:       "The Adventures of Stealy",
							AppIconURL:  pointers.NewStringPtr("https://bit.ly/1LixVJu"),
							ProjectType: "other",
						},
					},
					services.AppVersionsGetResponseElement{
						AppVersion: models.AppVersion{
							Platform: "android",
						},
						Version: "v1.12",
						AppInfo: services.AppData{
							Title:       "The Adventures of Stealy",
							AppIconURL:  pointers.NewStringPtr("https://bit.ly/1LixVJu"),
							ProjectType: "other",
						},
					},
				},
			},
		})
	})

	t.Run("ok - with platform filter", func(t *testing.T) {
		urlWithFilter := url + "?platform=ios"
		performControllerTest(t, httpMethod, urlWithFilter, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findAllFn: func(app *models.App, filterParams map[string]interface{}) ([]models.AppVersion, error) {
						require.Equal(t, app.ID.String(), "211afc15-127a-40f9-8cbe-1dadc1f86cdf")
						require.Equal(t, filterParams, map[string]interface{}{
							"platform": "ios",
						})
						return []models.AppVersion{
							models.AppVersion{
								ArtifactInfoData: json.RawMessage(`{"version":"v1.0"}`),
								Platform:         "ios",
							},
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{}, nil
					},
					getArtifactPublicPageURLFn: func(string, string, string, string) (string, error) {
						return "", nil
					},
					getAppDetailsFn: func(string, string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionsGetResponse{
				Data: []services.AppVersionsGetResponseElement{
					services.AppVersionsGetResponseElement{
						AppVersion: models.AppVersion{
							Platform: "ios",
						},
						Version: "v1.0",
					},
				},
			},
		})
	})

	t.Run("error - unexpected error in database", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findAllFn: func(app *models.App, filterParams map[string]interface{}) ([]models.AppVersion, error) {
						return []models.AppVersion{}, errors.New("SOME-SQL-ERROR")
					},
				},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when invalid JSON is stored in database for artifact info", func(t *testing.T) {
		urlWithFilter := url + "?platform=ios"
		performControllerTest(t, httpMethod, urlWithFilter, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findAllFn: func(app *models.App, filterParams map[string]interface{}) ([]models.AppVersion, error) {
						require.Equal(t, app.ID.String(), "211afc15-127a-40f9-8cbe-1dadc1f86cdf")
						require.Equal(t, filterParams, map[string]interface{}{
							"platform": "ios",
						})
						return []models.AppVersion{
							models.AppVersion{
								ArtifactInfoData: json.RawMessage(`invalid JSON`),
								Platform:         "ios",
							},
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{}, nil
					},
					getArtifactPublicPageURLFn: func(string, string, string, string) (string, error) {
						return "", nil
					},
					getAppDetailsFn: func(string, string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
				},
			},
			expectedInternalErr: "invalid character 'i' looking for beginning of value",
		})
	})

	t.Run("when error happens at fetching app data from Bitrise API", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID:        uuid.FromStringOrNil("211afc15-127a-40f9-8cbe-1dadc1f86cdf"),
				services.ContextKeyAuthorizedAppVersionID: uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879"),
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findAllFn: func(app *models.App, filterParams map[string]interface{}) ([]models.AppVersion, error) {
						require.Equal(t, app.ID.String(), "211afc15-127a-40f9-8cbe-1dadc1f86cdf")
						require.Equal(t, filterParams, map[string]interface{}{})
						return []models.AppVersion{
							models.AppVersion{
								ArtifactInfoData: json.RawMessage(`{"version":"v1.0"}`),
								Platform:         "ios",
							},
						}, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getArtifactsFn: func(string, string, string) ([]bitrise.ArtifactListElementResponseModel, error) {
						return []bitrise.ArtifactListElementResponseModel{}, nil
					},
					getArtifactPublicPageURLFn: func(string, string, string, string) (string, error) {
						return "", nil
					},
					getAppDetailsFn: func(string, string) (*bitrise.AppDetails, error) {
						return nil, errors.New("SOME-BITRISE-API-ERROR")
					},
				},
			},
			expectedInternalErr: "SOME-BITRISE-API-ERROR",
		})
	})
}

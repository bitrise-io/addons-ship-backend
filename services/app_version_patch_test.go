package services_test

import (
	"encoding/json"
	"net/http"
	"testing"

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

func Test_AppVersionPutHandler(t *testing.T) {
	httpMethod := "PATCH"
	url := "/apps/{app-slug}/version{version-id}"
	handler := services.AppVersionPutHandler

	testAppVersionID := uuid.FromStringOrNil("de438ddc-98e5-4226-a5f4-fd2d53474879")

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppVersionService"}, ControllerTestCase{
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
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						return &models.AppVersion{App: models.App{}, AppStoreInfoData: json.RawMessage(`{}`)}, nil
					},
					updateFn: func(appVersion *models.AppVersion, whitelist []string) ([]error, error) {
						return nil, nil
					},
				},
			},
			requestBody:        `{}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionPutResponse{
				Data: services.AppVersionPutResponseData{AppVersion: &models.AppVersion{}},
			},
		})
	})

	t.Run("ok - more complex", func(t *testing.T) {
		expectedAppStoreInfoModel := models.AppStoreInfo{ShortDescription: "Some short description"}
		expectedAppStoreInfo, err := json.Marshal(expectedAppStoreInfoModel)
		require.NoError(t, err)

		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						require.Equal(t, appVersion.ID, testAppVersionID)
						appVersion.AppStoreInfoData = json.RawMessage(`{}`)
						appVersion.App = models.App{}
						return appVersion, nil
					},
					updateFn: func(appVersion *models.AppVersion, whitelist []string) ([]error, error) {
						require.Equal(t, testAppVersionID, appVersion.ID)
						require.Equal(t, expectedAppStoreInfo, appVersion.AppStoreInfoData)
						return nil, nil
					},
				},
			},
			requestBody:        `{"short_description":"Some short description"}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppVersionPutResponse{
				Data: services.AppVersionPutResponseData{
					AppVersion: &models.AppVersion{
						Record: models.Record{ID: testAppVersionID},
					},
					AppStoreInfo: models.AppStoreInfo{
						ShortDescription: "Some short description",
					},
				},
			},
		})
	})

	t.Run("when request body is not a valid JSON", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						return &models.AppVersion{App: models.App{}, AppStoreInfoData: json.RawMessage(`{}`)}, nil
					},
					updateFn: func(appVersion *models.AppVersion, whitelist []string) ([]error, error) {
						return nil, nil
					},
				},
			},
			requestBody:        `invalid-request-body`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: httpresponse.StandardErrorRespModel{
				Message: "Invalid request body, JSON decode failed",
			},
		})
	})

	t.Run("when app version to update not found", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						return nil, gorm.ErrRecordNotFound
					},
					updateFn: func(appVersion *models.AppVersion, whitelist []string) ([]error, error) {
						return nil, nil
					},
				},
			},
			requestBody:        `{}`,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Not Found"},
		})
	})

	t.Run("when db error happens at finding the app version to update", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
					updateFn: func(appVersion *models.AppVersion, whitelist []string) ([]error, error) {
						return nil, nil
					},
				},
			},
			requestBody:         `{}`,
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when validation error happens at update", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						return &models.AppVersion{App: models.App{}, AppStoreInfoData: json.RawMessage(`{}`)}, nil
					},
					updateFn: func(appVersion *models.AppVersion, whitelist []string) ([]error, error) {
						return []error{errors.New("SOME-VALIDATION-ERROR")}, nil
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

	t.Run("when unexpected db error happens at update", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						return &models.AppVersion{App: models.App{}, AppStoreInfoData: json.RawMessage(`{}`)}, nil
					},
					updateFn: func(appVersion *models.AppVersion, whitelist []string) ([]error, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
			},
			requestBody:         `{}`,
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when app store info data contains an invalid JSON", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppVersionID: testAppVersionID,
			},
			env: &env.AppEnv{
				AppVersionService: &testAppVersionService{
					findFn: func(appVersion *models.AppVersion) (*models.AppVersion, error) {
						return &models.AppVersion{App: models.App{}, AppStoreInfoData: json.RawMessage(`invalid json`)}, nil
					},
					updateFn: func(appVersion *models.AppVersion, whitelist []string) ([]error, error) {
						appVersion.AppStoreInfoData = json.RawMessage(`invalid json`)
						return nil, nil
					},
				},
			},
			requestBody:         `{}`,
			expectedInternalErr: "invalid character 'i' looking for beginning of value",
		})
	})
}

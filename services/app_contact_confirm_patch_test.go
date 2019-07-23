package services_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	uuid "github.com/satori/go.uuid"
)

func Test_AppContactConfirmPatchHandler(t *testing.T) {
	httpMethod := "PATCH"
	url := "/confirm_email"
	handler := services.AppContactConfirmPatchHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppContactService"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppContactService: &testAppContactService{},
		},
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppContactID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppContactService: &testAppContactService{},
		},
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
				},
				AppContactService: &testAppContactService{
					findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
						return &models.AppContact{App: &models.App{}}, nil
					},
					updateFn: func(appContact *models.AppContact, whitelist []string) error {
						appContact.ConfirmedAt = time.Time{}
						return nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppContactPatchResponse{
				Data: services.AppContactPatchResponseData{
					AppContact: &models.AppContact{},
					App:        services.AppResponseData{},
				},
			},
		})
	})

	t.Run("ok - more complex", func(t *testing.T) {
		testApp := models.App{AppSlug: "an-app-slug", APIToken: "test-api-token", Plan: "gold"}
		testAppDetails := bitrise.AppDetails{Title: "Supe Duper App"}

		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppContactID: uuid.FromStringOrNil("8a230385-0113-4cf3-a9c6-469a313e587a"),
			},
			env: &env.AppEnv{
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						require.Equal(t, "test-api-token", apiToken)
						require.Equal(t, "an-app-slug", appSlug)
						return &testAppDetails, nil
					},
				},
				AppContactService: &testAppContactService{
					findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
						require.Equal(t, uuid.FromStringOrNil("8a230385-0113-4cf3-a9c6-469a313e587a"), appContact.ID)
						appContact.App = &testApp
						return appContact, nil
					},
					updateFn: func(appContact *models.AppContact, whitelist []string) error {
						require.Equal(t, uuid.FromStringOrNil("8a230385-0113-4cf3-a9c6-469a313e587a"), appContact.ID)
						require.Nil(t, appContact.ConfirmationToken)
						require.Equal(t, []string{"ConfirmedAt", "ConfirmationToken"}, whitelist)
						require.NotEqual(t, time.Time{}, appContact.ConfirmedAt)
						appContact.ConfirmedAt = time.Time{}
						return nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppContactPatchResponse{
				Data: services.AppContactPatchResponseData{
					AppContact: &models.AppContact{Record: models.Record{ID: uuid.FromStringOrNil("8a230385-0113-4cf3-a9c6-469a313e587a")}},
					App: services.AppResponseData{
						AppSlug:    testApp.AppSlug,
						Plan:       testApp.Plan,
						AppDetails: testAppDetails,
					},
				},
			},
		})
	})

	t.Run("when error happens at finding app contact", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				BitriseAPI: &testBitriseAPI{},
				AppContactService: &testAppContactService{
					findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
						return nil, gorm.ErrRecordNotFound
					},
					updateFn: func(appContact *models.AppContact, whitelist []string) error {
						appContact.ConfirmedAt = time.Time{}
						return nil
					},
				},
			},
			expectedInternalErr: "SQL Error: record not found",
		})
	})

	t.Run("when error happens at updating app contact", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{}, nil
					},
				},
				AppContactService: &testAppContactService{
					findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
						return &models.AppContact{App: &models.App{AppSlug: "an-app-slug", APIToken: "some-token"}}, nil
					},
					updateFn: func(appContact *models.AppContact, whitelist []string) error {
						return errors.New("SOME-SQL-ERROR")
					},
				},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when an error happens while getting app details", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return nil, errors.New("SOME-BITRISE-API-ERROR")
					},
				},
				AppContactService: &testAppContactService{
					findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
						return &models.AppContact{App: &models.App{AppSlug: "an-app-slug", APIToken: "some-token"}}, nil
					},
					updateFn: func(appContact *models.AppContact, whitelist []string) error {
						appContact.ConfirmedAt = time.Time{}
						return nil
					},
				},
			},
			expectedInternalErr: "SOME-BITRISE-API-ERROR",
		})
	})
}

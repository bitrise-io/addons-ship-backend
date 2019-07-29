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
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_AppContactPostHandler(t *testing.T) {
	httpMethod := "POST"
	url := "/apps/{app-slug}/contacts"
	handler := services.AppContactPostHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppContactService", "BitriseAPI", "Mailer"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppContactService: &testAppContactService{},
			BitriseAPI:        &testBitriseAPI{},
			Mailer:            &testMailer{},
		},
		requestBody: `{}`,
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppContactService: &testAppContactService{},
			BitriseAPI:        &testBitriseAPI{},
			Mailer:            &testMailer{},
		},
		requestBody: `{}`,
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppContactService: &testAppContactService{
					findFn: func(*models.AppContact) (*models.AppContact, error) {
						return nil, nil
					},
				},
				BitriseAPI: &testBitriseAPI{},
				Mailer:     &testMailer{},
			},
			requestBody:        `{}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   services.AppContactPostResponse{},
		})
	})

	t.Run("ok - more complex", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.FromStringOrNil("548bde58-2707-4c28-9474-4f35ba0176cb"),
			},
			env: &env.AppEnv{
				EmailConfirmLandingURL: "http://ship.bitrise.io/confirm_email",
				AppContactService: &testAppContactService{
					createFn: func(contact *models.AppContact) (*models.AppContact, error) {
						contact.App = &models.App{BitriseAPIToken: "test-api-token", AppSlug: "test-app-slug"}
						return contact, nil
					},
					findFn: func(*models.AppContact) (*models.AppContact, error) {
						return nil, gorm.ErrRecordNotFound
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						require.Equal(t, "test-api-token", apiToken)
						require.Equal(t, "test-app-slug", appSlug)
						return &bitrise.AppDetails{Title: "My awesome app"}, nil
					},
				},
				Mailer: &testMailer{
					sendEmailConfirmationFn: func(confirmURL string, contact *models.AppContact, appDetails *bitrise.AppDetails) error {
						require.Equal(t, "My awesome app", appDetails.Title)
						require.Equal(t, "http://ship.bitrise.io/confirm_email", confirmURL)
						require.NotNil(t, contact.ConfirmationToken)
						contact.ConfirmationToken = nil
						require.Equal(t, &models.AppContact{
							Email: "someones@email.addr",
							NotificationPreferencesData: json.RawMessage(`{"new_version":true,"successful_publish":false,"failed_publish":false}`),
							AppID: uuid.FromStringOrNil("548bde58-2707-4c28-9474-4f35ba0176cb"),
							App: &models.App{
								BitriseAPIToken: "test-api-token",
								AppSlug:         "test-app-slug",
							},
						}, contact)
						return nil
					},
				},
			},
			requestBody:        `{"email":"someones@email.addr","notification_preferences":{"new_version":true}}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppContactPostResponse{
				Data: &models.AppContact{
					Email: "someones@email.addr",
					NotificationPreferencesData: json.RawMessage(`{"new_version":true,"successful_publish":false,"failed_publish":false}`),
					App: &models.App{
						Record:          models.Record{ID: uuid.FromStringOrNil("548bde58-2707-4c28-9474-4f35ba0176cb")},
						BitriseAPIToken: "test-api-token",
						AppSlug:         "test-app-slug",
					},
				},
			},
		})
	})

	t.Run("when request body is invalid", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppContactService: &testAppContactService{
					findFn: func(*models.AppContact) (*models.AppContact, error) {
						return nil, nil
					},
				},
				BitriseAPI: &testBitriseAPI{},
				Mailer:     &testMailer{},
			},
			requestBody:        `invalid JSON`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Invalid request body, JSON decode failed"},
		})
	})

	t.Run("when error happens at creating new app contact", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppContactService: &testAppContactService{
					findFn: func(*models.AppContact) (*models.AppContact, error) {
						return nil, gorm.ErrRecordNotFound
					},
					createFn: func(contact *models.AppContact) (*models.AppContact, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
				BitriseAPI: &testBitriseAPI{},
				Mailer:     &testMailer{},
			},
			requestBody:         `{}`,
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when error happens at fetchin app details", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppContactService: &testAppContactService{
					findFn: func(*models.AppContact) (*models.AppContact, error) {
						return nil, gorm.ErrRecordNotFound
					},
					createFn: func(contact *models.AppContact) (*models.AppContact, error) {
						contact.App = &models.App{BitriseAPIToken: "test-api-token", AppSlug: "test-app-slug"}
						return contact, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return nil, errors.New("SOME-BITRISE-API-ERROR")
					},
				},
				Mailer: &testMailer{},
			},
			requestBody:         `{}`,
			expectedInternalErr: "SOME-BITRISE-API-ERROR",
		})
	})

	t.Run("when it's failed to send email notification", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppContactService: &testAppContactService{
					findFn: func(*models.AppContact) (*models.AppContact, error) {
						return nil, gorm.ErrRecordNotFound
					},
					createFn: func(contact *models.AppContact) (*models.AppContact, error) {
						contact.App = &models.App{BitriseAPIToken: "test-api-token", AppSlug: "test-app-slug"}
						return contact, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					getAppDetailsFn: func(apiToken, appSlug string) (*bitrise.AppDetails, error) {
						return &bitrise.AppDetails{Title: "My awesome app"}, nil
					},
				},
				Mailer: &testMailer{
					sendEmailConfirmationFn: func(confirmURL string, contact *models.AppContact, appDetails *bitrise.AppDetails) error {
						return errors.New("SOME-MAILER-ERROR")
					},
				},
			},
			requestBody:         `{}`,
			expectedInternalErr: "SOME-MAILER-ERROR",
		})
	})
}

package services_test

import (
	"net/http"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/go-crypto/crypto"
	"github.com/bitrise-io/go-utils/envutil"
	"github.com/c2fo/testify/require"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func Test_ProvisionHandler(t *testing.T) {
	httpMethod := "POST"
	url := "/provision"
	handler := services.ProvisionHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppService", "BitriseAPI"}, ControllerTestCase{
		env: &env.AppEnv{
			AppService: &testAppService{},
			BitriseAPI: &testBitriseAPI{},
		},
		requestBody: `{}`,
	})

	t.Run("ok", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return app, nil
					},
				},
				BitriseAPI: &testBitriseAPI{},
			},
			requestBody:        `{}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.ProvisionPostResponse{
				Envs: []services.Env{
					services.Env{Key: "ADDON_SHIP_API_URL"},
					services.Env{Key: "ADDON_SHIP_API_TOKEN"},
				},
			},
		})
	})

	t.Run("ok when app exists", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						require.Equal(t, "test-app-slug", app.AppSlug)
						return app, nil
					},
				},
				BitriseAPI: &testBitriseAPI{},
			},
			requestBody:        `{"app_slug":"test-app-slug","api_token":"test-bitrise-api-token","plan":"free"}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.ProvisionPostResponse{
				Envs: []services.Env{
					services.Env{Key: "ADDON_SHIP_API_URL"},
					services.Env{Key: "ADDON_SHIP_API_TOKEN"},
				},
			},
		})
	})

	t.Run("ok when app not exists", func(t *testing.T) {
		revokeFn, err := envutil.RevokableSetenv("APP_WEBHOOK_SECRET_ENCRYPT_KEY", "06042e86a7bd421c642c8c3e4ab13840")
		require.NoError(t, err)

		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			env: &env.AppEnv{
				AddonHostURL: "http://ship.addon.url",
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						require.Equal(t, "test-app-slug", app.AppSlug)
						return nil, gorm.ErrRecordNotFound
					},
					createFn: func(app *models.App) (*models.App, error) {
						require.Equal(t, "test-app-slug", app.AppSlug)
						require.Equal(t, "test-bitrise-api-token", app.BitriseAPIToken)
						require.Equal(t, "free", app.Plan)
						require.NotEmpty(t, app.APIToken)
						app.APIToken = "test-api-token"

						iv, err := crypto.GenerateIV()
						require.NoError(t, err)
						encryptedSecret, err := crypto.AES256GCMCipher("my-super-secret", iv, "06042e86a7bd421c642c8c3e4ab13840")
						require.NoError(t, err)

						app.EncryptedSecret = encryptedSecret
						app.EncryptedSecretIV = iv
						return app, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					registerWebhookFn: func(authToken, appSlug, secret, callbackURL string) error {
						require.Equal(t, "test-app-slug", appSlug)
						require.Equal(t, "my-super-secret", secret)
						require.Equal(t, "test-bitrise-api-token", authToken)
						return nil
					},
				},
			},
			requestBody:        `{"app_slug":"test-app-slug","api_token":"test-bitrise-api-token","plan":"free"}`,
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.ProvisionPostResponse{
				Envs: []services.Env{
					services.Env{Key: "ADDON_SHIP_API_URL", Value: "http://ship.addon.url"},
					services.Env{Key: "ADDON_SHIP_API_TOKEN", Value: "test-api-token"},
				},
			},
		})

		require.NoError(t, revokeFn())
	})

	t.Run("when request body is invalid", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return app, nil
					},
				},
				BitriseAPI: &testBitriseAPI{},
			},
			requestBody:        `invalid JSON`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   httpresponse.StandardErrorRespModel{Message: "Invalid request body, JSON decode failed"},
		})
	})

	t.Run("when database error happest at find", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
				BitriseAPI: &testBitriseAPI{},
			},
			requestBody:         `{}`,
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when database error happens at create", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			env: &env.AppEnv{
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						return nil, gorm.ErrRecordNotFound
					},
					createFn: func(app *models.App) (*models.App, error) {
						return nil, errors.New("SOME-SQL-ERROR")
					},
				},
				BitriseAPI: &testBitriseAPI{},
			},
			requestBody:         `{}`,
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})

	t.Run("when it's failed to get secret from app", func(t *testing.T) {
		revokeFn, err := envutil.RevokableSetenv("APP_WEBHOOK_SECRET_ENCRYPT_KEY", "06042e86a7bd421c642c8c3e4ab13840")
		require.NoError(t, err)

		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			env: &env.AppEnv{
				AddonHostURL: "http://ship.addon.url",
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						require.Equal(t, "test-app-slug", app.AppSlug)
						return nil, gorm.ErrRecordNotFound
					},
					createFn: func(app *models.App) (*models.App, error) {
						app.APIToken = "test-api-token"

						iv, err := crypto.GenerateIV()
						require.NoError(t, err)

						app.EncryptedSecretIV = iv
						return app, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					registerWebhookFn: func(authToken, appSlug, secret, callbackURL string) error {
						return nil
					},
				},
			},
			requestBody:         `{"app_slug":"test-app-slug","api_token":"test-bitrise-api-token","plan":"free"}`,
			expectedInternalErr: "cipher: message authentication failed",
		})

		require.NoError(t, revokeFn())
	})

	t.Run("when failed to register webhook", func(t *testing.T) {
		revokeFn, err := envutil.RevokableSetenv("APP_WEBHOOK_SECRET_ENCRYPT_KEY", "06042e86a7bd421c642c8c3e4ab13840")
		require.NoError(t, err)

		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			env: &env.AppEnv{
				AddonHostURL: "http://ship.addon.url",
				AppService: &testAppService{
					findFn: func(app *models.App) (*models.App, error) {
						require.Equal(t, "test-app-slug", app.AppSlug)
						return nil, gorm.ErrRecordNotFound
					},
					createFn: func(app *models.App) (*models.App, error) {
						iv, err := crypto.GenerateIV()
						require.NoError(t, err)
						encryptedSecret, err := crypto.AES256GCMCipher("my-super-secret", iv, "06042e86a7bd421c642c8c3e4ab13840")
						require.NoError(t, err)

						app.EncryptedSecret = encryptedSecret
						app.EncryptedSecretIV = iv
						return app, nil
					},
				},
				BitriseAPI: &testBitriseAPI{
					registerWebhookFn: func(authToken, appSlug, secret, callbackURL string) error {
						return errors.New("SOME-BITRISE-API-ERROR")
					},
				},
			},
			requestBody:         `{"app_slug":"test-app-slug","api_token":"test-bitrise-api-token","plan":"free"}`,
			expectedInternalErr: "SOME-BITRISE-API-ERROR",
		})

		require.NoError(t, revokeFn())
	})
}

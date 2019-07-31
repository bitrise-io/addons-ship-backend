package services

import (
	"net/http"
	"os"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// AuthenticateWithAddonAccessTokenHandlerFunc ...
func AuthenticateWithAddonAccessTokenHandlerFunc(env *env.AppEnv, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authentication")
		if authToken == "" {
			httpresponse.RespondWithUnauthorizedNoErr(w)
			return
		}

		if env.AddonAccessToken == "" {
			httpresponse.RespondWithInternalServerError(w, errors.New("No Addon Access Token set"))
			return
		}

		if authToken != env.AddonAccessToken {
			httpresponse.RespondWithUnauthorizedNoErr(w)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// AuthenticateWithDENSecretHandlerFunc ...
func AuthenticateWithDENSecretHandlerFunc(env *env.AppEnv, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Bitrise-Den-Webhook-Secret")
		if authToken == "" {
			httpresponse.RespondWithUnauthorizedNoErr(w)
			return
		}
		denAdminSecret, ok := os.LookupEnv("BITRISE_DEN_WEBHOOK_SECRET")
		if !ok || denAdminSecret == "" {
			httpresponse.RespondWithInternalServerError(w, errors.New("No value set for BITRISE_DEN_WEBHOOK_SECRET env"))
			return
		}

		if authToken != denAdminSecret {
			httpresponse.RespondWithUnauthorizedNoErr(w)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// AuthenticateWithSSOTokenHandlerFunc ...
func AuthenticateWithSSOTokenHandlerFunc(env *env.AppEnv, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if env.AppService == nil {
			httpresponse.RespondWithInternalServerError(w, errors.New("No App Service defined for handler"))
			return
		}
		if env.SsoTokenVerifier == nil {
			httpresponse.RespondWithInternalServerError(w, errors.New("No SSO Token Verifier defined for handler"))
			return
		}

		timestamp := r.FormValue("timestamp")
		token := r.FormValue("token")
		appSlug := r.FormValue("app_slug")

		logger := env.Logger
		logger.Info("Login form data",
			zap.String("timestamp", timestamp),
			zap.String("token", token),
			zap.String("app_slug", appSlug),
		)

		valid, err := env.SsoTokenVerifier.Verify(timestamp, token, appSlug)
		if err != nil {
			httpresponse.RespondWithInternalServerError(w, errors.WithStack(err))
			return
		}
		if !valid {
			httpresponse.RespondWithNotFoundErrorNoErr(w)
			return
		}

		app, err := env.AppService.Find(&models.App{AppSlug: appSlug})
		switch {
		case errors.Cause(err) == gorm.ErrRecordNotFound:
			httpresponse.RespondWithNotFoundErrorNoErr(w)
			return
		case err != nil:
			httpresponse.RespondWithInternalServerError(w, errors.Wrap(err, "SQL Error"))
			return
		}

		ctx := ContextWithAuthorizedAppID(r.Context(), app.ID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

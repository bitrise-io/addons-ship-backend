package services

import (
	"net/http"
	"os"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
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

// AuthenticateWithDENSecretnHandlerFunc ...
func AuthenticateWithDENSecretnHandlerFunc(env *env.AppEnv, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		denAuthHeaderKey, ok := os.LookupEnv("BITRISE_DEN_SERVER_ADMIN_SECRET_HEADER_KEY")
		if !ok || denAuthHeaderKey == "" {
			httpresponse.RespondWithInternalServerError(w, errors.New("No value set for BITRISE_DEN_SERVER_ADMIN_SECRET_HEADER_KEY env"))
			return
		}
		authToken := r.Header.Get(denAuthHeaderKey)
		if authToken == "" {
			httpresponse.RespondWithUnauthorizedNoErr(w)
			return
		}
		denAdminSecret, ok := os.LookupEnv("BITRISE_DEN_SERVER_ADMIN_SECRET")
		if !ok || denAdminSecret == "" {
			httpresponse.RespondWithInternalServerError(w, errors.New("No value set for BITRISE_DEN_SERVER_ADMIN_SECRET env"))
			return
		}

		if authToken != denAdminSecret {
			httpresponse.RespondWithUnauthorizedNoErr(w)
			return
		}

		h.ServeHTTP(w, r)
	})
}

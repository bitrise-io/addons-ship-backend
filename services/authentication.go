package services

import (
	"net/http"

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
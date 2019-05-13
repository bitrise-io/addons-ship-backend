package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httprequest"
	"github.com/bitrise-io/api-utils/httpresponse"
)

// AuthorizeForAppAccessHandlerFunc ...
func AuthorizeForAppAccessHandlerFunc(env *env.AppEnv, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken, err := httprequest.AuthTokenFromHeader(r.Header)
		if err != nil {
			httpresponse.RespondWithUnauthorizedNoErr(w)
			return
		}
		urlVars := env.RequestParams.Get(r)
		appSlug := urlVars["app-slug"]
		if appSlug == "" {
			httpresponse.RespondWithBadRequestErrorNoErr(w, "App Slug not provided")
			return
		}

		app, err := env.AppService.Find(&models.App{AppSlug: appSlug, APIToken: authToken})
		if err != nil {
			httpresponse.RespondWithNotFoundErrorNoErr(w)
			return
		}

		// Access granted
		ctx := ContextWithAuthorizedAppID(r.Context(), app.ID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httprequest"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// AuthorizeForAppAccessHandlerFunc ...
func AuthorizeForAppAccessHandlerFunc(env *env.AppEnv, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken, err := httprequest.AuthTokenFromHeader(r.Header)
		if err != nil {
			httpresponse.RespondWithUnauthorizedNoErr(w)
			return
		}
		if env.RequestParams == nil {
			httpresponse.RespondWithInternalServerError(w, errors.New("No Request Params provided"))
			return
		}
		urlVars := env.RequestParams.Get(r)
		appSlug := urlVars["app-slug"]
		if appSlug == "" {
			httpresponse.RespondWithBadRequestErrorNoErr(w, "App Slug not provided")
			return
		}

		if env.AppService == nil {
			httpresponse.RespondWithInternalServerError(w, errors.New("No App Service provided"))
			return
		}

		app, err := env.AppService.Find(&models.App{AppSlug: appSlug, APIToken: authToken})
		switch {
		case errors.Cause(err) == gorm.ErrRecordNotFound:
			httpresponse.RespondWithNotFoundErrorNoErr(w)
			return
		case err != nil:
			httpresponse.RespondWithInternalServerError(w, errors.WithStack(err))
			return
		}

		// Access granted
		ctx := ContextWithAuthorizedAppID(r.Context(), app.ID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthorizeForAppVersionAccessHandlerFunc ...
func AuthorizeForAppVersionAccessHandlerFunc(env *env.AppEnv, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if env.RequestParams == nil {
			httpresponse.RespondWithInternalServerError(w, errors.New("No Request Params provided"))
			return
		}
		urlVars := env.RequestParams.Get(r)
		appVersionParam := urlVars["version-id"]
		if appVersionParam == "" {
			httpresponse.RespondWithBadRequestErrorNoErr(w, "App Version ID not provided")
			return
		}
		appVersionID, err := uuid.FromString(appVersionParam)
		if err != nil {
			httpresponse.RespondWithBadRequestErrorNoErr(w, "App Version ID has invalid format")
			return
		}

		if env.AppVersionService == nil {
			httpresponse.RespondWithInternalServerError(w, errors.New("No App Version Service provided"))
			return
		}
		appVersion, err := env.AppVersionService.Find(&models.AppVersion{Record: models.Record{ID: appVersionID}})
		switch {
		case errors.Cause(err) == gorm.ErrRecordNotFound:
			httpresponse.RespondWithNotFoundErrorNoErr(w)
			return
		case err != nil:
			httpresponse.RespondWithInternalServerError(w, errors.WithStack(err))
			return
		}

		// Access granted
		ctx := ContextWithAuthorizedAppVersionID(r.Context(), appVersion.ID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

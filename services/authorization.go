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

		appVersionID, err := getUUIDFromRequest(env, r, "version-id")
		if err != nil {
			httpresponse.RespondWithBadRequestErrorNoErr(w, err.Error())
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

// AuthorizeForAppVersionScreenshotAccessHandlerFunc ...
func AuthorizeForAppVersionScreenshotAccessHandlerFunc(env *env.AppEnv, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if env.RequestParams == nil {
			httpresponse.RespondWithInternalServerError(w, errors.New("No Request Params provided"))
			return
		}

		appVersionID, err := GetAuthorizedAppVersionIDFromContext(r.Context())
		if err != nil {
			httpresponse.RespondWithBadRequestErrorNoErr(w, err.Error())
			return
		}

		screenshotID, err := getUUIDFromRequest(env, r, "screenshot-id")
		if err != nil {
			httpresponse.RespondWithBadRequestErrorNoErr(w, err.Error())
			return
		}

		if env.ScreenshotService == nil {
			httpresponse.RespondWithInternalServerError(w, errors.New("No Screenshot Service provided"))
			return
		}

		screenshot, err := env.ScreenshotService.Find(&models.Screenshot{Record: models.Record{ID: screenshotID}, AppVersionID: appVersionID})

		switch {
		case errors.Cause(err) == gorm.ErrRecordNotFound:
			httpresponse.RespondWithNotFoundErrorNoErr(w)
			return
		case err != nil:
			httpresponse.RespondWithInternalServerError(w, errors.WithStack(err))
			return
		}

		// Access granted
		ctx := ContextWithAuthorizedScreenshotID(r.Context(), screenshot.ID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUUIDFromRequest(env *env.AppEnv, r *http.Request, paramName string) (uuid.UUID, error) {
	urlVars := env.RequestParams.Get(r)
	param := urlVars[paramName]
	if param == "" {
		return uuid.UUID{}, errors.Errorf("Failed to fetch URL param %s", paramName)
	}

	paramUUID, err := uuid.FromString(param)
	if err != nil {
		return uuid.UUID{}, errors.Errorf("Invalid UUID format for %s", paramName)
	}

	return paramUUID, nil
}

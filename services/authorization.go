package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httprequest"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/api-utils/security"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// AuthorizeForAppDeprovisioningHandlerFunc ...
func AuthorizeForAppDeprovisioningHandlerFunc(env *env.AppEnv, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		app, err := env.AppService.Find(&models.App{AppSlug: appSlug})
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

// AuthorizeForAddonAPIAccessHandlerFunc ...
func AuthorizeForAddonAPIAccessHandlerFunc(env *env.AppEnv, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken, err := httprequest.AuthTokenFromHeader(r.Header)
		if err != nil {
			httpresponse.RespondWithUnauthorizedNoErr(w)
			return
		}

		if env.AppService == nil {
			httpresponse.RespondWithInternalServerError(w, errors.New("No App Service provided"))
			return
		}

		app, err := env.AppService.Find(&models.App{APIToken: authToken})
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

		authorizedAppID, err := GetAuthorizedAppIDFromContext(r.Context())
		if err != nil {
			httpresponse.RespondWithInternalServerError(w, err)
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
		appVersion, err := env.AppVersionService.Find(&models.AppVersion{Record: models.Record{ID: appVersionID}, AppID: authorizedAppID})
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
			httpresponse.RespondWithInternalServerError(w, err)
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

// AuthorizeForWebhookHandlerFunc ...
func AuthorizeForWebhookHandlerFunc(env *env.AppEnv, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload WebhookPayload
		defer httprequest.BodyCloseWithErrorLog(r)
		payloadBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			httpresponse.RespondWithInternalServerError(w, errors.Wrap(err, "Failed to read request body"))
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewReader(payloadBytes))
		err = json.Unmarshal(payloadBytes, &payload)
		if err != nil {
			httpresponse.RespondWithBadRequestErrorNoErr(w, "Invalid request body, JSON decode failed")
			return
		}

		if env.PublishTaskService == nil {
			httpresponse.RespondWithInternalServerError(w, errors.New("No Publish Task Service provided"))
			return
		}

		publishTask, err := env.PublishTaskService.Find(&models.PublishTask{TaskID: payload.TaskID})
		switch {
		case errors.Cause(err) == gorm.ErrRecordNotFound:
			httpresponse.RespondWithNotFoundErrorNoErr(w)
			return
		case err != nil:
			httpresponse.RespondWithInternalServerError(w, errors.WithStack(err))
			return
		}

		// Access granted
		ctx := ContextWithAuthorizedAppVersionID(r.Context(), publishTask.AppVersionID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthorizeForAppContactEmailConfirmationHandlerFunc ...
func AuthorizeForAppContactEmailConfirmationHandlerFunc(env *env.AppEnv, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBody := struct {
			ConfirmationToken string `json:"confirmation_token"`
		}{}
		defer httprequest.BodyCloseWithErrorLog(r)
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			httpresponse.RespondWithBadRequestErrorNoErr(w, "Invalid request body, JSON decode failed")
			return
		}

		if env.AppContactService == nil {
			httpresponse.RespondWithInternalServerError(w, errors.New("No App Contact Service provided"))
			return
		}

		appContact, err := env.AppContactService.Find(&models.AppContact{ConfirmationToken: &requestBody.ConfirmationToken})
		switch {
		case errors.Cause(err) == gorm.ErrRecordNotFound:
			httpresponse.RespondWithNotFoundErrorNoErr(w)
			return
		case err != nil:
			httpresponse.RespondWithInternalServerError(w, errors.WithStack(err))
			return
		}

		// Access granted
		ctx := ContextWithAuthorizedAppContactID(r.Context(), appContact.ID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthorizeForAppContactAccessHandlerFunc ...
func AuthorizeForAppContactAccessHandlerFunc(env *env.AppEnv, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if env.RequestParams == nil {
			httpresponse.RespondWithInternalServerError(w, errors.New("No Request Params provided"))
			return
		}

		appID, err := GetAuthorizedAppIDFromContext(r.Context())
		if err != nil {
			httpresponse.RespondWithInternalServerError(w, err)
			return
		}

		appContactID, err := getUUIDFromRequest(env, r, "contact-id")
		if err != nil {
			httpresponse.RespondWithBadRequestErrorNoErr(w, err.Error())
			return
		}

		if env.AppContactService == nil {
			httpresponse.RespondWithInternalServerError(w, errors.New("No App Contact Service provided"))
			return
		}

		appContact, err := env.AppContactService.Find(&models.AppContact{Record: models.Record{ID: appContactID}, AppID: appID})
		switch {
		case errors.Cause(err) == gorm.ErrRecordNotFound:
			httpresponse.RespondWithNotFoundErrorNoErr(w)
			return
		case err != nil:
			httpresponse.RespondWithInternalServerError(w, errors.WithStack(err))
			return
		}

		// Access granted
		ctx := ContextWithAuthorizedAppContactID(r.Context(), appContact.ID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthorizeBuildWebhookForAppAccessFunc ...
func AuthorizeBuildWebhookForAppAccessFunc(env *env.AppEnv, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appSlug := r.Header.Get("Bitrise-App-Id")
		if appSlug == "" {
			httpresponse.RespondWithNotFoundErrorNoErr(w)
			return
		}

		if env.AppService == nil {
			httpresponse.RespondWithInternalServerError(w, errors.New("No App Service provided"))
			return
		}

		app, err := env.AppService.Find(&models.App{AppSlug: appSlug})
		switch {
		case errors.Cause(err) == gorm.ErrRecordNotFound:
			httpresponse.RespondWithNotFoundErrorNoErr(w)
			return
		case err != nil:
			httpresponse.RespondWithInternalServerError(w, err)
			return
		}

		if len(app.EncryptedSecretIV) == 0 {
			ctx := ContextWithAuthorizedAppID(r.Context(), app.ID)
			h.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		appSecret, err := app.Secret()
		if err != nil {
			httpresponse.RespondWithInternalServerError(w, err)
			return
		}
		if appSecret == "" {
			ctx := ContextWithAuthorizedAppID(r.Context(), app.ID)
			h.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		requestPayloadSignature := r.Header.Get("Bitrise-Hook-Signature")
		payloadBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			httpresponse.RespondWithInternalServerError(w, errors.Wrap(err, "Failed to get request payload"))
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewReader(payloadBytes))

		signatureVerifier := security.NewSignatureVerifier(appSecret, string(payloadBytes), requestPayloadSignature)
		if !signatureVerifier.Verify() {
			httpresponse.RespondWithNotFoundErrorNoErr(w)
			return
		}

		// Access granted
		ctx := ContextWithAuthorizedAppID(r.Context(), app.ID)
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

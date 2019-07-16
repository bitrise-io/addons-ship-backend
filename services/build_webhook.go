package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httprequest"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// BuildWebhookPayload ...
type BuildWebhookPayload struct {
	AppSlug                string `json:"app_slug"`
	BuildSlug              string `json:"build_slug"`
	BuildNumber            int    `json:"build_number"`
	BuildStatus            string `json:"build_status"`
	BuildTriggeredWorkflow string `json:"build_triggered_workflow"`
}

// BuildWebhookHandler ...
func BuildWebhookHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppID, err := GetAuthorizedAppIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	switch r.Header.Get("Bitrise-Event-Type") {
	case "build/started":
		return httpresponse.RespondWithSuccess(w, nil)
	case "build/finished":
		if env.AppService == nil {
			return errors.New("No App Service defined for handler")
		}
		if env.AppSettingsService == nil {
			return errors.New("No App Settings Service defined for handler")
		}
		if env.AppVersionService == nil {
			return errors.New("No App Version Service defined for handler")
		}
		var params BuildWebhookPayload
		defer httprequest.BodyCloseWithErrorLog(r)
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
		}

		_, err := env.AppService.Find(&models.App{Record: models.Record{ID: authorizedAppID}})
		if err != nil {
			return errors.Wrap(err, "SQL Error")
		}

		appSettings, err := env.AppSettingsService.Find(&models.AppSettings{AppID: authorizedAppID})
		switch {
		case errors.Cause(err) == gorm.ErrRecordNotFound:
			return httpresponse.RespondWithNotFoundError(w)
		case err != nil:
			return errors.Wrap(err, "SQL Error")
		}

		if appSettings.IosWorkflow == "all" ||
			(params.BuildTriggeredWorkflow != "" && strings.Contains(appSettings.IosWorkflow, params.BuildTriggeredWorkflow)) {
			_, verrs, err := env.AppVersionService.Create(&models.AppVersion{
				Platform:    "ios",
				BuildNumber: fmt.Sprintf("%d", params.BuildNumber),
				BuildSlug:   params.BuildSlug,
				LastUpdate:  time.Now(),
				AppID:       authorizedAppID,
			})
			if len(verrs) > 0 {
				return httpresponse.RespondWithUnprocessableEntity(w, verrs)
			}
			if err != nil {
				return errors.Wrap(err, "SQL Error")
			}
		}

		if appSettings.AndroidWorkflow == "all" ||
			(params.BuildTriggeredWorkflow != "" && strings.Contains(appSettings.AndroidWorkflow, params.BuildTriggeredWorkflow)) {
			_, verrs, err := env.AppVersionService.Create(&models.AppVersion{
				Platform:    "android",
				BuildNumber: fmt.Sprintf("%d", params.BuildNumber),
				BuildSlug:   params.BuildSlug,
				LastUpdate:  time.Now(),
				AppID:       authorizedAppID,
			})
			if len(verrs) > 0 {
				return httpresponse.RespondWithUnprocessableEntity(w, verrs)
			}
			if err != nil {
				return errors.Wrap(err, "SQL Error")
			}
		}

		return httpresponse.RespondWithSuccess(w, nil)
	default:
		return errors.New("Invalid build event")
	}
}

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
	BuildStatus            int    `json:"build_status"`
	BuildTriggeredWorkflow string `json:"build_triggered_workflow"`
}

// BuildWebhookHandler ...
func BuildWebhookHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppID, err := GetAuthorizedAppIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	switch r.Header.Get("Bitrise-Event-Type") {
	case "build/triggered":
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
		if env.BitriseAPI == nil {
			return errors.New("No Bitrise API Service defined for handler")
		}
		if env.AppContactService == nil {
			return errors.New("No App Contact Service defined for handler")
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

		app := appSettings.App
		if appSettings.IosWorkflow == "all" ||
			(params.BuildTriggeredWorkflow != "" && strings.Contains(appSettings.IosWorkflow, params.BuildTriggeredWorkflow)) {
			latestAppVersion, err := env.AppVersionService.Latest(&models.AppVersion{AppID: app.ID, Platform: "android"})
			if err != nil && errors.Cause(err) != gorm.ErrRecordNotFound {
				return errors.Wrap(err, "SQL Error")
			}
			appVersion, err := prepareAppVersionForIosPlatform(env, w, r, app.BitriseAPIToken, app.AppSlug, params.BuildSlug)
			if err != nil {
				return err
			}
			appVersion.LastUpdate = time.Now()
			appVersion.AppID = authorizedAppID
			appVersion.BuildNumber = fmt.Sprintf("%d", params.BuildNumber)
			if latestAppVersion != nil {
				appVersion.AppStoreInfoData = latestAppVersion.AppStoreInfoData
			}
			appVersion, verrs, err := env.AppVersionService.Create(appVersion)
			if len(verrs) > 0 {
				return httpresponse.RespondWithUnprocessableEntity(w, verrs)
			}
			if err != nil {
				return errors.Wrap(err, "SQL Error")
			}
			if latestAppVersion != nil {
				err := env.WorkerService.EnqueueCopyUploadablesToNewAppVersion(latestAppVersion.ID.String(), appVersion.ID.String())
				if err != nil {
					return errors.Wrap(err, "Worker Error")
				}
			}

			if err := sendNotification(env, appVersion, app); err != nil {
				return errors.WithStack(err)
			}
		}

		if appSettings.AndroidWorkflow == "all" ||
			(params.BuildTriggeredWorkflow != "" && strings.Contains(appSettings.AndroidWorkflow, params.BuildTriggeredWorkflow)) {
			latestAppVersion, err := env.AppVersionService.Latest(&models.AppVersion{AppID: app.ID, Platform: "android"})
			if err != nil && errors.Cause(err) != gorm.ErrRecordNotFound {
				return errors.Wrap(err, "SQL Error")
			}
			appVersion, err := prepareAppVersionForAndroidPlatform(env, w, r, app.BitriseAPIToken, app.AppSlug, params.BuildSlug)
			if err != nil {
				return err
			}
			appVersion.LastUpdate = time.Now()
			appVersion.AppID = authorizedAppID
			appVersion.BuildNumber = fmt.Sprintf("%d", params.BuildNumber)
			if latestAppVersion != nil {
				appVersion.AppStoreInfoData = latestAppVersion.AppStoreInfoData
			}
			appVersion, verrs, err := env.AppVersionService.Create(appVersion)
			if len(verrs) > 0 {
				return httpresponse.RespondWithUnprocessableEntity(w, verrs)
			}
			if err != nil {
				return errors.Wrap(err, "SQL Error")
			}
			if latestAppVersion != nil {
				err := env.WorkerService.EnqueueCopyUploadablesToNewAppVersion(latestAppVersion.ID.String(), appVersion.ID.String())
				if err != nil {
					return errors.Wrap(err, "Worker Error")
				}
			}

			if err := sendNotification(env, appVersion, app); err != nil {
				return errors.WithStack(err)
			}
		}

		return httpresponse.RespondWithSuccess(w, nil)
	default:
		return errors.New("Invalid build event")
	}
}

func sendNotification(env *env.AppEnv, appVersion *models.AppVersion, app *models.App) error {
	appContacts, err := env.AppContactService.FindAll(app)
	if err != nil {
		return errors.WithStack(err)
	}
	appDetails, err := env.BitriseAPI.GetAppDetails(app.BitriseAPIToken, app.AppSlug)
	if err != nil {
		return errors.WithStack(err)
	}
	return env.Mailer.SendEmailNewVersion(appVersion, appContacts, env.AddonFrontendHostURL, appDetails)
}

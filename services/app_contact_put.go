package services

import (
	"encoding/json"
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httprequest"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

// AppContactPutHandler ...
func AppContactPutHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppContactID, err := GetAuthorizedAppContactIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	if env.AppContactService == nil {
		return errors.New("No App Contact Service defined for handler")
	}

	var params models.NotificationPreferences
	defer httprequest.BodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}

	appContact, err := env.AppContactService.Find(&models.AppContact{Record: models.Record{ID: authorizedAppContactID}})
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	notificationPreferences, err := json.Marshal(params)
	appContact.NotificationPreferencesData = notificationPreferences
	err = env.AppContactService.Update(appContact, []string{"NotificationPreferencesData"})
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	return nil
}

package services

import (
	"net/http"
	"time"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

// AppContactPatchResponseData ...
type AppContactPatchResponseData struct {
	AppContact *models.AppContact `json:"app_contact"`
	App        *models.App        `json:"app"`
}

// AppContactPatchResponse ...
type AppContactPatchResponse struct {
	Data AppContactPatchResponseData `json:"data"`
}

// AppContactConfirmPatchHandler ...
func AppContactConfirmPatchHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppContactID, err := GetAuthorizedAppContactIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if env.AppContactService == nil {
		return errors.New("No App Contact Service defined for handler")
	}
	appContact, err := env.AppContactService.Find(&models.AppContact{Record: models.Record{ID: authorizedAppContactID}})
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}
	appContact.ConfirmedAt = time.Now()
	appContact.ConfirmationToken = nil
	err = env.AppContactService.Update(appContact, []string{"ConfirmedAt", "ConfirmationToken"})
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	return httpresponse.RespondWithSuccess(w, AppContactPatchResponse{
		Data: AppContactPatchResponseData{AppContact: appContact, App: appContact.App},
	})
}

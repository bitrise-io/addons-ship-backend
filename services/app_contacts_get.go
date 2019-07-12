package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

// AppContactsGetResponse ...
type AppContactsGetResponse struct {
	Data []models.AppContact `json:"data"`
}

// AppContactsGetHandler ...
func AppContactsGetHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppID, err := GetAuthorizedAppIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	if env.AppContactService == nil {
		return errors.New("No App Contact Service defined for handler")
	}

	appContacts, err := env.AppContactService.FindAll(&models.App{Record: models.Record{ID: authorizedAppID}})
	if err != nil {
		return errors.WithStack(err)
	}

	return httpresponse.RespondWithSuccess(w, AppContactsGetResponse{Data: appContacts})
}

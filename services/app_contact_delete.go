package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

// AppContactDeleteResponse ...
type AppContactDeleteResponse struct {
	Data *models.AppContact `json:"data"`
}

// AppContactDeleteHandler ...
func AppContactDeleteHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
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

	err = env.AppContactService.Delete(appContact)
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	return httpresponse.RespondWithSuccess(w, AppContactDeleteResponse{Data: appContact})
}

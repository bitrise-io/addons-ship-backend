package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// AppGetHandler ...
func AppGetHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppID, err := GetAuthorizedAppIDFromContextErr(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if env.AppService == nil {
		return errors.New("No App Service defined for handler")
	}

	app, err := env.AppService.Find(&models.App{Record: models.Record{ID: authorizedAppID}})
	switch {
	case errors.Cause(err) == gorm.ErrRecordNotFound:
		return httpresponse.RespondWithNotFoundError(w)
	case err != nil:
		return errors.WithStack(err)
	}
	return httpresponse.RespondWithSuccess(w, app)
}

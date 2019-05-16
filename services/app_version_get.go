package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// AppVersionGetResponse ...
type AppVersionGetResponse struct {
	Data *models.AppVersion `json:"data"`
}

// AppVersionGetHandler ...
func AppVersionGetHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppVersionID, err := GetAuthorizedAppVersionIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if env.AppVersionService == nil {
		return errors.New("No App Version Service defined for handler")
	}

	appVersion, err := env.AppVersionService.Find(
		&models.AppVersion{Record: models.Record{ID: authorizedAppVersionID}},
	)
	switch {
	case errors.Cause(err) == gorm.ErrRecordNotFound:
		return httpresponse.RespondWithNotFoundError(w)
	case err != nil:
		return errors.Wrap(err, "SQL Error")
	}
	return httpresponse.RespondWithSuccess(w, AppVersionGetResponse{
		Data: appVersion,
	})
}

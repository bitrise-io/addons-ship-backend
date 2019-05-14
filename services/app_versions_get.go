package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// AppVersionsGetRespose ...
type AppVersionsGetRespose struct {
	Data []models.AppVersion `json:"data"`
}

// AppVersionsGetHandler ...
func AppVersionsGetHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppID, err := GetAuthorizedAppIDFromContextErr(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if env.AppVersionService == nil {
		return errors.New("No App Version Service defined for handler")
	}

	filterParams := map[string]interface{}{}
	if platformFilter := r.URL.Query().Get("platform"); platformFilter != "" {
		filterParams["platform"] = platformFilter
	}

	appVersions, err := env.AppVersionService.FindAll(
		&models.App{Record: models.Record{ID: authorizedAppID}},
		filterParams,
	)
	switch {
	case errors.Cause(err) == gorm.ErrRecordNotFound:
		return httpresponse.RespondWithNotFoundError(w)
	case err != nil:
		return errors.Wrap(err, "SQL Error")
	}
	return httpresponse.RespondWithSuccess(w, AppVersionsGetRespose{
		Data: appVersions,
	})
}

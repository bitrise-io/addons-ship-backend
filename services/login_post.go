package services

import (
	"fmt"
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/pkg/errors"
)

// LoginPostHandler ...
func LoginPostHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppID, err := GetAuthorizedAppIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if env.AppService == nil {
		return errors.New("No App Service defined for handler")
	}

	app, err := env.AppService.Find(&models.App{Record: models.Record{ID: authorizedAppID}})
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	redirectURL := fmt.Sprintf("%s/apps/%s", env.AddonFrontendHostURL, app.AppSlug)
	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
	return nil
}

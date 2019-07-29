package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httprequest"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/go-crypto/crypto"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// ProvisionPostParams ...
type ProvisionPostParams struct {
	AppSlug         string `json:"app_slug"`
	BitriseAPIToken string `json:"bitrise_api_token"`
	Plan            string `json:"plan"`
}

// Env ...
type Env struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ProvisionPostResponse ...
type ProvisionPostResponse struct {
	Envs []Env `json:"envs"`
}

// ProvisionHandler ...
func ProvisionHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	if env.AppService == nil {
		return errors.New("No App Service defined for handler")
	}
	var params ProvisionPostParams
	defer httprequest.BodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}

	if env.BitriseAPI == nil {
		return errors.New("No Bitrise API Service defined for handler")
	}

	app, err := env.AppService.Find(&models.App{AppSlug: params.AppSlug, BitriseAPIToken: params.BitriseAPIToken})
	switch {
	case errors.Cause(err) == gorm.ErrRecordNotFound:
		app, err = env.AppService.Create(&models.App{
			AppSlug:         params.AppSlug,
			BitriseAPIToken: params.BitriseAPIToken,
			Plan:            params.Plan,
			APIToken:        crypto.SecureRandomHash(50),
		})
		if err != nil {
			return errors.Wrap(err, "SQL Error")
		}
		secret, err := app.Secret()
		if err != nil {
			return errors.WithStack(err)
		}
		err = env.BitriseAPI.RegisterWebhook(app.BitriseAPIToken, app.AppSlug, secret, fmt.Sprintf("%s/webhook", env.AddonHostURL))
		if err != nil {
			return errors.WithStack(err)
		}
	case err != nil:
		return errors.Wrap(err, "SQL Error")
	}

	envs := []Env{
		Env{Key: "ADDON_SHIP_API_URL", Value: env.AddonHostURL},
		Env{Key: "ADDON_SHIP_API_TOKEN", Value: app.APIToken},
	}
	return httpresponse.RespondWithSuccess(w, ProvisionPostResponse{Envs: envs})
}

package services

import (
	"encoding/json"
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httprequest"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/go-crypto/crypto"
	"github.com/pkg/errors"
)

type appContactPostParams struct {
	Email       string                         `json:"email"`
	Preferences models.NotificationPreferences `json:"notification_preferences"`
}

// AppContactPostResponse ...
type AppContactPostResponse struct {
	Data *models.AppContact `json:"data"`
}

// AppContactPostHandler ...
func AppContactPostHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppID, err := GetAuthorizedAppIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	if env.AppContactService == nil {
		return errors.New("No App Contact Service defined for handler")
	}
	if env.BitriseAPI == nil {
		return errors.New("No Bitrise API Service defined for handler")
	}
	if env.Mailer == nil {
		return errors.New("No Mailer defined for handler")
	}

	var params appContactPostParams
	defer httprequest.BodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}
	appContact, err := env.AppContactService.Find(&models.AppContact{
		AppID: authorizedAppID,
		Email: params.Email,
	})
	if err == nil {
		return httpresponse.RespondWithSuccess(w, AppContactPostResponse{Data: appContact})
	}
	notificationPreferences, err := json.Marshal(params.Preferences)
	if err != nil {
		return errors.WithStack(err)
	}
	confirmationToken := crypto.SecureRandomHash(24)
	appContact, verrs, err := env.AppContactService.Create(&models.AppContact{
		AppID: authorizedAppID,
		Email: params.Email,
		NotificationPreferencesData: notificationPreferences,
		ConfirmationToken:           &confirmationToken,
	})
	if len(verrs) > 0 {
		return httpresponse.RespondWithUnprocessableEntity(w, verrs)
	}
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}
	appDetails, err := env.BitriseAPI.GetAppDetails(appContact.App.BitriseAPIToken, appContact.App.AppSlug)
	if err != nil {
		return errors.WithStack(err)
	}

	err = env.Mailer.SendEmailConfirmation(env.EmailConfirmLandingURL, appContact, appDetails)
	if err != nil {
		return errors.WithStack(err)
	}

	return httpresponse.RespondWithSuccess(w, AppContactPostResponse{Data: appContact})
}

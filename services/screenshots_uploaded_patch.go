package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// ScreenshotsUploadedPatchResponse ...
type ScreenshotsUploadedPatchResponse struct {
	Data []ScreenshotData `json:"data"`
}

// ScreenshotsUploadedPatchHandler ...
func ScreenshotsUploadedPatchHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppVersionID, err := GetAuthorizedAppVersionIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	if env.ScreenshotService == nil {
		return errors.New("No Screenshot Service defined for handler")
	}

	screenshotsToUpdate, err := prepareScreenshotsToUpdate(env.ScreenshotService, authorizedAppVersionID)
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	verrs, err := env.ScreenshotService.BatchUpdate(screenshotsToUpdate, []string{"Uploaded"})
	if len(verrs) > 0 {
		return httpresponse.RespondWithUnprocessableEntity(w, verrs)
	}
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	if env.AWS == nil {
		return errors.New("No AWS Provider defined for handler")
	}
	responseData, err := newScreenshotGetResponseData(screenshotsToUpdate, env.AWS)
	if err != nil {
		return errors.WithStack(err)
	}
	return httpresponse.RespondWithSuccess(w, ScreenshotsUploadedPatchResponse{
		Data: responseData,
	})
}

func prepareScreenshotsToUpdate(screenshotService dataservices.ScreenshotService, appVersionID uuid.UUID) ([]models.Screenshot, error) {
	var screenshotsToUpdate []models.Screenshot

	screenshots, err := screenshotService.FindAll(&models.AppVersion{Record: models.Record{ID: appVersionID}})
	if err != nil {
		return []models.Screenshot{}, err
	}
	for _, screenshot := range screenshots {
		screenshot.Uploaded = true
		screenshotsToUpdate = append(screenshotsToUpdate, screenshot)
	}
	return screenshotsToUpdate, nil
}

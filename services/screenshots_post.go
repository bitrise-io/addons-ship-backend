package services

import (
	"encoding/json"
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httprequest"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/api-utils/providers"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type screenshotsPostParamsElement struct {
	Filename   string `json:"filename"`
	Filesize   int64  `json:"filesize"`
	DeviceType string `json:"device_type"`
	ScreenSize string `json:"screen_size"`
}

type screenshotsPostParams struct {
	Screenshots []screenshotsPostParamsElement `json:"screenshots"`
}

// ScreenshotsPostResponse ...
type ScreenshotsPostResponse struct {
	Data []ScreenshotData `json:"data"`
}

// ScreenshotsPostHandler ...
func ScreenshotsPostHandler(env *env.AppEnv, w http.ResponseWriter, r *http.Request) error {
	authorizedAppVersionID, err := GetAuthorizedAppVersionIDFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	var params screenshotsPostParams
	defer httprequest.BodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}

	if env.ScreenshotService == nil {
		return errors.New("No Screenshot Service defined for handler")
	}

	createdScreenshots, verrs, err := env.ScreenshotService.BatchCreate(screenshotCreateParamsFromRequestParams(params.Screenshots, authorizedAppVersionID))
	if len(verrs) > 0 {
		return httpresponse.RespondWithUnprocessableEntity(w, verrs)
	}
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}

	if env.AWS == nil {
		return errors.New("No AWS Provider defined for handler")
	}
	responseData, err := newScreenshotPostResponseData(createdScreenshots, env.AWS)
	if err != nil {
		return errors.WithStack(err)
	}

	return httpresponse.RespondWithSuccess(w, ScreenshotsPostResponse{
		Data: responseData,
	})
}

func screenshotCreateParamsFromRequestParams(params []screenshotsPostParamsElement, appVersionID uuid.UUID) []*models.Screenshot {
	var createParams []*models.Screenshot
	for _, param := range params {
		createParams = append(createParams, &models.Screenshot{
			AppVersionID: appVersionID,
			Uploadable: models.Uploadable{
				Filename: param.Filename,
				Filesize: param.Filesize,
			},
			DeviceType: param.DeviceType,
			ScreenSize: param.ScreenSize,
		})
	}
	return createParams
}

func newScreenshotPostResponseData(screenshots []*models.Screenshot, awsProvider providers.AWSInterface) ([]ScreenshotData, error) {
	data := []ScreenshotData{}
	for _, screenshot := range screenshots {
		presignedURL, err := awsProvider.GeneratePresignedPUTURL(screenshot.AWSPath(), presignedURLExpirationInterval, screenshot.Filesize)
		if err != nil {
			return []ScreenshotData{}, errors.WithStack(err)
		}
		data = append(data, ScreenshotData{Screenshot: *screenshot, UploadURL: presignedURL})
	}
	return data, nil
}

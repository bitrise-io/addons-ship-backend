package worker

import (
	"net/url"

	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/gocraft/work"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

var copyUploadablesToNewAppVersion = "copy_uploadables_to_new_app_version"

// CopyUploadablesToNewAppVersion ...
func (c *Context) CopyUploadablesToNewAppVersion(job *work.Job) error {
	c.env.Logger.Info("[i] Job CopyUploadablesToNewAppVersion started")
	appVersionFromID := job.ArgString("from_id")
	if appVersionFromID == "" {
		c.env.Logger.Error("Failed to get ID of app version to copy screenshots from")
		return errors.New("Failed to get from_id")
	}
	appVersionToID := job.ArgString("to_id")
	if appVersionToID == "" {
		c.env.Logger.Error("Failed to get ID of app version to copy screenshots to")
		return errors.New("Failed to get to_id")
	}

	c.env.Logger.Info("[i] CopyUploadablesToNewAppVersion: Copying screenshots...")
	originalScreenshots, err := c.env.ScreenshotService.FindAll(&models.AppVersion{Record: models.Record{ID: uuid.FromStringOrNil(appVersionFromID)}})
	if err != nil {
		return errors.Wrap(err, "SQL Error")
	}
	newAppVersionID := uuid.FromStringOrNil(appVersionToID)
	if len(originalScreenshots) > 0 {
		screenShotsToCreate := []*models.Screenshot{}
		for _, sc := range originalScreenshots {
			screenShotsToCreate = append(screenShotsToCreate, &models.Screenshot{
				UploadableObject: sc.UploadableObject,
				DeviceType:       sc.DeviceType,
				ScreenSize:       sc.ScreenSize,
				AppVersionID:     newAppVersionID,
			})
		}
		createsScreenshots, verrs, err := c.env.ScreenshotService.BatchCreate(screenShotsToCreate)
		if err != nil {
			return errors.Wrap(err, "SQL Error")
		}
		if len(verrs) > 0 {
			return errors.Errorf("Validation errors: %#v", verrs)
		}

		c.env.Logger.Info("[i] CopyUploadablesToNewAppVersion: Copying S3 files...")

		for idx, sc := range originalScreenshots {
			from := url.QueryEscape(sc.AWSPath())
			to := createsScreenshots[idx].AWSPath()

			err = c.env.AWS.CopyObject(from, to)
			if err != nil {
				c.env.Logger.Error("[!] CopyUploadablesToNewAppVersion: Failed to copy AWS file", zap.Any("error", err))
				return errors.WithStack(err)
			}
		}
	}

	c.env.Logger.Info("[i] CopyUploadablesToNewAppVersion: Copying feature graphic...")
	originalFeatureGraphic, err := c.env.FeatureGraphicService.Find(&models.FeatureGraphic{AppVersionID: uuid.FromStringOrNil(appVersionFromID)})
	if err != nil && errors.Cause(err) != gorm.ErrRecordNotFound {
		return errors.Wrap(err, "SQL Error")
	}

	if originalFeatureGraphic != nil {
		createsFeatureGraphic, verrs, err := c.env.FeatureGraphicService.Create(&models.FeatureGraphic{
			UploadableObject: originalFeatureGraphic.UploadableObject,
			AppVersionID:     newAppVersionID,
		})
		if err != nil {
			return errors.Wrap(err, "SQL Error")
		}
		if len(verrs) > 0 {
			return errors.Errorf("Validation error: %#v", verrs)
		}

		from := url.QueryEscape(originalFeatureGraphic.AWSPath())
		to := createsFeatureGraphic.AWSPath()
		err = c.env.AWS.CopyObject(from, to)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	c.env.Logger.Info("[i] Job CopyUploadablesToNewAppVersion finished")
	return nil
}

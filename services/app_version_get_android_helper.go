package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

func appVersionGetAndroidHelper(env *env.AppEnv, w http.ResponseWriter,
	r *http.Request, appVersion *models.AppVersion,
	artifacts []bitrise.ArtifactListElementResponseModel) error {
	publishEnabled := false
	publicInstallPageEnabled := false
	publicInstallPageArtifactSlug := ""
	var selectedArtifact bitrise.ArtifactListElementResponseModel

	for _, artifact := range artifacts {
		if artifact.IsAAB() {
			publishEnabled = true
			selectedArtifact = artifact
		}
		if artifact.IsUniversalAPK() {
			publicInstallPageEnabled = true
			publicInstallPageArtifactSlug = artifact.Slug
			selectedArtifact = artifact
		}
		// TODO: check the split APK condition
	}
	var artifactPublicInstallPageURL string
	if publicInstallPageEnabled {
		var err error
		artifactPublicInstallPageURL, err = env.BitriseAPI.GetArtifactPublicInstallPageURL(
			appVersion.App.BitriseAPIToken,
			appVersion.App.AppSlug,
			appVersion.BuildSlug,
			publicInstallPageArtifactSlug,
		)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	appDetails, err := env.BitriseAPI.GetAppDetails(appVersion.App.BitriseAPIToken, appVersion.App.AppSlug)
	if err != nil {
		return errors.WithStack(err)
	}

	responseData, err := newArtifactVersionGetResponse(appVersion, selectedArtifact.ArtifactMeta, artifactPublicInstallPageURL, appDetails)
	if err != nil {
		return errors.WithStack(err)
	}
	responseData.PublishEnabled = publishEnabled

	return httpresponse.RespondWithSuccess(w, AppVersionGetResponse{
		Data: responseData,
	})
}

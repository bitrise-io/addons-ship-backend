package services

import (
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

func appVersionGetIosHelper(env *env.AppEnv, w http.ResponseWriter,
	r *http.Request, appVersion *models.AppVersion,
	artifacts []bitrise.ArtifactListElementResponseModel) error {
	publishEnabled := false
	publicInstallPageEnabled := false
	publicInstallPageArtifactSlug := ""
	var selectedArtifact bitrise.ArtifactListElementResponseModel

	for _, artifact := range artifacts {
		if artifact.IsIPA() {
			if artifact.HasAppStoreDistributionType() {
				publishEnabled = true
				selectedArtifact = artifact
			}
			if artifact.HasDebugDistributionType() {
				publicInstallPageEnabled = true
				publicInstallPageArtifactSlug = artifact.Slug
				if selectedArtifact == (bitrise.ArtifactListElementResponseModel{}) {
					selectedArtifact = artifact
				}
			}
		}
		if artifact.IsXCodeArchive() {
			publishEnabled = true
			if selectedArtifact == (bitrise.ArtifactListElementResponseModel{}) {
				selectedArtifact = artifact
			}
		}
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

	responseData, err := newArtifactVersionGetResponse(appVersion, selectedArtifact, artifactPublicInstallPageURL, appDetails, publishEnabled)
	if err != nil {
		return errors.WithStack(err)
	}

	return httpresponse.RespondWithSuccess(w, AppVersionGetResponse{
		Data: responseData,
	})
}

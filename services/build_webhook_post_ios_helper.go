package services

import (
	"encoding/json"
	"net/http"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

func prepareAppVersionForIosPlatform(env *env.AppEnv, w http.ResponseWriter,
	r *http.Request, apiToken, appSlug, buildSlug string) (*models.AppVersion, error) {
	artifacts, err := env.BitriseAPI.GetArtifacts(apiToken, appSlug, buildSlug)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	publishEnabled := false
	publicInstallPageEnabled := false
	publicInstallPageArtifactSlug := ""
	var selectedArtifact *bitrise.ArtifactListElementResponseModel

	for _, artifact := range artifacts {
		if artifact.IsIPA() {
			if artifact.HasAppStoreDistributionType() {
				publishEnabled = true
				selectedArtifact = &artifact
			}
			if artifact.HasDebugDistributionType() {
				publicInstallPageEnabled = true
				publicInstallPageArtifactSlug = artifact.Slug
				if selectedArtifact == nil {
					selectedArtifact = &artifact
				}
			}
		}
		if artifact.IsXCodeArchive() {
			publishEnabled = true
			if selectedArtifact == nil {
				selectedArtifact = &artifact
			}
		}
	}
	var artifactPublicInstallPageURL string
	if publicInstallPageEnabled {
		var err error
		artifactPublicInstallPageURL, err = env.BitriseAPI.GetArtifactPublicInstallPageURL(
			apiToken,
			appSlug,
			buildSlug,
			publicInstallPageArtifactSlug,
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	if selectedArtifact == nil {
		return nil, httpresponse.RespondWithNotFoundError(w)
	}

	appDetails, err := env.BitriseAPI.GetAppDetails(apiToken, appSlug)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var supportedDeviceTypes []string
	for _, familyID := range selectedArtifact.ArtifactMeta.AppInfo.DeviceFamilyList {
		switch familyID {
		case 1:
			supportedDeviceTypes = append(supportedDeviceTypes, "iPhone", "iPod Touch")
		case 2:
			supportedDeviceTypes = append(supportedDeviceTypes, "iPad")
		default:
			supportedDeviceTypes = append(supportedDeviceTypes, "Unknown")
		}
	}
	appInfo := models.AppInfo{
		MinimumOS:            selectedArtifact.ArtifactMeta.AppInfo.MinimumOS,
		BundleID:             selectedArtifact.ArtifactMeta.AppInfo.BundleID,
		SupportedDeviceTypes: supportedDeviceTypes,
	}
	appInfoData, err := json.Marshal(appInfo)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	provisioningInfo := models.ProvisioningInfo{
		ExpireDate:       selectedArtifact.ArtifactMeta.ProvisioningInfo.ExpireDate,
		DistributionType: selectedArtifact.ArtifactMeta.ProvisioningInfo.DistributionType,
	}
	provisioningInfoData, err := json.Marshal(provisioningInfo)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &models.AppVersion{
		Platform:             "ios",
		Version:              selectedArtifact.ArtifactMeta.AppInfo.Version,
		BuildSlug:            buildSlug,
		AppInfoData:          appInfoData,
		ProvisioningInfoData: provisioningInfoData,
	}, nil
}

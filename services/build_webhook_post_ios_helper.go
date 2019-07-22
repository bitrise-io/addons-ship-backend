package services

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/pkg/errors"
)

func prepareAppVersionForIosPlatform(env *env.AppEnv, w http.ResponseWriter,
	r *http.Request, apiToken, appSlug, buildSlug string) (*models.AppVersion, error) {
	artifacts, err := env.BitriseAPI.GetArtifacts(apiToken, appSlug, buildSlug)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	selectedArtifact, _, _, _ := selectIosArtifact(artifacts)
	if selectedArtifact == nil {
		return nil, errors.New("No artifact found")
	}

	if selectedArtifact.ArtifactMeta == nil {
		return nil, errors.New("No artifact meta data found for artifact")
	}

	if reflect.DeepEqual(selectedArtifact.ArtifactMeta.AppInfo, bitrise.AppInfo{}) {
		return nil, errors.New("No artifact app info found for artifact")
	}
	if reflect.DeepEqual(selectedArtifact.ArtifactMeta.ProvisioningInfo, bitrise.ProvisioningInfo{}) {
		return nil, errors.New("No artifact provisioning info found for artifact")
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
	artifactInfo := models.ArtifactInfo{
		Version:              selectedArtifact.ArtifactMeta.AppInfo.Version,
		MinimumOS:            selectedArtifact.ArtifactMeta.AppInfo.MinimumOS,
		BundleID:             selectedArtifact.ArtifactMeta.AppInfo.BundleID,
		SupportedDeviceTypes: supportedDeviceTypes,
		ExpireDate:           selectedArtifact.ArtifactMeta.ProvisioningInfo.ExpireDate,
		DistributionType:     selectedArtifact.ArtifactMeta.ProvisioningInfo.DistributionType,
	}
	artifactInfoData, err := json.Marshal(artifactInfo)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &models.AppVersion{
		Platform:         "ios",
		BuildSlug:        buildSlug,
		ArtifactInfoData: artifactInfoData,
	}, nil
}

func selectIosArtifact(artifacts []bitrise.ArtifactListElementResponseModel) (*bitrise.ArtifactListElementResponseModel, bool, bool, string) {
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
	return selectedArtifact, publishEnabled, publicInstallPageEnabled, publicInstallPageArtifactSlug
}

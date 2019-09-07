package services

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/pkg/errors"
)

func prepareAppVersionForIosPlatform(w http.ResponseWriter, r *http.Request, artifacts []bitrise.ArtifactListElementResponseModel, buildSlug string) (*models.AppVersion, error) {
	selectedArtifact, _, _, _, _ := selectIosArtifact(artifacts)
	if selectedArtifact == nil || reflect.DeepEqual(*selectedArtifact, bitrise.ArtifactListElementResponseModel{}) {
		return nil, errors.New("No iOS artifact found")
	}

	if selectedArtifact.ArtifactMeta == nil {
		return nil, errors.New("No artifact meta data found for artifact")
	}

	if reflect.DeepEqual(selectedArtifact.ArtifactMeta.AppInfo, bitrise.AppInfo{}) {
		return nil, errors.New("No artifact app info found for artifact")
	}
	// if reflect.DeepEqual(selectedArtifact.ArtifactMeta.ProvisioningInfo, bitrise.ProvisioningInfo{}) {
	// 	return nil, errors.New("No artifact provisioning info found for artifact")
	// }

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
		// ExpireDate:           selectedArtifact.ArtifactMeta.ProvisioningInfo.ExpireDate,
	}
	artifactInfoData, err := json.Marshal(artifactInfo)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &models.AppVersion{
		Platform:         "ios",
		BuildSlug:        buildSlug,
		ArtifactInfoData: artifactInfoData,
		Scheme:           selectedArtifact.ArtifactMeta.Scheme,
	}, nil
}

func hasIosArtifact(artifacts []bitrise.ArtifactListElementResponseModel) bool {
	for _, artifact := range artifacts {
		if artifact.IsIPA() || artifact.IsXCodeArchive() {
			return true
		}
	}

	return false
}

func selectIosArtifact(artifacts []bitrise.ArtifactListElementResponseModel) (*bitrise.ArtifactListElementResponseModel, bool, bool, string, string) {
	publishEnabled := false
	publicInstallPageEnabled := false
	ipaIPAExportMethod := ""
	publicInstallPageArtifactSlug := ""
	var selectedArtifact bitrise.ArtifactListElementResponseModel
	for _, artifact := range artifacts {
		if artifact.IsIPA() {
			if artifact.ArtifactMeta != nil && artifact.ArtifactMeta.ProvisioningInfo.IPAExportMethod != "" {
				ipaIPAExportMethod = artifact.ArtifactMeta.ProvisioningInfo.IPAExportMethod
			}

			if artifact.HasAppStoreIPAExportMethod() {
				publishEnabled = true
			}
			if artifact.HasDebugIPAExportMethod() {
				publicInstallPageEnabled = true
				publicInstallPageArtifactSlug = artifact.Slug
			}
		}
		if artifact.IsXCodeArchive() {
			publishEnabled = true
			selectedArtifact = artifact
		}
	}
	return &selectedArtifact, publishEnabled, publicInstallPageEnabled, ipaIPAExportMethod, publicInstallPageArtifactSlug
}

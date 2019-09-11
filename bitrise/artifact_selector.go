package bitrise

import (
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/pkg/errors"
)

// ArtifactSelector ...
type ArtifactSelector struct {
	artifacts []ArtifactListElementResponseModel
}

// NewArtifactSelector ...
func NewArtifactSelector(artifacts []ArtifactListElementResponseModel) ArtifactSelector {
	return ArtifactSelector{
		artifacts: artifacts,
	}
}

// PrepareAndroidAppVersions ...
func (s *ArtifactSelector) PrepareAndroidAppVersions(buildSlug, buildNumber, buildCommitMessage, module string) ([]models.AppVersion, error, error) {
	appVersions := []models.AppVersion{}
	artifacts, settingErr, err := pickArtifactsByModule(s.artifacts, module)
	if settingErr != nil {
		return nil, settingErr, nil
	}
	if err != nil {
		return nil, nil, err
	}
	flavorGroups := map[string][]ArtifactListElementResponseModel{}
	for _, artifact := range artifacts {
		if artifact.ArtifactMeta != nil {
			flavorGroups[artifact.ArtifactMeta.ProductFlavour] = append(flavorGroups[artifact.ArtifactMeta.ProductFlavour], artifact)
		} else {
			return nil, nil, errors.New("No artifact meta data found for artifact")
		}
	}
	groupKeys := []string{}
	for key := range flavorGroups {
		groupKeys = append(groupKeys, key)
	}
	sort.Strings(groupKeys)

	for _, key := range groupKeys {
		group := flavorGroups[key]
		buildTypeGroups := groupByBuildType(group)
		keys := []string{}
		for key := range buildTypeGroups {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		artifactInfo := models.ArtifactInfo{
			MinimumSDK:  group[0].ArtifactMeta.AppInfo.MinimumSDKVersion,
			PackageName: group[0].ArtifactMeta.AppInfo.PackageName,
			Version:     group[0].ArtifactMeta.AppInfo.VersionName,
			VersionCode: group[0].ArtifactMeta.AppInfo.VersionCode,
			Module:      group[0].ArtifactMeta.Module,
			BuildType:   strings.Join(keys, ", "),
		}
		artifactInfoData, err := json.Marshal(artifactInfo)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
		appVersion := models.AppVersion{
			Platform:         "android",
			BuildSlug:        buildSlug,
			BuildNumber:      buildNumber,
			ArtifactInfoData: artifactInfoData,
			LastUpdate:       time.Now(),
			CommitMessage:    buildCommitMessage,
			ProductFlavour:   group[0].ArtifactMeta.ProductFlavour,
		}
		appVersions = append(appVersions, appVersion)
	}
	return appVersions, nil, nil
}

func groupByBuildType(artifacts []ArtifactListElementResponseModel) map[string][]ArtifactListElementResponseModel {
	buildTypeGroups := map[string][]ArtifactListElementResponseModel{}
	for _, artifact := range artifacts {
		buildTypeGroups[artifact.ArtifactMeta.BuildType] = append(buildTypeGroups[artifact.ArtifactMeta.BuildType], artifact)
	}
	return buildTypeGroups
}

func groupByModule(artifacts []ArtifactListElementResponseModel) (map[string][]ArtifactListElementResponseModel, error) {
	moduleGroups := map[string][]ArtifactListElementResponseModel{}
	for _, artifact := range artifacts {
		if artifact.ArtifactMeta != nil {
			moduleGroups[artifact.ArtifactMeta.Module] = append(moduleGroups[artifact.ArtifactMeta.Module], artifact)
		} else {
			return nil, errors.New("No artifact meta data found for artifact")
		}
	}
	return moduleGroups, nil
}

func pickArtifactsByModule(artifacts []ArtifactListElementResponseModel, module string) ([]ArtifactListElementResponseModel, error, error) {
	pickedArtifacts := []ArtifactListElementResponseModel{}
	moduleGroups, err := groupByModule(artifacts)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	if len(moduleGroups) == 1 {
		for key := range moduleGroups {
			pickedArtifacts = moduleGroups[key]
			break
		}
	} else if len(moduleGroups) > 1 && module == "" {
		return nil, errors.New("No module setting found"), nil
	} else if len(moduleGroups) > 1 {
		pickedArtifacts = moduleGroups[module]
	}
	return pickedArtifacts, nil, nil
}

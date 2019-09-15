package bitrise

import (
	"encoding/json"
	"path/filepath"
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
	artifacts, settingErr := pickArtifactsByModule(s.artifacts, module)
	if settingErr != nil {
		return nil, settingErr, nil
	}
	flavourGroups := map[string][]ArtifactListElementResponseModel{}
	for _, artifact := range artifacts {
		if artifact.ArtifactMeta != nil {
			flavourGroups[artifact.ArtifactMeta.ProductFlavour] = append(flavourGroups[artifact.ArtifactMeta.ProductFlavour], artifact)
		}
	}
	groupKeys := []string{}
	for key := range flavourGroups {
		groupKeys = append(groupKeys, key)
	}
	sort.Strings(groupKeys)

	for _, key := range groupKeys {
		group := flavourGroups[key]
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

// Select ...
func (s *ArtifactSelector) Select(module, flavour string) ([]string, error) {
	artifactSlugs := []string{}
	artifacts, settingErr := pickArtifactsByModule(s.artifacts, module)
	if settingErr != nil {
		return nil, settingErr
	}
	buildTypeGroups := groupByBuildType(artifacts)
	artifacts = buildTypeGroups["release"]
	flavourGroups := groupByFlavour(artifacts)
	selectedGroup := flavourGroups[flavour]
	for _, artifact := range selectedGroup {
		if artifact.IsStandaloneAPK() || artifact.IsUniversalAPK() || artifact.IsAAB() {
			artifactSlugs = append(artifactSlugs, artifact.Slug)
			continue
		}
		if len(artifact.ArtifactMeta.Split) > 0 && artifact.ArtifactMeta.Aab == "" {
			artifactSlugs = append(artifactSlugs, artifact.Slug)
			continue
		}
	}

	return artifactSlugs, nil
}

// PublishAndShareInfo ...
func (s *ArtifactSelector) PublishAndShareInfo(appVersion *models.AppVersion) (bool, bool, string, bool, bool, error) {
	publishEnabled := false
	publicInstallPageEnabled := false
	publicInstallPageArtifactSlug := ""
	split := false
	universalAvailable := false
	artifactInfo, err := appVersion.ArtifactInfo()
	if err != nil {
		return false, false, "", false, false, errors.WithStack(err)
	}
	if artifactInfo.BuildType == "release" {
		publishEnabled = true
	}
	for _, artifact := range s.artifacts {
		if artifact.ArtifactMeta != nil {
			if artifact.ArtifactMeta.ProductFlavour == appVersion.ProductFlavour &&
				artifact.ArtifactMeta.BuildType == artifactInfo.BuildType &&
				artifact.ArtifactMeta.Module == artifactInfo.Module {
				if artifact.IsUniversalAPK() {
					universalAvailable = true
					if artifact.IsPublicPageEnabled {
						publicInstallPageEnabled = true
						publicInstallPageArtifactSlug = artifact.Slug
					}
				}
				if artifact.IsStandaloneAPK() {
					if artifact.IsPublicPageEnabled {
						publicInstallPageEnabled = true
						publicInstallPageArtifactSlug = artifact.Slug
					}
				}
				if len(artifact.ArtifactMeta.Split) > 0 {
					split = true
				}
			}
		}
	}
	return publishEnabled, publicInstallPageEnabled, publicInstallPageArtifactSlug, split, universalAvailable, nil
}

// HasAndroidArtifact ...
func (s *ArtifactSelector) HasAndroidArtifact() bool {
	for _, artifact := range s.artifacts {
		ext := filepath.Ext(artifact.Title)
		for _, androidExt := range []string{".apk", ".aab"} {
			if ext == androidExt {
				return true
			}
		}
	}
	return false
}

func groupByBuildType(artifacts []ArtifactListElementResponseModel) map[string][]ArtifactListElementResponseModel {
	buildTypeGroups := map[string][]ArtifactListElementResponseModel{}
	for _, artifact := range artifacts {
		buildTypeGroups[artifact.ArtifactMeta.BuildType] = append(buildTypeGroups[artifact.ArtifactMeta.BuildType], artifact)
	}
	return buildTypeGroups
}
func groupByFlavour(artifacts []ArtifactListElementResponseModel) map[string][]ArtifactListElementResponseModel {
	flavourGroups := map[string][]ArtifactListElementResponseModel{}
	for _, artifact := range artifacts {
		if artifact.ArtifactMeta != nil {
			flavourGroups[artifact.ArtifactMeta.ProductFlavour] = append(flavourGroups[artifact.ArtifactMeta.ProductFlavour], artifact)
		}
	}
	return flavourGroups
}

func groupByModule(artifacts []ArtifactListElementResponseModel) map[string][]ArtifactListElementResponseModel {
	moduleGroups := map[string][]ArtifactListElementResponseModel{}
	for _, artifact := range artifacts {
		if artifact.ArtifactMeta != nil {
			moduleGroups[artifact.ArtifactMeta.Module] = append(moduleGroups[artifact.ArtifactMeta.Module], artifact)
		}
	}
	return moduleGroups
}

func pickArtifactsByModule(artifacts []ArtifactListElementResponseModel, module string) ([]ArtifactListElementResponseModel, error) {
	pickedArtifacts := []ArtifactListElementResponseModel{}
	moduleGroups := groupByModule(artifacts)
	if len(moduleGroups) == 1 {
		for key := range moduleGroups {
			pickedArtifacts = moduleGroups[key]
			break
		}
	} else if len(moduleGroups) > 1 && module == "" {
		return nil, errors.New("No module setting found")
	} else if len(moduleGroups) > 1 {
		pickedArtifacts = moduleGroups[module]
	}
	return pickedArtifacts, nil
}

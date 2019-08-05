package bitrise

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"time"
)

var debugDistributionTypes = [...]string{"development", "ad-hoc"}

// AppInfo ...
type AppInfo struct {
	AppName           string `json:"app_name"`
	MinimumOS         string `json:"min_OS_version"`
	MinimumSDKVersion string `json:"min_sdk_version"`
	BundleID          string `json:"bundle_id"`
	DeviceFamilyList  []int  `json:"device_family_list"`
	PackageName       string `json:"package_name"`
	VersionName       string `json:"version_name"`
	Version           string `json:"version"`
}

// ProvisioningInfo ...
type ProvisioningInfo struct {
	ExpireDate       time.Time `json:"expire_date"`
	DistributionType string    `json:"distribution_type"`
}

// ArtifactMeta ...
type ArtifactMeta struct {
	AppInfo          AppInfo          `json:"app_info"`
	ProvisioningInfo ProvisioningInfo `json:"provisioning_info"`
	Size             string           `json:"file_size_bytes"`
}

// ArtifactData ...
type ArtifactData struct {
	Meta ArtifactMeta
	Slug string
}

// ArtifactListElementResponseModel ....
type ArtifactListElementResponseModel struct {
	Title               string        `json:"title"`
	ArtifactType        *string       `json:"artifact_type"`
	ArtifactMeta        *ArtifactMeta `json:"artifact_meta"`
	IsPublicPageEnabled bool          `json:"is_public_page_enabled"`
	Slug                string        `json:"slug"`
	FileSizeBytes       *int64        `json:"file_size_bytes"`
}

type artifactListResponseModel struct {
	Data   []ArtifactListElementResponseModel `json:"data"`
	Paging pagingResponseModel                `json:"paging"`
}

// ArtifactShowResponseItemModel ...
type ArtifactShowResponseItemModel struct {
	Title                *string         `json:"title"`
	ArtifactType         *string         `json:"artifact_type"`
	ArtifactMeta         json.RawMessage `json:"artifact_meta"`
	DownloadPath         *string         `json:"expiring_download_url"`
	IsPublicPageEnabled  bool            `json:"is_public_page_enabled"`
	Slug                 string          `json:"slug"`
	PublicInstallPageURL string          `json:"public_install_page_url"`
	FileSizeBytes        *int64          `json:"file_size_bytes"`
}

type artifactShowResponseModel struct {
	Data ArtifactShowResponseItemModel `json:"data"`
}

// HasDebugDistributionType ...
func (a ArtifactListElementResponseModel) HasDebugDistributionType() bool {
	if a.ArtifactMeta == nil ||
		a.ArtifactMeta.ProvisioningInfo == (ProvisioningInfo{}) ||
		a.ArtifactMeta.ProvisioningInfo.DistributionType == "" {
		return false
	}
	for _, artifactType := range debugDistributionTypes {
		if a.ArtifactMeta.ProvisioningInfo.DistributionType == artifactType {
			return true
		}
	}
	return false
}

// HasAppStoreDistributionType ...
func (a ArtifactListElementResponseModel) HasAppStoreDistributionType() bool {
	if a.ArtifactMeta == nil || a.ArtifactMeta.ProvisioningInfo == (ProvisioningInfo{}) {
		return false
	}
	return a.ArtifactMeta.ProvisioningInfo.DistributionType == "app-store"
}

// IsIPA ...
func (a ArtifactListElementResponseModel) IsIPA() bool {
	return filepath.Ext(a.Title) == ".ipa"
}

// IsXCodeArchive ...
func (a ArtifactListElementResponseModel) IsXCodeArchive() bool {
	return strings.Contains(strings.ToLower(a.Title), "xcarchive") && filepath.Ext(a.Title) == ".zip"
}

// IsAAB ...
func (a ArtifactListElementResponseModel) IsAAB() bool {
	return filepath.Ext(a.Title) == ".aab"
}

// IsUniversalAPK ...
func (a ArtifactListElementResponseModel) IsUniversalAPK() bool {
	return strings.Contains(strings.ToLower(a.Title), "universal") && filepath.Ext(a.Title) == ".apk"
}

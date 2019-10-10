package bitrise

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"time"
)

var debugIPAExportMethods = [...]string{"development", "ad-hoc"}

// AppInfo ...
type AppInfo struct {
	AppName           string `json:"app_name"`
	MinimumOS         string `json:"min_OS_version"`
	MinimumSDKVersion string `json:"min_sdk_version"`
	BundleID          string `json:"bundle_id"`
	BuildNumber       string `json:"build_number"`
	DeviceFamilyList  []int  `json:"device_family_list"`
	PackageName       string `json:"package_name"`
	VersionName       string `json:"version_name"`
	VersionCode       string `json:"version_code"`
	Version           string `json:"version"`
}

// ProvisioningInfo ...
type ProvisioningInfo struct {
	ExpireDate      time.Time `json:"expire_date"`
	IPAExportMethod string    `json:"ipa_export_method"`
}

// ArtifactMeta ...
type ArtifactMeta struct {
	AppInfo          AppInfo          `json:"app_info"`
	ProvisioningInfo ProvisioningInfo `json:"provisioning_info"`
	Size             string           `json:"file_size_bytes"`
	Scheme           string           `json:"scheme"`
	Module           string           `json:"module"`
	ProductFlavor    string           `json:"product_flavour"`
	BuildType        string           `json:"build_type"`
	Include          bool             `json:"include"`
	Universal        string           `json:"universal"`
	Aab              string           `json:"aab"`
	Apk              string           `json:"apk"`
	Split            []string         `json:"split"`
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

// HasDebugIPAExportMethod ...
func (a ArtifactListElementResponseModel) HasDebugIPAExportMethod() bool {
	if a.ArtifactMeta == nil ||
		a.ArtifactMeta.ProvisioningInfo == (ProvisioningInfo{}) ||
		a.ArtifactMeta.ProvisioningInfo.IPAExportMethod == "" {
		return false
	}
	for _, artifactType := range debugIPAExportMethods {
		if a.ArtifactMeta.ProvisioningInfo.IPAExportMethod == artifactType {
			return true
		}
	}
	return false
}

// HasAppStoreIPAExportMethod ...
func (a ArtifactListElementResponseModel) HasAppStoreIPAExportMethod() bool {
	if a.ArtifactMeta == nil || a.ArtifactMeta.ProvisioningInfo == (ProvisioningInfo{}) {
		return false
	}
	return a.ArtifactMeta.ProvisioningInfo.IPAExportMethod == "app-store"
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
	if a.ArtifactMeta == nil {
		return false
	}
	return filepath.Base(a.ArtifactMeta.Aab) == a.Title
}

// IsStandaloneAPK ...
func (a ArtifactListElementResponseModel) IsStandaloneAPK() bool {
	if a.ArtifactMeta == nil {
		return false
	}
	return filepath.Base(a.ArtifactMeta.Apk) == a.Title
}

// IsUniversalAPK ...
func (a ArtifactListElementResponseModel) IsUniversalAPK() bool {
	if a.ArtifactMeta == nil {
		return false
	}
	return filepath.Base(a.ArtifactMeta.Universal) == a.Title
}

package bitrise

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/nulls"
)

// AppInfo ...
type AppInfo struct {
	MinimumOS         string `json:"min_OS_version"`
	MinimumSDKVersion string `json:"min_sdk_version"`
	BundleID          string `json:"bundle_id"`
	DeviceFamilyList  []int  `json:"device_family_list"`
	PackageName       string `json:"package_name"`
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

type artifactListElementResponseModel struct {
	Title               *string         `json:"title"`
	ArtifactType        *string         `json:"artifact_type"`
	ArtifactMeta        json.RawMessage `json:"artifact_meta"`
	IsPublicPageEnabled bool            `json:"is_public_page_enabled"`
	Slug                string          `json:"slug"`
	FileSizeBytes       *int64          `json:"file_size_bytes"`
}

type artifactListResponseModel struct {
	Data   []artifactListElementResponseModel `json:"data"`
	Paging pagingResponseModel                `json:"paging"`
}

type artifactShowResponseItemModel struct {
	Title                nulls.String    `json:"title"`
	ArtifactType         nulls.String    `json:"artifact_type"`
	ArtifactMeta         json.RawMessage `json:"artifact_meta"`
	DownloadPath         *string         `json:"expiring_download_url"`
	IsPublicPageEnabled  bool            `json:"is_public_page_enabled"`
	Slug                 string          `json:"slug"`
	PublicInstallPageURL string          `json:"public_install_page_url"`
	FileSizeBytes        *int64          `json:"file_size_bytes"`
}

type artifactShowResponseModel struct {
	Data artifactShowResponseItemModel `json:"data"`
}

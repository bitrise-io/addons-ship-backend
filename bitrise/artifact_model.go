package bitrise

import "time"

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

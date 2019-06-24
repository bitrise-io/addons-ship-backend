package bitrise

// ProvisioningProfile ...
type ProvisioningProfile struct {
	Filename string `json:"upload_file_name"`
	Slug     string `json:"slug"`
}

type provisioningProfileListResponseModel struct {
	ProvisioningProfiles []ProvisioningProfile `json:"data"`
}

// CodeSigningIdentity ...
type CodeSigningIdentity struct {
	Filename string `json:"upload_file_name"`
	Slug     string `json:"slug"`
}

type codeSigningIdentityListResponseModel struct {
	CodeSigningIdentities []CodeSigningIdentity `json:"data"`
}

// AndroidKeystoreFile ...
type AndroidKeystoreFile struct {
	Filename string `json:"upload_file_name"`
	Slug     string `json:"slug"`
}

type androidKeystoreFileListResponseModel struct {
	AndroidKeystoreFiles []AndroidKeystoreFile `json:"data"`
}

// GenericProjectFile ...
type GenericProjectFile struct {
	Filename string `json:"upload_file_name"`
	Slug     string `json:"slug"`
}

type genericProjectFileListResponseModel struct {
	GenericProjectFiles []GenericProjectFile `json:"data"`
}

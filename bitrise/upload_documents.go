package bitrise

// ProvisioningProfile ...
type ProvisioningProfile struct {
	Filename    string `json:"upload_file_name"`
	Slug        string `json:"slug"`
	DownloadURL string `json:"download_url"`
}

type provisioningProfileListResponseModel struct {
	ProvisioningProfiles []ProvisioningProfile `json:"data"`
}

type provisioningProfileShowResponseModel struct {
	Data ProvisioningProfile `json:"data"`
}

// CodeSigningIdentity ...
type CodeSigningIdentity struct {
	Filename            string `json:"upload_file_name"`
	Slug                string `json:"slug"`
	CertificatePassword string `json:"certificate_password"`
	DownloadURL         string `json:"download_url"`
}

type codeSigningIdentityListResponseModel struct {
	CodeSigningIdentities []CodeSigningIdentity `json:"data"`
}

type codeSigningIdentityShowResponseModel struct {
	Data CodeSigningIdentity `json:"data"`
}

// ExposedMetadataStore ...
type ExposedMetadataStore struct {
	Password           string `json:"PASSWORD"`
	Alias              string `json:"ALIAS"`
	PrivateKeyPassword string `json:"PRIVATE_KEY_PASSWORD"`
}

// AndroidKeystoreFile ...
type AndroidKeystoreFile struct {
	Filename             string               `json:"upload_file_name"`
	Slug                 string               `json:"slug"`
	DownloadURL          string               `json:"download_url"`
	UserEnvKey           string               `json:"user_env_key"`
	ExposedMetadataStore ExposedMetadataStore `json:"exposed_meta_datastore"`
}

type androidKeystoreFileListResponseModel struct {
	AndroidKeystoreFiles []AndroidKeystoreFile `json:"data"`
}

type androidKeystoreFileShowResponseModel struct {
	Data *AndroidKeystoreFile `json:"data"`
}

// GenericProjectFile ...
type GenericProjectFile struct {
	Filename    string `json:"upload_file_name"`
	Slug        string `json:"slug"`
	DownloadURL string `json:"download_url"`
}

type genericProjectFileListResponseModel struct {
	GenericProjectFiles []GenericProjectFile `json:"data"`
}

type genericProjectFileShowResponseModel struct {
	Data *GenericProjectFile `json:"data"`
}

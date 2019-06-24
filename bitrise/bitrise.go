package bitrise

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

var (
	apiBaseURL         = "https://api.bitrise.io/v0.1"
	validArtifactTypes = [...]string{"android-apk", "ios-ipa"}
)

// APIInterface ...
type APIInterface interface {
	GetArtifactData(string, string, string) (*ArtifactData, error)
	GetArtifactPublicInstallPageURL(string, string, string, string) (string, error)
	GetAppDetails(authToken, appSlug string) (*AppDetails, error)
	GetProvisioningProfiles(authToken, appSlug string) ([]ProvisioningProfile, error)
	GetCodeSigningIdentities(authToken, appSlug string) ([]CodeSigningIdentity, error)
	GetAndroidKeystoreFiles(authToken, appSlug string) ([]AndroidKeystoreFile, error)
	GetServiceAccountFiles(authToken, appSlug string) ([]GenericProjectFile, error)
}

// API ...
type API struct {
	*http.Client
	url string
}

// New ...
func New() *API {
	url, ok := os.LookupEnv("BITRISE_API_ROOT_URL")
	if !ok {
		url = apiBaseURL
	}
	return &API{
		Client: &http.Client{},
		url:    url,
	}
}

func (a *API) doRequest(authToken, method, path string) (*http.Response, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", a.url, path), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Set("Bitrise-Addon-Auth-Token", authToken)
	resp, err := a.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Failed to fetch artifact data: status: %d", resp.StatusCode)
	}
	return resp, nil
}

// GetArtifactData ...
func (a *API) GetArtifactData(authToken, appSlug, buildSlug string) (*ArtifactData, error) {
	responseModel, err := a.listArtifacts(authToken, appSlug, buildSlug, "")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if responseModel.Paging.Next == "" {
		artifactData, err := getInstallableArtifactsFromResponseModel(responseModel)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if artifactData == nil {
			return nil, errors.New("No matching artifact found")
		}
		return artifactData, nil
	}
	next := responseModel.Paging.Next
	for next != "" {
		responseModel, err = a.listArtifacts(authToken, appSlug, buildSlug, next)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		artifactData, err := getInstallableArtifactsFromResponseModel(responseModel)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if artifactData != nil {
			return artifactData, nil
		}
		next = responseModel.Paging.Next
	}
	return nil, errors.New("No matching artifact found")
}

// GetArtifactPublicInstallPageURL ...
func (a *API) GetArtifactPublicInstallPageURL(authToken, appSlug, buildSlug, artifactSlug string) (string, error) {
	resp, err := a.doRequest(authToken, "GET", fmt.Sprintf("/apps/%s/builds/%s/artifacts/%s", appSlug, buildSlug, artifactSlug))
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer httpresponse.BodyCloseWithErrorLog(resp)
	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("Failed to fetch artifact data: status: %d", resp.StatusCode)
	}
	var responseModel artifactShowResponseModel
	if err := json.NewDecoder(resp.Body).Decode(&responseModel); err != nil {
		return "", errors.WithStack(err)
	}
	return responseModel.Data.PublicInstallPageURL, nil
}

// GetAppDetails ...
func (a *API) GetAppDetails(authToken, appSlug string) (*AppDetails, error) {
	resp, err := a.doRequest(authToken, "GET", "/apps/"+appSlug)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer httpresponse.BodyCloseWithErrorLog(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Failed to fetch artifact data: status: %d", resp.StatusCode)
	}
	var responseModel appShowResponseModel
	if err := json.NewDecoder(resp.Body).Decode(&responseModel); err != nil {
		return nil, errors.WithStack(err)
	}
	return &responseModel.Data, nil
}

// GetProvisioningProfiles ...
func (a *API) GetProvisioningProfiles(authToken, appSlug string) ([]ProvisioningProfile, error) {
	resp, err := a.doRequest(authToken, "GET", fmt.Sprintf("/apps/%s/provisioning-profiles", appSlug))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer httpresponse.BodyCloseWithErrorLog(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Failed to fetch artifact data: status: %d", resp.StatusCode)
	}
	var responseModel provisioningProfileListResponseModel
	if err := json.NewDecoder(resp.Body).Decode(&responseModel); err != nil {
		return nil, errors.WithStack(err)
	}
	return responseModel.ProvisioningProfiles, nil
}

// GetCodeSigningIdentities ...
func (a *API) GetCodeSigningIdentities(authToken, appSlug string) ([]CodeSigningIdentity, error) {
	resp, err := a.doRequest(authToken, "GET", fmt.Sprintf("/apps/%s/build-certificates", appSlug))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer httpresponse.BodyCloseWithErrorLog(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Failed to fetch artifact data: status: %d", resp.StatusCode)
	}
	var responseModel codeSigningIdentityListResponseModel
	if err := json.NewDecoder(resp.Body).Decode(&responseModel); err != nil {
		return nil, errors.WithStack(err)
	}
	return responseModel.CodeSigningIdentities, nil
}

// GetAndroidKeystoreFiles ...
func (a *API) GetAndroidKeystoreFiles(authToken, appSlug string) ([]AndroidKeystoreFile, error) {
	resp, err := a.doRequest(authToken, "GET", fmt.Sprintf("/apps/%s/android-keystore-files", appSlug))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer httpresponse.BodyCloseWithErrorLog(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Failed to fetch artifact data: status: %d", resp.StatusCode)
	}
	var responseModel androidKeystoreFileListResponseModel
	if err := json.NewDecoder(resp.Body).Decode(&responseModel); err != nil {
		return nil, errors.WithStack(err)
	}
	return responseModel.AndroidKeystoreFiles, nil
}

// GetServiceAccountFiles ...
func (a *API) GetServiceAccountFiles(authToken, appSlug string) ([]GenericProjectFile, error) {
	resp, err := a.doRequest(authToken, "GET", fmt.Sprintf("/apps/%s/generic-project-files", appSlug))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer httpresponse.BodyCloseWithErrorLog(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Failed to fetch artifact data: status: %d", resp.StatusCode)
	}
	var responseModel genericProjectFileListResponseModel
	if err := json.NewDecoder(resp.Body).Decode(&responseModel); err != nil {
		return nil, errors.WithStack(err)
	}
	serviceAccountFiles := []GenericProjectFile{}
	for _, genFile := range responseModel.GenericProjectFiles {
		if filepath.Ext(genFile.Filename) == ".json" {
			serviceAccountFiles = append(serviceAccountFiles, genFile)
		}
	}

	return serviceAccountFiles, nil
}

func (a *API) listArtifacts(authToken, appSlug, buildSlug, next string) (*artifactListResponseModel, error) {
	path := fmt.Sprintf("/apps/%s/builds/%s/artifacts", appSlug, buildSlug)
	if next != "" {
		path = fmt.Sprintf("%s?next=%s", path, next)
	}
	resp, err := a.doRequest(authToken, "GET", path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer httpresponse.BodyCloseWithErrorLog(resp)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Failed to fetch artifact data: status: %d", resp.StatusCode)
	}
	var responseModel artifactListResponseModel
	if err := json.NewDecoder(resp.Body).Decode(&responseModel); err != nil {
		return nil, errors.WithStack(err)
	}
	return &responseModel, nil
}

func getInstallableArtifactsFromResponseModel(respModel *artifactListResponseModel) (*ArtifactData, error) {
	for _, buildArtifact := range respModel.Data {
		if validArtifact(buildArtifact) {
			var artifactMeta ArtifactMeta
			err := json.Unmarshal([]byte(buildArtifact.ArtifactMeta), &artifactMeta)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return &ArtifactData{
				Meta: artifactMeta,
				Slug: buildArtifact.Slug,
			}, nil
		}
	}
	return nil, nil
}

func validArtifact(artifact artifactListElementResponseModel) bool {
	for _, artifactType := range validArtifactTypes {
		if artifact.ArtifactType == nil {
			return false
		}
		if artifactType == *artifact.ArtifactType {
			return true
		}
	}
	return false
}

package bitrise

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

var (
	apiBaseURL         = "https://api.bitrise.io/v0.1"
	validArtifactTypes = [...]string{"android-apk", "ios-ipa"}
)

// APIInterface ...
type APIInterface interface {
	GetArtifactMetadata(authToken, appSlug, buildSlug string) (*ArtifactMeta, error)
}

// API ...
type API struct {
	*http.Client
}

// New ...
func New() *API {
	return &API{
		Client: &http.Client{},
	}
}

func (a *API) doRequest(authToken, method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
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

// GetArtifactMetadata ...
func (a *API) GetArtifactMetadata(authToken, appSlug, buildSlug string) (*ArtifactMeta, error) {
	responseModel, err := a.listArtifacts(authToken, appSlug, buildSlug, "")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if responseModel.Paging.Next == "" {
		artifactMeta, err := getInstallableArtifactsFromResponseModel(responseModel)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if artifactMeta == nil {
			return nil, errors.New("No matching artifact found")
		}
		return artifactMeta, nil
	}
	next := responseModel.Paging.Next
	for next != "" {
		responseModel, err = a.listArtifacts(authToken, appSlug, buildSlug, next)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		artifacMeta, err := getInstallableArtifactsFromResponseModel(responseModel)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if artifacMeta != nil {
			return artifacMeta, nil
		}
		next = responseModel.Paging.Next
	}
	return nil, errors.New("No matching artifact found")
}

func (a *API) listArtifacts(authToken, appSlug, buildSlug, next string) (*artifactListResponseModel, error) {
	url := fmt.Sprintf("%s/apps/%s/builds/%s/artifacts", apiBaseURL, appSlug, buildSlug)
	if next != "" {
		url = fmt.Sprintf("%s?next=%s", url, next)
	}
	resp, err := a.doRequest(authToken, "GET", url)
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

func getInstallableArtifactsFromResponseModel(respModel *artifactListResponseModel) (*ArtifactMeta, error) {
	for _, buildArtifact := range respModel.Data {
		if validArtifact(buildArtifact) {
			var artifactMeta ArtifactMeta
			err := json.Unmarshal([]byte(buildArtifact.ArtifactMeta), &artifactMeta)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return &artifactMeta, nil
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

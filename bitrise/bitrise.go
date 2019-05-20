package bitrise

import (
	"encoding/json"

	bitriseapiclient "github.com/bitrise-io/bitrise-api-client/client"
	"github.com/bitrise-io/bitrise-api-client/client/build_artifact"
	"github.com/bitrise-io/bitrise-api-client/models"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/pkg/errors"
)

var validArtifactTypes = [...]string{"android-apk", "ios-ipa"}

// APIInterface ...
type APIInterface interface {
	GetArtifactMetadata(authToken, appSlug, buildSlug string) (*ArtifactMeta, error)
}

// API ...
type API struct {
	*bitriseapiclient.Bitrise
}

// New ...
func New() *API {
	return &API{
		Bitrise: bitriseapiclient.Default,
	}
}

func validArtifact(artifact *models.V0ArtifactListElementResponseModel) bool {
	for _, artifactType := range validArtifactTypes {
		if artifactType == artifact.ArtifactType {
			return true
		}
	}
	return false
}

// GetArtifactMetadata ...
func (a *API) GetArtifactMetadata(authToken, appSlug, buildSlug string) (*ArtifactMeta, error) {
	buildArtifacts, err := a.BuildArtifact.ArtifactList(&build_artifact.ArtifactListParams{
		AppSlug: appSlug, BuildSlug: buildSlug,
	}, httptransport.APIKeyAuth("Bitrise-Addon-Auth-Token", "header", authToken))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, buildArtifact := range buildArtifacts.Payload.Data {
		if validArtifact(buildArtifact) {
			var artifactMeta ArtifactMeta
			// use artifact meta!!!!!
			err := json.Unmarshal([]byte(buildArtifact.ArtifactType), &artifactMeta)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return &artifactMeta, nil
		}
	}
	return nil, errors.New("No installable artifact found")
}

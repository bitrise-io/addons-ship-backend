package bitrise

// BuildDetails ...
type BuildDetails struct {
	CommitMessage string `json:"commit_message"`
}

type buildShowResponseModel struct {
	Data BuildDetails `json:"data"`
}

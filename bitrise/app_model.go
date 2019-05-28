package bitrise

// AppDetails ...
type AppDetails struct {
	Title     string  `json:"title"`
	AvatarURL *string `json:"avatar_url"`
}

type appShowResponseModel struct {
	Data AppDetails `json:"data"`
}

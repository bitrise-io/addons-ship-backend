package models

// UploadableObject ...
type UploadableObject struct {
	Filename string `json:"filename"`
	Filesize int64  `json:"filesize"`
	Uploaded bool   `json:"uploaded"`
}

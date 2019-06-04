package models

// Uploadable ...
type Uploadable struct {
	Filename string `json:"filename"`
	Filesize int64  `json:"filesize"`
	Uploaded bool   `json:"uploaded"`
}

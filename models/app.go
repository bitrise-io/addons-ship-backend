package models

import uuid "github.com/satori/go.uuid"

// App ...
type App struct {
	Model
	AppSlug         string
	Plan            string
	BitriseAPIToken string
	APIToken        string
}

// BeforeCreate ...
func (a *App) BeforeCreate() error {
	a.ID = uuid.NewV4()
	return nil
}

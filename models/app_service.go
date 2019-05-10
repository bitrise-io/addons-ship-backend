package models

import "github.com/jinzhu/gorm"

// AppService ...
type AppService struct {
	DB *gorm.DB
}

// Create ...
func (a *AppService) Create(app *App) (*App, error) {
	return &App{}, nil
}

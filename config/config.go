package config

import (
	"os"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/jinzhu/gorm"
)

// AppConfig ...
type AppConfig struct {
	AppService dataservices.AppService
	Port       string
}

// New ...
func New(db *gorm.DB) (conf AppConfig) {
	var ok bool
	conf.Port, ok = os.LookupEnv("PORT")
	if !ok {
		conf.Port = "80"
	}
	conf.AppService = &models.AppService{DB: db}
	return
}

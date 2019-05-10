package config

import (
	"os"

	"github.com/jinzhu/gorm"
)

// AppConfig ...
type AppConfig struct {
	DB   *gorm.DB
	Port string
}

// New ...
func New(db *gorm.DB) (conf AppConfig) {
	var ok bool
	conf.Port, ok = os.LookupEnv("PORT")
	if !ok {
		conf.Port = "80"
	}
	conf.DB = db
	return
}

package config

import "os"

// AppConfig ...
type AppConfig struct {
	Port string
}

// New ...
func New() (conf AppConfig) {
	var ok bool
	conf.Port, ok = os.LookupEnv("PORT")
	if !ok {
		conf.Port = "80"
	}
	return
}

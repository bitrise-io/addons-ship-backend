package env

import (
	"os"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/providers"
	"github.com/jinzhu/gorm"
)

// AppEnv ...
type AppEnv struct {
	AppService        dataservices.AppService
	AppVersionService dataservices.AppVersionService
	Port              string
	Environment       string
	RequestParams     providers.RequestParamsInterface
}

// New ...
func New(db *gorm.DB) (env AppEnv) {
	var ok bool
	env.Port, ok = os.LookupEnv("PORT")
	if !ok {
		env.Port = "80"
	}
	env.Environment, ok = os.LookupEnv("ENVIRONMENT")
	if !ok {
		env.Environment = "development"
	}
	env.AppService = &models.AppService{DB: db}
	env.AppVersionService = &models.AppVersionService{DB: db}
	env.RequestParams = &providers.RequestParamsProvider{}
	return
}

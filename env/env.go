package env

import (
	"os"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/logging"
	"github.com/bitrise-io/api-utils/providers"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	// ServerEnvProduction ...
	ServerEnvProduction = "production"
	// ServerEnvDevelopment ...
	ServerEnvDevelopment = "development"
)

// AppEnv ...
type AppEnv struct {
	Port                  string
	Environment           string
	Logger                *zap.Logger
	AppService            dataservices.AppService
	AppVersionService     dataservices.AppVersionService
	ScreenshotService     dataservices.ScreenshotService
	FeatureGraphicService dataservices.FeatureGraphicService
	BitriseAPI            bitrise.APIInterface
	RequestParams         providers.RequestParamsInterface
	AWS                   providers.AWSInterface
}

// New ...
func New(db *gorm.DB) (*AppEnv, error) {
	var ok bool
	env := &AppEnv{}
	env.Port, ok = os.LookupEnv("PORT")
	if !ok {
		env.Port = "80"
	}
	env.Environment, ok = os.LookupEnv("ENVIRONMENT")
	if !ok {
		env.Environment = ServerEnvDevelopment
	}
	env.Logger = logging.WithContext(nil)
	env.AppService = &models.AppService{DB: db}
	env.AppVersionService = &models.AppVersionService{DB: db}
	env.ScreenshotService = &models.ScreenshotService{DB: db}
	env.FeatureGraphicService = &models.FeatureGraphicService{DB: db}
	if env.Environment == ServerEnvDevelopment {
		env.BitriseAPI = &bitrise.APIDev{}
	} else {
		env.BitriseAPI = bitrise.New()
	}
	env.RequestParams = &providers.RequestParams{}

	awsConfig, err := awsConfig()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	env.AWS = &providers.AWS{Config: awsConfig}
	return env, nil
}

func awsConfig() (providers.AWSConfig, error) {
	awsBucket, ok := os.LookupEnv("AWS_BUCKET")
	if !ok {
		return providers.AWSConfig{}, errors.New("No AWS_BUCKET env var defined")
	}
	awsRegion, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		return providers.AWSConfig{}, errors.New("No AWS_REGION env var defined")
	}
	awsAccessKeyID, ok := os.LookupEnv("AWS_ACCESS_KEY_ID")
	if !ok {
		return providers.AWSConfig{}, errors.New("No AWS_ACCESS_KEY_ID env var defined")
	}
	awsSecretAccessKey, ok := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	if !ok {
		return providers.AWSConfig{}, errors.New("No AWS_SECRET_ACCESS_KEY env var defined")
	}
	return providers.AWSConfig{
		Bucket:          awsBucket,
		Region:          awsRegion,
		AccessKeyID:     awsAccessKeyID,
		SecretAccessKey: awsSecretAccessKey,
	}, nil
}

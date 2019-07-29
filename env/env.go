package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/mailer"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/redis"
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
	Port                   string
	Environment            string
	AddonAccessToken       string
	AddonHostURL           string
	AddonFrontendHostURL   string
	Logger                 *zap.Logger
	AppService             dataservices.AppService
	AppContactService      dataservices.AppContactService
	AppVersionService      dataservices.AppVersionService
	ScreenshotService      dataservices.ScreenshotService
	FeatureGraphicService  dataservices.FeatureGraphicService
	AppSettingsService     dataservices.AppSettingsService
	AppVersionEventService dataservices.AppVersionEventService
	PublishTaskService     dataservices.PublishTaskService
	BitriseAPI             bitrise.APIInterface
	RequestParams          providers.RequestParamsInterface
	AWS                    providers.AWSInterface
	Redis                  redis.Interface
	RedisExpirationTime    int
	LogStoreService        dataservices.LogStore
	WorkerService          dataservices.WorkerService
	Mailer                 mailer.Interface
	EmailConfirmLandingURL string
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
	env.AddonAccessToken, ok = os.LookupEnv("ADDON_ACCESS_TOKEN")
	if !ok {
		return nil, errors.New("No value set for env ADDON_ACCESS_TOKEN")
	}
	env.AddonHostURL, ok = os.LookupEnv("ADDON_HOST_URL")
	if !ok {
		return nil, errors.New("No value set for env ADDON_HOST_URL")
	}
	env.AddonFrontendHostURL, ok = os.LookupEnv("ADDON_FRONTEND_HOST_URL")
	if !ok {
		return nil, errors.New("No value set for env ADDON_FRONTEND_HOST_URL")
	}
	env.Logger = logging.WithContext(nil)
	env.AppService = &models.AppService{DB: db}
	env.AppContactService = &models.AppContactService{DB: db}
	env.AppVersionService = &models.AppVersionService{DB: db}
	env.ScreenshotService = &models.ScreenshotService{DB: db}
	env.FeatureGraphicService = &models.FeatureGraphicService{DB: db}
	env.AppSettingsService = &models.AppSettingsService{DB: db}
	env.AppVersionEventService = &models.AppVersionEventService{DB: db}
	env.PublishTaskService = &models.PublishTaskService{DB: db}
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

	redisExpiration := int64(1000)
	redisExpirationStr, ok := os.LookupEnv("REDIS_KEY_EXPIRATION_TIME")
	if ok {
		redisExpiration, err = strconv.ParseInt(redisExpirationStr, 10, 64)
		if err != nil {
			fmt.Println("Invalid Redis expiration time, setting default to 1000 seconds...")
		}
	}
	env.RedisExpirationTime = int(redisExpiration)
	env.Redis = redis.New()
	env.LogStoreService = &models.LogStoreService{Redis: redis.New(), Expiration: env.RedisExpirationTime}

	awsMailRegion, ok := os.LookupEnv("AWS_MAIL_REGION")
	if !ok {
		return nil, errors.New("No AWS_MAIL_REGION env var defined")
	}
	env.Mailer = &mailer.SES{FromEmail: "ship@bitrise.io", Config: providers.AWSConfig{
		Region:          awsMailRegion,
		AccessKeyID:     awsConfig.AccessKeyID,
		SecretAccessKey: awsConfig.SecretAccessKey,
	}}
	env.EmailConfirmLandingURL, ok = os.LookupEnv("EMAIL_CONFIRM_LANDING_URL")
	if !ok {
		return nil, errors.New("No value set for env EMAIL_CONFIRM_LANDING_URL")
	}

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

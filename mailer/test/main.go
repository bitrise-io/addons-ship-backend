package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/mailer"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/providers"
	"github.com/bitrise-io/go-utils/pointers"
	"github.com/pkg/errors"
)

func main() {
	awsConfig, err := awsConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	awsMailRegion, ok := os.LookupEnv("AWS_MAIL_REGION")
	if !ok {
		fmt.Println("No AWS_MAIL_REGION env var defined")
		os.Exit(1)
	}
	ses := mailer.SES{FromEmail: "test.ship@bitrise.io", Config: providers.AWSConfig{
		Region:          awsMailRegion,
		AccessKeyID:     awsConfig.AccessKeyID,
		SecretAccessKey: awsConfig.SecretAccessKey,
	}}
	targetEmail := os.Getenv("TARGET_EMAIL")
	if targetEmail == "" {
		fmt.Println("No TARGET_EMAIL env var defined")
		os.Exit(1)
	}
	emailName := os.Getenv("MAIL_TO_SEND")
	switch emailName {
	case "confirmation":
		err := ses.SendEmailConfirmation("Your test app", "http://here.you.can.confirm", &models.AppContact{
			Email: targetEmail,
			NotificationPreferencesData: json.RawMessage(`{}`),
			ConfirmationToken:           pointers.NewStringPtr("your-confirmation-token"),
		})
		if err != nil {
			failEmailSend(err)
		}
	case "new_version":
		err := ses.SendEmailNewVersion(&models.AppVersion{
			ArtifactInfoData: json.RawMessage(`{"version":"1.1.0"}`),
			BuildNumber:      "28",
			Platform:         "ios",
			App:              models.App{AppSlug: "test-app-slug-1"},
		}, []models.AppContact{models.AppContact{
			Email: targetEmail,
			NotificationPreferencesData: json.RawMessage(`{"new_version":true}`),
		}}, "http://bitrise.io",
			&bitrise.AppDetails{Title: "Standup Timer"})
		if err != nil {
			failEmailSend(err)
		}
	case "publish_succeeded":
		err := ses.SendEmailPublish(targetEmail, true)
		if err != nil {
			failEmailSend(err)
		}
	case "publish_failed":
		err := ses.SendEmailPublish(targetEmail, false)
		if err != nil {
			failEmailSend(err)
		}
	case "notifications":
		err := ses.SendEmailNotifications(targetEmail)
		if err != nil {
			failEmailSend(err)
		}
	default:
		failEmailSend(errors.New("No MAIL_TO_SEND env var defined"))
	}
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

func failEmailSend(err error) {
	fmt.Printf("Failed to send email: %s\n", err)
	os.Exit(1)
}

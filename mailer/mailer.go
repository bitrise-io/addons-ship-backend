package mailer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/templates"
)

// Interface ...
type Interface interface {
	SendEmailConfirmation(confirmURL string, contact *models.AppContact, appDetails *bitrise.AppDetails) error
	SendEmailNewVersion(appVersion *models.AppVersion, contacts []models.AppContact, frontendBaseURL string, appDetails *bitrise.AppDetails) error
	SendEmailPublish(appVersion *models.AppVersion, contacts []models.AppContact, appDetails *bitrise.AppDetails, frontendBaseURL string, publishSucceeded bool) error
}

// Request ...
type Request struct {
	From    string
	To      []string
	Subject string
}

// SESEmailInput ...
func (r *Request) SESEmailInput(template string, data map[string]interface{}) (*ses.SendEmailInput, error) {
	var toAddresses []*string
	for _, address := range r.To {
		toAddresses = append(toAddresses, aws.String(address))
	}
	body, err := templates.Get(template, data)
	if err != nil {
		return nil, err
	}
	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: toAddresses,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(r.Subject),
			},
		},
		Source: aws.String(r.From),
	}, nil
}

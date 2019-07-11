package mailer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/templates"
)

// Interface ...
type Interface interface {
	SendMail(r *Request, template string, data map[string]interface{}) error
	SendEmailConfirmation(appTitle, addonBaseURL string, contact *models.AppContact) error
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

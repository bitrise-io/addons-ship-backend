package mailer

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/providers"
	"github.com/pkg/errors"
)

// SES ...
type SES struct {
	FromEmail string
	Config    providers.AWSConfig
}

// SendMail ...
func (m *SES) SendMail(r *Request, template string, data map[string]interface{}) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(m.Config.Region),
		Credentials: credentials.NewStaticCredentials(
			m.Config.AccessKeyID,
			m.Config.SecretAccessKey,
			"",
		),
	})
	if err != nil {
		return err
	}
	svc := ses.New(sess)
	input, err := r.SESEmailInput(template, data)
	if err != nil {
		return err
	}
	_, err = svc.SendEmail(input)
	if err != nil {
		return err
	}

	return nil
}

// SendEmailConfirmation ...
func (m *SES) SendEmailConfirmation(appTitle, addonBaseURL string, contact *models.AppContact) error {
	notificationPreferences, err := contact.NotificationPreferences()
	if err != nil {
		return errors.WithStack(err)
	}
	nameForHey := strings.Split(contact.Email, "@")[0]

	return m.SendMail(&Request{
		To:      []string{contact.Email},
		From:    m.FromEmail,
		Subject: "Ship wants to send you notifications about the activity of this app. 🔔",
	},
		"email/confirmation.html",
		map[string]interface{}{
			"Name":              func() string { return nameForHey },
			"AppTitle":          func() string { return appTitle },
			"NewVersion":        func() bool { return notificationPreferences.NewVersion },
			"SuccessfulPublish": func() bool { return notificationPreferences.SuccessfulPublish },
			"FailedPublish":     func() bool { return notificationPreferences.FailedPublish },
			"URL": func() string {
				return fmt.Sprintf("%s/apps/%s/contacts/%s/confirm?token=%s", addonBaseURL, contact.App.AppSlug, contact.ID, contact.ConfirmationToken)
			},
		})
}

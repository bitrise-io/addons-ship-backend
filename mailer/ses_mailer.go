package mailer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/bitrise-io/api-utils/providers"
)

// SES ...
type SES struct {
	Config providers.AWSConfig
}

// SendMail ...
func (m *SES) SendMail(r *Request, template string, data map[string]interface{}) (bool, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(m.Config.Region),
		Credentials: credentials.NewStaticCredentials(
			m.Config.AccessKeyID,
			m.Config.SecretAccessKey,
			"",
		),
	})
	if err != nil {
		return false, err
	}
	svc := ses.New(sess)
	input, err := r.SESEmailInput(template, data)
	if err != nil {
		return false, err
	}
	_, err = svc.SendEmail(input)
	if err != nil {
		return false, err
	}

	return true, nil
}

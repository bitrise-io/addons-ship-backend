package mailer

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/api-utils/providers"
	"github.com/pkg/errors"
)

// SES ...
type SES struct {
	FromEmail string
	Config    providers.AWSConfig
}

func (m *SES) sendMail(r *Request, template string, data map[string]interface{}) error {
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
func (m *SES) SendEmailConfirmation(confirmURL string, contact *models.AppContact, appDetails *bitrise.AppDetails) error {
	appIconURL := "https://bitrise-public-content-production.s3.amazonaws.com/emails/invitation-app-32x32.png"
	if appDetails.AvatarURL != nil {
		appIconURL = *appDetails.AvatarURL
	}
	nameForHey := getUsernameFromEmail(contact.Email)
	var confirmationToken string
	if contact.ConfirmationToken != nil {
		confirmationToken = *contact.ConfirmationToken
	} else {
		return errors.New("Confirmation token is empty")
	}

	return m.sendMail(&Request{
		To:      []string{contact.Email},
		From:    m.FromEmail,
		Subject: "Ship wants to send you notifications about the activity of this app. 🔔",
	},
		"email/confirmation.html",
		map[string]interface{}{
			"CurrentTime": func() time.Time { return time.Now() },
			"Name":        func() string { return nameForHey },
			"AppTitle":    func() string { return appDetails.Title },
			"AppIconURL":  func() string { return appIconURL },
			"AppURL":      func() string { return fmt.Sprintf("%s?token=%s", confirmURL, confirmationToken) },
		},
	)
}

// SendEmailNewVersion ...
func (m *SES) SendEmailNewVersion(appVersion *models.AppVersion, contacts []models.AppContact, frontendBaseURL string, appDetails *bitrise.AppDetails) error {
	artifactInfo, err := appVersion.ArtifactInfo()
	if err != nil {
		return errors.WithStack(err)
	}
	appIconURL := "https://bitrise-public-content-production.s3.amazonaws.com/emails/invitation-app-32x32.png"
	if appDetails.AvatarURL != nil {
		appIconURL = *appDetails.AvatarURL
	}
	for _, contact := range contacts {
		notificationPreferences, err := contact.NotificationPreferences()
		if err != nil {
			return errors.WithStack(err)
		}
		if !notificationPreferences.NewVersion {
			return nil
		}
		nameForHey := getUsernameFromEmail(contact.Email)
		err = m.sendMail(&Request{
			To:      []string{contact.Email},
			From:    m.FromEmail,
			Subject: "New app version is available on Ship.",
		},
			"email/new_version.html",
			map[string]interface{}{
				"CurrentTime": func() time.Time { return time.Now() },
				"Name":        func() string { return nameForHey },
				"AppTitle":    func() string { return appDetails.Title },
				"AppIconURL":  func() string { return appIconURL },
				"NewVersion":  func() string { return artifactInfo.Version },
				"BuildNumber": func() string { return appVersion.BuildNumber },
				"AppPlatform": func() string { return appVersion.Platform },
				"AppURL":      func() string { return fmt.Sprintf("%s/apps/%s", frontendBaseURL, appVersion.App.AppSlug) },
			})
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// SendEmailPublish ...
func (m *SES) SendEmailPublish(appVersion *models.AppVersion, contacts []models.AppContact, appDetails *bitrise.AppDetails, frontendBaseURL string, publishSucceeded bool) error {
	artifactInfo, err := appVersion.ArtifactInfo()
	if err != nil {
		return errors.WithStack(err)
	}
	appIconURL := "https://bitrise-public-content-production.s3.amazonaws.com/emails/invitation-app-32x32.png"
	if appDetails.AvatarURL != nil {
		appIconURL = *appDetails.AvatarURL
	}
	var publishTarget string
	if appVersion.Platform == "ios" {
		publishTarget = "App Store Connect"
	} else if appVersion.Platform == "android" {
		publishTarget = "Play Store"
	}
	for _, contact := range contacts {
		notificationPreferences, err := contact.NotificationPreferences()
		if err != nil {
			return errors.WithStack(err)
		}
		if !notificationPreferences.NewVersion {
			return nil
		}
		nameForHey := getUsernameFromEmail(contact.Email)
		err = m.sendMail(&Request{
			To:      []string{contact.Email},
			From:    m.FromEmail,
			Subject: "App publish notification on Ship.",
		},
			"email/publish.html",
			map[string]interface{}{
				"CurrentTime":      func() time.Time { return time.Now() },
				"Name":             func() string { return nameForHey },
				"AppTitle":         func() string { return appDetails.Title },
				"AppIconURL":       func() string { return appIconURL },
				"Version":          func() string { return artifactInfo.Version },
				"BuildNumber":      func() string { return appVersion.BuildNumber },
				"AppPlatform":      func() string { return appVersion.Platform },
				"AppURL":           func() string { return fmt.Sprintf("%s/apps/%s", frontendBaseURL, appVersion.App.AppSlug) },
				"PublishSucceeded": func() bool { return publishSucceeded },
				"PublishURL":       func() string { return "https://bitrise.io" },
				"PublishTarget":    func() string { return publishTarget },
			})
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func getUsernameFromEmail(email string) string {
	return strings.Split(email, "@")[0]
}

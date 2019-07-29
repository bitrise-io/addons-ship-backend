package services_test

import (
	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/models"
)

type testMailer struct {
	sendEmailConfirmationFn func(confirmURL string, contact *models.AppContact, appDetails *bitrise.AppDetails) error
	sendEmailNewVersionFn   func(appVersion *models.AppVersion, contacts []models.AppContact, frontendBaseURL string, appDetails *bitrise.AppDetails) error
	sendEmailPublishFn      func(appVersion *models.AppVersion, contacts []models.AppContact, appDetails *bitrise.AppDetails, frontendBaseURL string, publishSucceeded bool) error
}

func (m *testMailer) SendEmailConfirmation(confirmURL string, contact *models.AppContact, appDetails *bitrise.AppDetails) error {
	if m.sendEmailConfirmationFn == nil {
		panic("You have to override Mailer.SendEmailConfirmation function in tests")
	}
	return m.sendEmailConfirmationFn(confirmURL, contact, appDetails)
}

func (m *testMailer) SendEmailNewVersion(appVersion *models.AppVersion, contacts []models.AppContact, frontendBaseURL string, appDetails *bitrise.AppDetails) error {
	if m.sendEmailNewVersionFn == nil {
		panic("You have to override Mailer.SendEmailNewVersion function in tests")
	}
	return m.sendEmailNewVersionFn(appVersion, contacts, frontendBaseURL, appDetails)
}

func (m *testMailer) SendEmailPublish(appVersion *models.AppVersion, contacts []models.AppContact, appDetails *bitrise.AppDetails, frontendBaseURL string, publishSucceeded bool) error {
	if m.sendEmailPublishFn == nil {
		panic("You have to override Mailer.SendEmailPublish function in tests")
	}
	return m.sendEmailPublishFn(appVersion, contacts, appDetails, frontendBaseURL, publishSucceeded)
}

package services_test

import (
	"github.com/bitrise-io/addons-ship-backend/bitrise"
	"github.com/bitrise-io/addons-ship-backend/models"
)

type testMailer struct {
	sendEmailConfirmationFn func(appTitle, addonBaseURL string, contact *models.AppContact) error
	sendEmailNewVersionFn   func(appVersion *models.AppVersion, contacts []models.AppContact, frontendBaseURL string, appDetails *bitrise.AppDetails) error
}

func (m *testMailer) SendEmailConfirmation(appTitle, addonBaseURL string, contact *models.AppContact) error {
	if m.sendEmailConfirmationFn == nil {
		panic("You have to override Mailer.SendEmailConfirmation function in tests")
	}
	return m.sendEmailConfirmationFn(appTitle, addonBaseURL, contact)
}

func (m *testMailer) SendEmailNewVersion(appVersion *models.AppVersion, contacts []models.AppContact, frontendBaseURL string, appDetails *bitrise.AppDetails) error {
	if m.sendEmailNewVersionFn == nil {
		panic("You have to override Mailer.SendEmailNewVersion function in tests")
	}
	return m.sendEmailNewVersionFn(appVersion, contacts, frontendBaseURL, appDetails)
}

package services_test

import "github.com/bitrise-io/addons-ship-backend/models"

type testMailer struct {
	sendEmailConfirmationFn func(appTitle, addonBaseURL string, contact *models.AppContact) error
}

func (m *testMailer) SendEmailConfirmation(appTitle, addonBaseURL string, contact *models.AppContact) error {
	if m.sendEmailConfirmationFn == nil {
		panic("You have to override Mailer.SendEmailConfirmation function in tests")
	}
	return m.sendEmailConfirmationFn(appTitle, addonBaseURL, contact)
}

package main

import (
	"fmt"
	"os"

	"github.com/bitrise-io/addons-ship-backend/mailer"
	"github.com/bitrise-io/api-utils/providers"
)

func main() {
	templateData := map[string]interface{}{
		"Name": func() string { return "Gergely Bekesi" },
		"URL":  func() string { return "http://www.bitrise.io" },
	}
	m := mailer.SES{Config: providers.AWSConfig{
		Region:          os.Getenv("AWS_MAIL_REGION"),
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}}

	_, err := m.SendMail(&mailer.Request{
		From:    "ship@bitrise.io",
		To:      []string{"gergely.bekesi@bitrise.io"},
		Subject: "Test",
	}, "mail.html", templateData)
	if err != nil {
		fmt.Println(err)
	}
}

package utils

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/afroash/mastupeti/initializers"
)

// func SendEmail() {
// 	// load env variables
// 	initializers.LoadEnvVariables()

// 	//setup authentication for sending email

// 	//construct the outgoing email.

// 	// send email and check for error.

// 	// if no error return success.

// }

func SendEmail(userEmail, subject, body, to string) error {
	//load env variables
	initializers.LoadEnvVariables()
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		os.Getenv("EMAIL_USERNAME"),
		os.Getenv("EMAIL_PASSWORD"),
		os.Getenv("EMAIL_HOST"),
	)

	// Construct the email message.
	msg := []byte("To: " + to + "\r\n" +
		"Subject: Message from Site\r\n" +
		"From: app_go@masterash.co.uk\r\n" +
		"\r\n" +
		"From: " + userEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Body: " + body + "\r\n")

	err := smtp.SendMail("email-smtp.eu-west-2.amazonaws.com:587", auth, "app_go@masterash.co.uk", []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

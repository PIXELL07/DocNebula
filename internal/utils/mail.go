package utils

import (
	"net/smtp"
)

func SendResetEmail(to string, token string) error {

	link := "https://DocNebula.app/reset-password?token=" + token

	msg := []byte(
		"Subject: Reset Password\n\n" +
			"Click this link to reset your password:\n" +
			link,
	)

	return smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", "your@email.com", "APP_PASSWORD", "smtp.gmail.com"),
		"your@email.com",
		[]string{to},
		msg,
	)
}

package mails

import "net/smtp"

const mailUsername = "mail.haha.ru@gmail.com"
const mailPassword = "qwerty152342007"

func SendMessage(message, to string) error {
	auth := smtp.PlainAuth(
		"",
		mailUsername,
		mailPassword,
		"smtp.gmail.com",
	)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		mailUsername,
		[]string{to},
		[]byte(mime + message),
	)
	if err != nil {
		return err
	}

	return nil
}

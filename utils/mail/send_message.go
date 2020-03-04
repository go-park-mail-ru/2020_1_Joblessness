package mail

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const key = "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"

func SendMessage(htmlContent, toMail string) error {
	from := mail.NewEmail("hh.ru", "maratishimbaev8@gmail.com")
	to := mail.NewEmail("User", toMail)

	subject := "Резюме"
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	client := sendgrid.NewSendClient(key)

	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}

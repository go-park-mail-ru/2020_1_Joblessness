package mail

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"os"
)

func SendMessage(htmlContent, toMail string) error {
	from := mail.NewEmail("hh.ru", "maratishimbaev8@gmail.com")
	to := mail.NewEmail("User", toMail)

	subject := "Резюме"
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_KEY"))

	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}

package mail

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"os"
)

func SendMessage(htmlContent, toMail string) error {
	from := mail.NewEmail("hh.ru", "maratishimbaev8@gmail.com")
	to := mail.NewEmail("User", toMail)

	subject := "Резюме"
	plainTextContent := "Резюме"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_KEY"))

	fmt.Println(os.Getenv("SENDGRID_KEY"))

	response, err := client.Send(message)
	if err != nil {
		return err
	}

	fmt.Println(response.StatusCode)
	fmt.Println(response.Body)
	fmt.Println(response.Headers)

	return nil
}

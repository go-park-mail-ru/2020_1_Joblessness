package mail

import (
	"bytes"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"joblessness/haha/models/base"
	"os"
)

type Mail struct {
	To string `json:"to"`
}

func SummaryToHTML(summary baseModels.Summary) (htmlString string, err error) {
	t := template.New("summary.html")

	t, err = t.ParseFiles("haha/utils/mail/summary.html")
	if err != nil {
		return htmlString, err
	}

	data := struct {
		Summary baseModels.Summary
	}{
		summary,
	}

	var html bytes.Buffer
	err = t.Execute(&html, data)
	if err != nil {
		return htmlString, err
	}

	return html.String(), err
}

func SendMessage(htmlContent, toMail string) (err error) {
	from := mail.NewEmail("hh.ru", "maratishimbaev8@gmail.com")
	to := mail.NewEmail("User", toMail)

	subject := "Резюме"
	plainTextContent := "Резюме"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_KEY"))

	_, err = client.Send(message)
	if err != nil {
		return err
	}

	return err
}

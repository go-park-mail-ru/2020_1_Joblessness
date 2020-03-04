package mail

import (
	"bytes"
	"html/template"
	"joblessness/haha/models"
)

const templatePath = "haha/utils/mail/templates/summary.html"

func SummaryToHTML(summary models.Summary) (string, error) {
	t := template.New("summary.html")

	t, err := t.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	data := struct {
		Summary models.Summary
	}{
		summary,
	}

	var html bytes.Buffer
	err = t.Execute(&html, data)
	if err != nil {
		return "", err
	}

	return html.String(), nil
}

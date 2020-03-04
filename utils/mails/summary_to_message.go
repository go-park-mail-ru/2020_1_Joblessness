package mails

import (
	_models "../../models"
	"bytes"
	"html/template"
)

const templatePath = "../utils/mails/templates/summary.html"

func SummaryToHTML(summary _models.Summary) (string, error) {
	t := template.New("summary.html")

	t, err := t.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	data := struct {
		Summary _models.Summary
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

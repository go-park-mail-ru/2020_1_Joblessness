package mail

import (
	"bytes"
	"html/template"
	"joblessness/haha/models"
)

func SummaryToHTML(t *template.Template, summary models.Summary) (string, error) {
	data := struct {
		Summary models.Summary
	}{
		summary,
	}

	var html bytes.Buffer
	err := t.Execute(&html, data)
	if err != nil {
		return "", err
	}

	return html.String(), nil
}

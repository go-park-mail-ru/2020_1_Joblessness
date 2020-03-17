package pdf

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"io"
	"joblessness/haha/models"
)

func SummaryToPdf(w io.Writer, summary models.Summary) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 16)
	tr := pdf.UnicodeTranslatorFromDescriptor("cp1252")

	name := fmt.Sprintf("Name: %s %s\n", summary.Author.FirstName, summary.Author.LastName)
	personal := fmt.Sprintf("Birthday: %s\nGender: %s\n", summary.Author.Birthday, summary.Author.Gender)
	contacts := fmt.Sprintf("Email: %s\nPhone: %s\n", summary.Author.Email, summary.Author.Phone)
	general := fmt.Sprintf("Educatiaon:\n %s\nExpirience:\n %s\n", summary.Educations, summary.Experiences)

	pdf.CellFormat(190, 7, "SUMMARY", "0", 0, "CM", false, 0, "")
	pdf.Ln(-1)
	pdf.MultiCell(100, 7, tr(name), "0", "LM", false)
	pdf.Ln(20)
	pdf.CellFormat(190, 7, "PERSONAL INFORMATION\n", "0", 0, "CM", false, 0, "")
	pdf.Ln(-1)
	pdf.MultiCell(100, 7, tr(personal), "0", "LM", false)
	pdf.Ln(20)
	pdf.CellFormat(190, 7, "CONTACTS", "0", 0, "CM", false, 0, "")
	pdf.Ln(-1)
	pdf.MultiCell(190, 7, tr(contacts), "0", "LM", false)
	pdf.Ln(20)
	pdf.CellFormat(190, 7, "GENERAL INFORMATION", "0", 0, "CM", false, 0, "")
	pdf.Ln(-1)
	pdf.MultiCell(190, 7, tr(general), "0", "LM", false)

	return pdf.Output(w)
}

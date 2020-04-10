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
	experience := experienceToStr(summary.Experiences)
	education := educationToStr(summary.Educations)

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
	pdf.MultiCell(190, 7, tr(experience), "0", "LM", false)
	pdf.Ln(-1)
	pdf.MultiCell(190, 7, tr(education), "0", "LM", false)

	return pdf.Output(w)
}

func experienceToStr(experience []models.Experience) (result string) {
	result = "Experience\n"
	for _, v := range experience {
		result += fmt.Sprintf("Company name: %s\nRole: %s\nResponsibilities: %s\nYears: %d-%d",
			v.CompanyName, v.Role, v.Responsibilities, v.Start.Year(), v.Stop.Year())
	}
	return result
}

func educationToStr(education []models.Education) (result string) {
	result = "Education\n"
	for _, v := range education {
		result += fmt.Sprintf("Institution: %s\nSpecialization: %s\nGraduated: %d\nType: %s\n",
			v.Institution, v.Speciality, v.Graduated.Year(), v.Type)
	}
	return result
}

package pdf

import (
	"fmt"
	"github.com/kpawlik/gofpdf"
	"io"
	"joblessness/haha/models"
)

func SummaryToPdf(w io.Writer, summary models.Summary) error {
	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.SetFontLocation("./haha/utils/pdf/")
	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 16)
	tr := pdf.UnicodeTranslatorFromDescriptor("cp1251")

	name := fmt.Sprintf("Name: %s %s\n", summary.Author.FirstName, summary.Author.LastName)
	personal := fmt.Sprintf("Birthday: %d-%d-%d\nGender: %s\n", summary.Author.Birthday.Year(),
		summary.Author.Birthday.Month(), summary.Author.Birthday.Day(), genderToStr(summary.Author.Gender))
	contacts := fmt.Sprintf("Email: %s\nPhone: %s\n", summary.Author.Email, summary.Author.Phone)
	experience := experienceToStr(summary.Experiences)
	education := educationToStr(summary.Educations)

	pdf.CellFormat(190, 7, "SUMMARY", "0", 0, "CM", false, 0, "")
	pdf.Ln(20)
	pdf.CellFormat(190, 7, "PERSONAL INFORMATION\n", "0", 0, "CM", false, 0, "")
	pdf.Ln(-1)
	pdf.MultiCell(100, 7, tr(name), "0", "LM", false)
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
		result += fmt.Sprintf("Company name: %s\nRole: %s\nResponsibilities: %s\nYears: %d-%d\n",
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

func genderToStr(gender string) string {
	switch gender {
	case "0":
		return "Мужчина"
	case "1":
		return "Женщина"
	case "2":
		return "Женщина в теле мужчины"
	case "3":
		return "Мужчина в теле мужчины"
	case "4":
		return "Женщина в теле женщины"
	case "5":
		return "Мужчина в теле женщины"
	case "6":
		return "Несколько мужчин"
	case "7":
		return "Мужчина и женщина"
	case "8":
		return "Женщина в теле мужчины"
	case "9":
		return "Все выше перечисленное"
	}
	return "Не определился"
}

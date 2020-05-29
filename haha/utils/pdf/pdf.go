package pdf

import (
	"fmt"
	"github.com/kpawlik/gofpdf"
	"io"
	"joblessness/haha/models/base"
	"path/filepath"
	"runtime"
)

var (
	_, file, _, _ = runtime.Caller(0)
	fontLocation  = filepath.Dir(file) + "/font/"
)

func SummaryToPdf(w io.Writer, summary *baseModels.Summary) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFontLocation(fontLocation)

	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 16)
	tr := pdf.UnicodeTranslatorFromDescriptor("cp1251")

	name := fmt.Sprintf("Имя: %s %s\n", summary.Author.FirstName, summary.Author.LastName)
	personal := fmt.Sprintf("День рождения: %d-%d-%d\nПол: %s\n", summary.Author.Birthday.Year(),
		summary.Author.Birthday.Month(), summary.Author.Birthday.Day(), genderToStr(summary.Author.Gender))
	contacts := fmt.Sprintf("Почта: %s\nТелефон: %s\n", summary.Author.Email, summary.Author.Phone)
	experience := experienceToStr(summary.Experiences)
	education := educationToStr(summary.Educations)

	pdf.CellFormat(190, 7, tr("РЕЗЮМЕ"), "0", 0, "CM", false, 0, "")
	pdf.Ln(20)
	pdf.CellFormat(190, 7, tr("ПЕРСОНАЛЬНАЯ ИНФОРМАЦИЯ\n"), "0", 0, "CM", false, 0, "")
	pdf.Ln(-1)
	pdf.MultiCell(100, 7, tr(name), "0", "LM", false)
	pdf.Ln(-1)
	pdf.MultiCell(100, 7, tr(personal), "0", "LM", false)
	pdf.Ln(20)
	pdf.CellFormat(190, 7, tr("КОНТАКТЫ"), "0", 0, "CM", false, 0, "")
	pdf.Ln(-1)
	pdf.MultiCell(190, 7, tr(contacts), "0", "LM", false)
	pdf.Ln(20)
	pdf.CellFormat(190, 7, tr("ОСНОВНАЯ ИНФОРМАЦИЯ"), "0", 0, "CM", false, 0, "")
	pdf.Ln(-1)
	pdf.MultiCell(190, 7, tr(experience), "0", "LM", false)
	pdf.Ln(-1)
	pdf.MultiCell(190, 7, tr(education), "0", "LM", false)

	return pdf.Output(w)
}

func experienceToStr(experience []baseModels.Experience) (result string) {
	if len(experience) == 0 {
		return ""
	}

	result = "ОПЫТ\n"
	for _, v := range experience {
		result += fmt.Sprintf("Компания: %s\nДолжность: %s\nОбязанности: %s\nПериод: %d-%d\n",
			v.CompanyName, v.Role, v.Responsibilities, v.Start.Year(), v.Stop.Year())
	}
	return result
}

func educationToStr(education []baseModels.Education) (result string) {
	if len(education) == 0 {
		return ""
	}

	result = "ОБРАЗОВАНИЕ\n"
	for _, v := range education {
		result += fmt.Sprintf("Институт: %s\nСпециальность: %s\nВыпуск: %d\nТип: %s\n",
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

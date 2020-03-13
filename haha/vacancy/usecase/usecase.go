package usecase

import (
	"joblessness/haha/models"
	"joblessness/haha/vacancy"
)

type VacancyUseCase struct {
	vacancyRepo vacancy.VacancyRepository
}

func NewVacancyUseCase(vacancyRepo vacancy.VacancyRepository) *VacancyUseCase {
	return &VacancyUseCase{vacancyRepo}
}

func (u *VacancyUseCase) CreateVacancy(vacancy models.Vacancy) (vacancyID uint64, err error) {
	return u.vacancyRepo.CreateVacancy(vacancy)
}

func (u *VacancyUseCase) GetVacancies() ([]models.Vacancy, error) {
	return u.vacancyRepo.GetVacancies()
}

func (u *VacancyUseCase) GetVacancy(vacancyID uint64) (models.Vacancy, error) {
	return u.vacancyRepo.GetVacancy(vacancyID)
}

func (u *VacancyUseCase) ChangeVacancy(vacancy models.Vacancy) error {
	return u.vacancyRepo.ChangeVacancy(vacancy)
}

func (u *VacancyUseCase) DeleteVacancy(vacancyID uint64) error {
	return u.vacancyRepo.DeleteVacancy(vacancyID)
}

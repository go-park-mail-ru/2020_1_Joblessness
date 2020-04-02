package vacancyUseCase

import (
	"joblessness/haha/models"
	"joblessness/haha/vacancy/interfaces"
	"strconv"
)

type VacancyUseCase struct {
	vacancyRepo vacancyInterfaces.VacancyRepository
}

func NewVacancyUseCase(vacancyRepo vacancyInterfaces.VacancyRepository) *VacancyUseCase {
	return &VacancyUseCase{vacancyRepo}
}

func (u *VacancyUseCase) CreateVacancy(vacancy *models.Vacancy) (vacancyID uint64, err error) {
	return u.vacancyRepo.CreateVacancy(vacancy)
}

func (u *VacancyUseCase) GetVacancies(page string) ([]models.Vacancy, error) {
	pageInt, _ := strconv.Atoi(page)
	return u.vacancyRepo.GetVacancies(pageInt)
}

func (u *VacancyUseCase) GetVacancy(vacancyID uint64) (*models.Vacancy, error) {
	return u.vacancyRepo.GetVacancy(vacancyID)
}

func (u *VacancyUseCase) ChangeVacancy(vacancy *models.Vacancy) error {
	return u.vacancyRepo.ChangeVacancy(vacancy)
}

func (u *VacancyUseCase) DeleteVacancy(vacancyID uint64) error {
	return u.vacancyRepo.DeleteVacancy(vacancyID)
}

func (u *VacancyUseCase) GetOrgVacancies(userID uint64) ([]models.Vacancy, error) {
	return u.vacancyRepo.GetOrgVacancies(userID)
}
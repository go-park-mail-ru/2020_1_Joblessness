package vacancyUseCase

import (
	"github.com/microcosm-cc/bluemonday"
	"joblessness/haha/models/base"
	"joblessness/haha/vacancy/interfaces"
	"strconv"
)

type VacancyUseCase struct {
	vacancyRepo vacancyInterfaces.VacancyRepository
	policy      *bluemonday.Policy
}

func NewVacancyUseCase(vacancyRepo vacancyInterfaces.VacancyRepository, policy *bluemonday.Policy) *VacancyUseCase {
	return &VacancyUseCase{
		vacancyRepo: vacancyRepo,
		policy:      policy,
	}
}

func (u *VacancyUseCase) CreateVacancy(vacancy *baseModels.Vacancy) (vacancyID uint64, err error) {
	return u.vacancyRepo.CreateVacancy(vacancy)
}

func (u *VacancyUseCase) GetVacancies(page string) (baseModels.Vacancies, error) {
	pageInt, _ := strconv.Atoi(page)
	res, err := u.vacancyRepo.GetVacancies(pageInt)
	if err != nil {
		return nil, err
	}

	res.Sanitize(u.policy)
	return res, nil
}

func (u *VacancyUseCase) GetVacancy(vacancyID uint64) (*baseModels.Vacancy, error) {
	res, err := u.vacancyRepo.GetVacancy(vacancyID)
	if err != nil {
		return nil, err
	}

	res.Sanitize(u.policy)
	return res, nil
}

func (u *VacancyUseCase) ChangeVacancy(vacancy *baseModels.Vacancy) (err error) {
	if err = u.vacancyRepo.CheckAuthor(vacancy.ID, vacancy.Organization.ID); err != nil {
		return err
	}

	return u.vacancyRepo.ChangeVacancy(vacancy)
}

func (u *VacancyUseCase) DeleteVacancy(vacancyID, authorID uint64) (err error) {
	if err = u.vacancyRepo.CheckAuthor(vacancyID, authorID); err != nil {
		return err
	}

	return u.vacancyRepo.DeleteVacancy(vacancyID)
}

func (u *VacancyUseCase) GetOrgVacancies(userID uint64) (baseModels.Vacancies, error) {
	res, err := u.vacancyRepo.GetOrgVacancies(userID)
	if err != nil {
		return nil, err
	}

	res.Sanitize(u.policy)
	return res, nil
}

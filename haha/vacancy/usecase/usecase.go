package vacancyUseCase

import (
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"joblessness/haha/models/base"
	"joblessness/haha/utils/chat"
	"joblessness/haha/vacancy/interfaces"
	"strconv"
	"time"
)

type VacancyUseCase struct {
	vacancyRepo vacancyInterfaces.VacancyRepository
	room chat.Room
	policy      *bluemonday.Policy
}

func NewVacancyUseCase(vacancyRepo vacancyInterfaces.VacancyRepository,
	room chat.Room,
	policy *bluemonday.Policy) *VacancyUseCase {
	return &VacancyUseCase{
		vacancyRepo: vacancyRepo,
		room: room,
		policy:      policy,
	}
}

func (u *VacancyUseCase) announceVacancy(vacancy *baseModels.Vacancy) (err error) {
	users, orgName, err := u.vacancyRepo.GetRelatedUsers(vacancy.Organization.ID)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("Похоже, у компании %s появилась новая вакансия %s, Вам это может быть интересно",
		orgName, vacancy.Name)

	for _, id := range users {
		u.room.SendGeneratedMessage(&chat.Message{
			Message:   message,
			UserOneId: vacancy.Organization.ID,
			UserOne:   "",
			UserTwoId: id,
			UserTwo:   "",
			Created:   time.Now(),
			VacancyID: vacancy.ID,
		})
	}

	return nil
}

func (u *VacancyUseCase) CreateVacancy(vacancy *baseModels.Vacancy) (vacancyID uint64, err error) {
	vacancyID, err = u.vacancyRepo.CreateVacancy(vacancy)
	if err != nil {
		return vacancyID, err
	}

	_ = u.announceVacancy(vacancy)
	return vacancyID, err
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

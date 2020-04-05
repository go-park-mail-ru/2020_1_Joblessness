package summaryUseCase

import (
	"joblessness/haha/models"
	"joblessness/haha/summary/interfaces"
	"strconv"
)

type SummaryUseCase struct {
	summaryRepo summaryInterfaces.SummaryRepository
}

func NewSummaryUseCase(summaryRepo summaryInterfaces.SummaryRepository) *SummaryUseCase {
	return &SummaryUseCase{summaryRepo}
}

func (u *SummaryUseCase) CreateSummary(summary *models.Summary) (summaryID uint64, err error) {
	return u.summaryRepo.CreateSummary(summary)
}

func (u *SummaryUseCase) GetAllSummaries(page string) (summaries []models.Summary, err error) {
	pageInt, _ := strconv.Atoi(page)
	return u.summaryRepo.GetAllSummaries(pageInt)
}

func (u *SummaryUseCase) GetUserSummaries(userID uint64) (summaries []models.Summary, err error) {
	return u.summaryRepo.GetUserSummaries(userID)
}

func (u *SummaryUseCase) GetSummary(summaryID uint64) (summary *models.Summary, err error) {
	return u.summaryRepo.GetSummary(summaryID)
}

func (u *SummaryUseCase) ChangeSummary(summary *models.Summary) (err error) {
	return u.summaryRepo.ChangeSummary(summary)
}

func (u *SummaryUseCase) DeleteSummary(summaryID uint64) (err error) {
	return u.summaryRepo.DeleteSummary(summaryID)
}

func (u *SummaryUseCase) SendSummary(sendSummary *models.SendSummary) (err error) {
	res, err := u.summaryRepo.IsPersonSummary(sendSummary.SummaryID, sendSummary.UserID)
	if err != nil {
		return err
	} else if !res {
		return summaryInterfaces.ErrPersonIsNotOwner
	}

	err = u.summaryRepo.SendSummary(sendSummary)
	if err == summaryInterfaces.ErrSummaryAlreadySend {
		err = u.summaryRepo.RefreshSummary(sendSummary.SummaryID, sendSummary.VacancyID)
	}

	return err
}

func (u *SummaryUseCase) ResponseSummary(sendSummary *models.SendSummary)  (err error) {
	res, err := u.summaryRepo.IsOrganizationVacancy(sendSummary.VacancyID, sendSummary.OrganizationID)
	if err != nil {
		return err
	} else if !res {
		return summaryInterfaces.ErrOrgIsNotOwner
	}
	err = u.summaryRepo.ResponseSummary(sendSummary)

	return err
}

func (u *SummaryUseCase) GetOrgSendSummaries(userID uint64) (summaries models.OrgSummaries, err error) {
	return u.summaryRepo.GetOrgSendSummaries(userID)
}

func (u *SummaryUseCase) GetUserSendSummaries(userID uint64) (summaries models.OrgSummaries, err error) {
	return u.summaryRepo.GetUserSendSummaries(userID)
}
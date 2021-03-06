package summaryUseCase

import (
	"errors"
	"github.com/microcosm-cc/bluemonday"
	"joblessness/haha/models/base"
	"joblessness/haha/summary/interfaces"
	"strconv"
)

type SummaryUseCase struct {
	summaryRepo summaryInterfaces.SummaryRepository
	policy      *bluemonday.Policy
}

func NewSummaryUseCase(summaryRepo summaryInterfaces.SummaryRepository, policy *bluemonday.Policy) *SummaryUseCase {
	return &SummaryUseCase{
		summaryRepo: summaryRepo,
		policy:      policy,
	}
}

func (u *SummaryUseCase) CreateSummary(summary *baseModels.Summary) (summaryID uint64, err error) {
	return u.summaryRepo.CreateSummary(summary)
}

func (u *SummaryUseCase) GetAllSummaries(page string) (summaries baseModels.Summaries, err error) {
	pageInt, _ := strconv.Atoi(page)
	res, err := u.summaryRepo.GetAllSummaries(pageInt)
	if err != nil {
		return nil, err
	}

	res.Sanitize(u.policy)
	return res, nil
}

func (u *SummaryUseCase) GetUserSummaries(page string, userID uint64) (summaries baseModels.Summaries, err error) {
	pageInt, _ := strconv.Atoi(page)
	res, err := u.summaryRepo.GetUserSummaries(pageInt, userID)
	if err != nil {
		return nil, err
	}

	res.Sanitize(u.policy)
	return res, nil
}

func (u *SummaryUseCase) GetSummary(summaryID uint64) (summary *baseModels.Summary, err error) {
	res, err := u.summaryRepo.GetSummary(summaryID)
	if err != nil {
		return nil, err
	}

	res.Sanitize(u.policy)
	return res, nil
}

func (u *SummaryUseCase) ChangeSummary(summary *baseModels.Summary) (err error) {
	if err = u.summaryRepo.CheckAuthor(summary.ID, summary.Author.ID); err != nil {
		return err
	}

	return u.summaryRepo.ChangeSummary(summary)
}

func (u *SummaryUseCase) DeleteSummary(summaryID uint64, authorID uint64) (err error) {
	if err = u.summaryRepo.CheckAuthor(summaryID, authorID); err != nil {
		return err
	}

	return u.summaryRepo.DeleteSummary(summaryID)
}

func (u *SummaryUseCase) SendSummary(sendSummary *baseModels.SendSummary) (err error) {
	if err := u.summaryRepo.CheckAuthor(sendSummary.SummaryID, sendSummary.UserID); err != nil {
		return err
	}

	err = u.summaryRepo.SendSummary(sendSummary)
	if errors.Is(err, summaryInterfaces.ErrSummaryAlreadySent) {
		err = u.summaryRepo.RefreshSummary(sendSummary.SummaryID, sendSummary.VacancyID)
	}

	return err
}

func (u *SummaryUseCase) GetOrgSendSummaries(userID uint64) (summaries baseModels.OrgSummaries, err error) {
	res, err := u.summaryRepo.GetOrgSendSummaries(userID)
	if err != nil {
		return nil, err
	}

	res.Sanitize(u.policy)
	return res, nil
}

func (u *SummaryUseCase) GetUserSendSummaries(userID uint64) (summaries baseModels.OrgSummaries, err error) {
	res, err := u.summaryRepo.GetUserSendSummaries(userID)
	if err != nil {
		return nil, err
	}

	res.Sanitize(u.policy)
	return res, nil
}

func (u *SummaryUseCase) SendSummaryByMail(summaryID, authorID uint64, to string) (err error) {
	if err = u.summaryRepo.CheckAuthor(summaryID, authorID); err != nil {
		return err
	}

	return u.summaryRepo.SendSummaryByMail(summaryID, to)
}

package recommendationUseCase

import (
	baseModels "joblessness/haha/models/base"
	recommendationInterfaces "joblessness/haha/recommendation/interfaces"
)

type UseCase struct {
	repository recommendationInterfaces.Repository
}

func NewUseCase(repository recommendationInterfaces.Repository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) GetRecommendedVacancies(userID uint64) (recommendations []baseModels.Vacancy, err error) {
	return u.repository.GetRecommendedVacancies(userID)
}

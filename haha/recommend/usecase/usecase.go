package recommendUseCase

import (
	"joblessness/haha/models/base"
	"joblessness/haha/recommend/interfaces"
)

type UseCase struct {
	repository recommendInterfaces.RecommendRepository
}

func NewUseCase(repository recommendInterfaces.RecommendRepository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) GetRecommendedVacancies(userID uint64) (recommendations []baseModels.Vacancy, err error) {
	return u.repository.GetRecommendedVacancies(userID)
}

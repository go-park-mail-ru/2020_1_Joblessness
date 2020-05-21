package recommendUseCase

import (
	"joblessness/haha/models/base"
	"joblessness/haha/recommend/interfaces"
)

const pageSize = 10

type UseCase struct {
	repository recommendInterfaces.RecommendRepository
}

func NewUseCase(repository recommendInterfaces.RecommendRepository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) GetRecommendedVacancies(userID, pageNumber uint64) (recommendations []baseModels.Vacancy, err error) {
	recommendations, recommendCount, err := u.repository.GetRecommendedVacancies(userID, pageSize, pageNumber*pageSize)
	if err != nil {
		return u.repository.GetPopularVacancies(pageSize, (pageNumber-1)*pageSize)
	}

	if len(recommendations) < pageSize {
		offset := (pageNumber-1)*pageSize - recommendCount
		vacancies, err := u.repository.GetPopularVacancies(uint64(pageSize-len(recommendations)), offset)
		if err != nil {
			return recommendations, err
		}

		recommendations = append(recommendations, vacancies...)
	}

	return recommendations, err
}

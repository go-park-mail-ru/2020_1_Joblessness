package recommendInterfaces

import "joblessness/haha/models/base"

type RecommendUseCase interface {
	GetRecommendedVacancies(userID, pageNumber uint64) (recommendations []baseModels.Vacancy, err error)
}

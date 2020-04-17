package recommendationInterfaces

import baseModels "joblessness/haha/models/base"

type UseCase interface {
	GetRecommendedVacancies(userID uint64) (recommendations []baseModels.Vacancy, err error)
}

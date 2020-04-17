package recommendationInterfaces

import "joblessness/haha/models/base"

type Repository interface {
	GetRecommendedVacancies(userID uint64) (recommendations []baseModels.Vacancy, err error)
}

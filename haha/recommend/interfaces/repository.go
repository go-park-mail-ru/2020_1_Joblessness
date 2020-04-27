package recommendInterfaces

import "joblessness/haha/models/base"

type RecommendRepository interface {
	GetRecommendedVacancies(userID uint64) (recommendations []baseModels.Vacancy, err error)
}

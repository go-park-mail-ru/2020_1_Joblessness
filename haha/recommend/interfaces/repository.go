package recommendInterfaces

import "joblessness/haha/models/base"

type RecommendRepository interface {
	GetPopularVacancies(limit, offset uint64) (vacancies []baseModels.Vacancy, err error)
	GetRecommendedVacancies(userID, limit, offset uint64) (recommendations []baseModels.Vacancy, recommendCount uint64, err error)
}

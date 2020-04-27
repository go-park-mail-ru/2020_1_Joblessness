package recommendPostgres

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/muesli/regommend"
	"joblessness/haha/models/base"
	"joblessness/haha/models/postgres"
	"joblessness/haha/recommend/interfaces"
	"joblessness/haha/vacancy/interfaces"
	"strconv"
	"strings"
)

type repository struct {
	db                *sql.DB
	vacancyRepository vacancyInterfaces.VacancyRepository
}

func NewRecommendRepository(db *sql.DB, vacancyRepository vacancyInterfaces.VacancyRepository) *repository {
	return &repository{
		db:                db,
		vacancyRepository: vacancyRepository,
	}
}

func (r *repository) GetRecommendedVacancies(userID uint64) (recommendations []baseModels.Vacancy, err error) {
	var hasUser bool
	checkUser := `SELECT COUNT(*) <> 0 FROM users WHERE id = $1`
	err = r.db.QueryRow(checkUser, userID).Scan(&hasUser)
	if err != nil {
		return recommendations, recommendInterfaces.ErrNoUser
	}

	getUsersWithResponses := `
		SELECT u.id, array_agg(r.vacancy_id)
		FROM users u
		JOIN summary s ON u.id = s.author
		JOIN response r ON s.id = r.summary_id
		GROUP BY u.id
		HAVING u.person_id IS NOT NULL`
	rows, err := r.db.Query(getUsersWithResponses)
	if err != nil {
		return recommendations, err
	}

	table := regommend.Table("recommendations")

	for rows.Next() {
		var userID uint64
		var vacanciesRaw string

		if err = rows.Scan(&userID, &vacanciesRaw); err != nil {
			return recommendations, err
		}

		vacanciesMap := make(map[interface{}]float64)

		vacanciesRaw = vacanciesRaw[1 : len(vacanciesRaw)-1]
		vacancies := strings.Split(vacanciesRaw, ",")

		for _, vacancy := range vacancies {
			vacanciesMap[vacancy] = 1.0
		}

		table.Add(int(userID), vacanciesMap)
	}

	table.Add(123, map[interface{}]float64{})

	recs, err := table.Recommend(int(userID))
	if err != nil {
		return recommendations, err
	}

	if len(recs) == 0 {
		return recommendations, recommendInterfaces.ErrNoRecommendation
	}

	recKeys := make([]int, len(recs))
	for i := range recs {
		key, _ := recs[i].Key.(string)
		recKeys[i], _ = strconv.Atoi(key)
		fmt.Printf("key: %s, distance: %f\n", recs[i].Key, recs[i].Distance)
	}

	getRecommendations := `
		SELECT id, organization_id, name, description, with_tax, responsibilities, conditions, keywords, salary_from, salary_to
		FROM vacancy
		WHERE id = ANY($1)`
	rows, err = r.db.Query(getRecommendations, pq.Array(recKeys))
	if err != nil {
		return recommendations, err
	}

	for rows.Next() {
		var vacancy pgModels.Vacancy

		err = rows.Scan(&vacancy.ID, &vacancy.OrganizationID, &vacancy.Name, &vacancy.Description, &vacancy.WithTax, &vacancy.Responsibilities, &vacancy.Conditions, &vacancy.Keywords, &vacancy.SalaryFrom, &vacancy.SalaryTo)
		if err != nil {
			return recommendations, err
		}
		user, organization, err := r.vacancyRepository.GetVacancyOrganization(vacancy.OrganizationID)
		if err != nil {
			return recommendations, err
		}

		recommendations = append(recommendations, *pgModels.ToBaseVacancy(&vacancy, user, organization))
	}

	return recommendations, err
}

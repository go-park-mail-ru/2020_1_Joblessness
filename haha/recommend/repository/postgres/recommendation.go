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

func (r *repository) GetPopularVacancies(limit, offset uint64) (vacancies []baseModels.Vacancy, err error) {
	getVacancies := `
		SELECT v.id, v.organization_id, v.name, v.description, v.with_tax, v.responsibilities, v.conditions, v.keywords, v.salary_from, v.salary_to, COUNT(*) count
		FROM vacancy v
		JOIN response r ON v.id = r.vacancy_id
		GROUP BY v.id
		ORDER BY count
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(getVacancies, limit, offset)
	if err != nil {
		return vacancies, err
	}
	defer rows.Close()

	for rows.Next() {
		var vacancy pgModels.Vacancy

		var count uint64
		err = rows.Scan(&vacancy.ID, &vacancy.OrganizationID, &vacancy.Name, &vacancy.Description, &vacancy.WithTax, &vacancy.Responsibilities, &vacancy.Conditions, &vacancy.Keywords, &vacancy.SalaryFrom, &vacancy.SalaryTo, &count)
		if err != nil {
			return vacancies, err
		}
		user, organization, err := r.vacancyRepository.GetVacancyOrganization(vacancy.OrganizationID)
		if err != nil {
			return vacancies, err
		}

		vacancies = append(vacancies, *pgModels.ToBaseVacancy(&vacancy, user, organization))
	}

	if len(vacancies) == 0 {
		return vacancies, recommendInterfaces.ErrNoRecommendation
	}

	return vacancies, err
}

func (r *repository) GetRecommendedVacancies(userID, limit, offset uint64) (recommends []baseModels.Vacancy, recommendCount uint64, err error) {
	var hasUser bool
	checkUser := `SELECT COUNT(*) <> 0 FROM users WHERE id = $1`
	err = r.db.QueryRow(checkUser, userID).Scan(&hasUser)
	if err != nil {
		return recommends, recommendCount, recommendInterfaces.ErrNoUser
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
		return recommends, recommendCount, err
	}
	defer rows.Close()

	table := regommend.Table("recommends")

	for rows.Next() {
		var userID uint64
		var vacanciesRaw string

		if err = rows.Scan(&userID, &vacanciesRaw); err != nil {
			return recommends, recommendCount, err
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
		return recommends, recommendCount, recommendInterfaces.ErrNoRecommendation
	}

	if len(recs) == 0 {
		return recommends, recommendCount, recommendInterfaces.ErrNoRecommendation
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
		WHERE id = ANY($1)
		LIMIT $2 OFFSET $3`
	rows, err = r.db.Query(getRecommendations, pq.Array(recKeys), limit, int(offset))
	if err != nil {
		return recommends, recommendCount, err
	}

	for rows.Next() {
		var vacancy pgModels.Vacancy

		err = rows.Scan(&vacancy.ID, &vacancy.OrganizationID, &vacancy.Name, &vacancy.Description, &vacancy.WithTax, &vacancy.Responsibilities, &vacancy.Conditions, &vacancy.Keywords, &vacancy.SalaryFrom, &vacancy.SalaryTo)
		if err != nil {
			return recommends, recommendCount, err
		}
		user, organization, err := r.vacancyRepository.GetVacancyOrganization(vacancy.OrganizationID)
		if err != nil {
			return recommends, recommendCount, err
		}

		recommends = append(recommends, *pgModels.ToBaseVacancy(&vacancy, user, organization))
		recommendCount++
	}

	return recommends, recommendCount, err
}

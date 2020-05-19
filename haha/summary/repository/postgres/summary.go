package summaryPostgres

import (
	"database/sql"
	"fmt"
	"joblessness/haha/models/base"
	pgModels "joblessness/haha/models/postgres"
	"joblessness/haha/summary/interfaces"
	"joblessness/haha/utils/mail"
)

type GetOptions struct {
	userID uint64
	page   int
}

type SummaryRepository struct {
	db *sql.DB
}

func NewSummaryRepository(db *sql.DB) *SummaryRepository {
	return &SummaryRepository{db}
}

func (r *SummaryRepository) CreateSummary(summary *baseModels.Summary) (summaryID uint64, err error) {
	summaryDB, educationDBs, experienceDBs := pgModels.ToPgSummary(summary)

	createSummary := `INSERT INTO summary (author, keywords, name, salary_from, salary_to)
					  VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	err = r.db.QueryRow(createSummary, summaryDB.AuthorID, summaryDB.Keywords, summaryDB.Name,
		summaryDB.SalaryFrom, summaryDB.SalaryTo).Scan(&summaryDB.ID)
	if err != nil {
		return summaryID, err
	}

	createEducation := `INSERT INTO education (summary_id, institution, speciality, graduated, type)
						VALUES ($1, $2, $3, $4, $5)`

	for _, educationDB := range educationDBs {
		_, err = r.db.Exec(createEducation, summaryDB.ID, educationDB.Institution, educationDB.Speciality,
			educationDB.Graduated, educationDB.Type)
		if err != nil {
			return summaryID, err
		}
	}

	createExperience := `INSERT INTO experience (summary_id, company_name, role, responsibilities, start, stop)
						 VALUES ($1, $2, $3, $4, $5, $6)`

	for _, experienceDB := range experienceDBs {
		_, err = r.db.Exec(createExperience, summaryDB.ID, experienceDB.CompanyName, experienceDB.Role,
			experienceDB.Responsibilities, experienceDB.Start, experienceDB.Stop)
		if err != nil {
			return summaryID, err
		}
	}

	return summaryDB.ID, nil
}

func (r *SummaryRepository) GetEducationsBySummaryID(summaryID uint64) ([]*pgModels.Education, error) {
	getEducations := `SELECT institution, speciality, graduated, type
					  FROM education WHERE summary_id = $1`

	rows, err := r.db.Query(getEducations, summaryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	educationDBs := make([]*pgModels.Education, 0)

	for rows.Next() {
		educationDB := pgModels.Education{SummaryID: summaryID}

		err = rows.Scan(&educationDB.Institution, &educationDB.Speciality, &educationDB.Graduated,
			&educationDB.Type)
		if err != nil {
			return nil, err
		}

		educationDBs = append(educationDBs, &educationDB)
	}

	return educationDBs, nil
}

func (r *SummaryRepository) GetExperiencesBySummaryID(summaryID uint64) ([]*pgModels.Experience, error) {
	getExperience := `SELECT company_name, role, responsibilities, start, stop
					  FROM experience WHERE summary_id = $1`

	rows, err := r.db.Query(getExperience, summaryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	experienceDBs := make([]*pgModels.Experience, 0)

	for rows.Next() {
		experienceDB := pgModels.Experience{SummaryID: summaryID}

		err = rows.Scan(&experienceDB.CompanyName, &experienceDB.Role, &experienceDB.Responsibilities,
			&experienceDB.Start, &experienceDB.Stop)
		if err != nil {
			return nil, err
		}

		experienceDBs = append(experienceDBs, &experienceDB)
	}

	return experienceDBs, nil
}

func (r *SummaryRepository) GetSummaryAuthor(authorID uint64) (*pgModels.User, *pgModels.Person, error) {
	user := pgModels.User{ID: authorID}
	var person pgModels.Person

	getUser := `SELECT tag, email, phone, avatar, name, surname, gender, birthday
				FROM users 
				JOIN person p on users.person_id = p.id
				WHERE users.id = $1`
	err := r.db.QueryRow(getUser, user.ID).Scan(&user.Tag, &user.Email, &user.Phone, &user.Avatar,
		&person.Name, &person.LastName, &person.Gender, &person.Birthday)

	return &user, &person, err
}

func (r *SummaryRepository) GetSummaries(opt *GetOptions) (baseModels.Summaries, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if opt.userID == 0 {
		getSummaries := `SELECT id, author, keywords, name, salary_from, salary_to
					 	 FROM summary
					 	 LIMIT $1 OFFSET $2;`
		rows, err = r.db.Query(getSummaries, 10, opt.page*10)
		if err != nil {
			return nil, err
		}
	} else {
		getSummaries := `SELECT id, author, keywords, name, salary_from, salary_to
					 	 FROM summary WHERE author = $1
						LIMIT $2 OFFSET $3;`
		rows, err = r.db.Query(getSummaries, opt.userID, 10, opt.page*10)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	summaries := make(baseModels.Summaries, 0)

	for rows.Next() {
		var summaryDB pgModels.Summary

		err = rows.Scan(&summaryDB.ID, &summaryDB.AuthorID, &summaryDB.Keywords, &summaryDB.Name, &summaryDB.SalaryFrom,
			&summaryDB.SalaryTo)
		if err != nil {
			return nil, err
		}

		educationDBs, err := r.GetEducationsBySummaryID(summaryDB.ID)
		if err != nil {
			return nil, err
		}

		experienceDBs, err := r.GetExperiencesBySummaryID(summaryDB.ID)
		if err != nil {
			return nil, err
		}

		userDB, personDB, err := r.GetSummaryAuthor(summaryDB.AuthorID)
		if err != nil {
			return nil, err
		}

		summaries = append(summaries, pgModels.ToBaseSummary(&summaryDB, educationDBs, experienceDBs, userDB, personDB))
	}

	return summaries, nil
}

func (r *SummaryRepository) GetAllSummaries(page int) (summaries baseModels.Summaries, err error) {
	return r.GetSummaries(&GetOptions{0, page})
}

func (r *SummaryRepository) GetUserSummaries(page int, userID uint64) (summaries baseModels.Summaries, err error) {
	return r.GetSummaries(&GetOptions{userID, page})
}

func (r *SummaryRepository) GetSummary(summaryID uint64) (*baseModels.Summary, error) {
	summaryDB := pgModels.Summary{ID: summaryID}

	getSummary := `SELECT author, keywords, name, salary_from, salary_to
				   FROM summary WHERE id = $1`
	err := r.db.QueryRow(getSummary, summaryID).Scan(&summaryDB.AuthorID, &summaryDB.Keywords, &summaryDB.Name,
		&summaryDB.SalaryFrom, &summaryDB.SalaryTo)
	if err != nil {
		return &baseModels.Summary{}, err
	}

	educationDBs, err := r.GetEducationsBySummaryID(summaryDB.ID)
	if err != nil {
		return &baseModels.Summary{}, err
	}

	experienceDBs, err := r.GetExperiencesBySummaryID(summaryDB.ID)
	if err != nil {
		return &baseModels.Summary{}, err
	}

	userDB, personDB, err := r.GetSummaryAuthor(summaryDB.AuthorID)

	return pgModels.ToBaseSummary(&summaryDB, educationDBs, experienceDBs, userDB, personDB), err
}

func (r *SummaryRepository) CheckAuthor(summaryID uint64, authorID uint64) (err error) {
	var isAuthor bool

	checkAuthor := `SELECT author = $1 FROM summary WHERE id = $2`
	if err = r.db.QueryRow(checkAuthor, authorID, summaryID).Scan(&isAuthor); err != nil {
		return err
	}

	if !isAuthor {
		return fmt.Errorf("%w, person id: %d, summary id: %d", summaryInterfaces.ErrPersonIsNotOwner, authorID, summaryID)
	}

	return err
}

func (r *SummaryRepository) ChangeSummary(summary *baseModels.Summary) (err error) {
	summaryDB, educationDBs, experienceDBs := pgModels.ToPgSummary(summary)

	changeSummary := `UPDATE summary
					  SET keywords = COALESCE(NULLIF($1, ''), keywords)
					  WHERE id = $2`
	_, err = r.db.Exec(changeSummary, summaryDB.Keywords, summaryDB.ID)
	if err != nil {
		return fmt.Errorf("%w, summary id: %d", summaryInterfaces.ErrSummaryNotFound, summaryDB.ID)
	}

	changeEducation := `UPDATE education
						SET institution = COALESCE(NULLIF($1, ''), institution), 
						    speciality = COALESCE(NULLIF($2, ''), speciality),
						    graduated = COALESCE(NULLIF($3, ''), graduated), 
						    type = COALESCE(NULLIF($4, ''), type)
						WHERE summary_id = $5`

	for _, educationDB := range educationDBs {
		_, err = r.db.Exec(changeEducation, &educationDB.Institution, &educationDB.Speciality, &educationDB.Graduated,
			&educationDB.Type, &educationDB.SummaryID)
		if err != nil {
			return err
		}
	}

	changeExperience := `UPDATE experience
						SET company_name = COALESCE(NULLIF($1, ''), company_name),
						    role = COALESCE(NULLIF($2, ''), role), 
						    responsibilities = COALESCE(NULLIF($3, ''), responsibilities),
						    start = COALESCE(NULLIF($4, ''), start),
						    stop = COALESCE(NULLIF($5, ''), stop)
						WHERE summary_id = $6`

	for _, experienceDB := range experienceDBs {
		_, err = r.db.Exec(changeExperience, &experienceDB.CompanyName, &experienceDB.Role,
			&experienceDB.Responsibilities, &experienceDB.Start, &experienceDB.Stop,
			&experienceDB.SummaryID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *SummaryRepository) DeleteSummary(summaryID uint64) (err error) {
	deleteSummary := `DELETE FROM summary
					  WHERE id = $1`
	_, err = r.db.Exec(deleteSummary, summaryID)
	if err != nil {
		return fmt.Errorf("%w, summary id: %d", summaryInterfaces.ErrSummaryNotFound, summaryID)
	}

	//TODO Убрал удаление связанных строк CASCADE есть в бд
	return nil
}

func (r *SummaryRepository) SendSummary(sendSummary *baseModels.SendSummary) (err error) {
	setLike := `INSERT INTO response (summary_id, vacancy_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING;`
	rows, err := r.db.Exec(setLike, sendSummary.SummaryID, sendSummary.VacancyID)
	if err != nil {
		return err
	}
	if rowsAf, _ := rows.RowsAffected(); rowsAf == 0 {
		return summaryInterfaces.ErrSummaryAlreadySent
	}

	return nil
}

func (r *SummaryRepository) RefreshSummary(summaryID, vacancyID uint64) (err error) {
	setLike := `UPDATE response 
				SET date = CURRENT_TIMESTAMP,
				    approved = false,
				    rejected = false
				WHERE summary_id = $1
				AND vacancy_id = $2;`
	rows, err := r.db.Exec(setLike, summaryID, vacancyID)
	if err != nil {
		return err
	}
	if rowsAf, _ := rows.RowsAffected(); rowsAf == 0 {
		return summaryInterfaces.ErrNoSummaryToRefresh
	}

	return nil
}

func (r *SummaryRepository) GetOrgSendSummaries(userID uint64) (summaries baseModels.OrgSummaries, err error) {
	getSummary := `SELECT u.id, u.tag, v.id, s.id, u.avatar, p.name, p.surname, r.approved, r.rejected, r.interview_date
				   FROM vacancy v 
				   JOIN response r on v.id = r.vacancy_id
				       AND r.approved = false
				       AND r.rejected = false
				   JOIN summary s on r.summary_id = s.id
				   JOIN users u on s.author = u.id
				   JOIN person p on u.person_id = p.id
				   WHERE v.organization_id = $1
				   order by r.date desc`

	rows, err := r.db.Query(getSummary, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summaries = make(baseModels.OrgSummaries, 0)

	for rows.Next() {
		var vacancyDB baseModels.VacancyResponse

		err = rows.Scan(&vacancyDB.UserID, &vacancyDB.Tag, &vacancyDB.VacancyID, &vacancyDB.SummaryID, &vacancyDB.Avatar,
			&vacancyDB.FirstName, &vacancyDB.LastName, &vacancyDB.Accepted, &vacancyDB.Denied, &vacancyDB.InterviewDate)
		if err != nil {
			return nil, err
		}

		summaries = append(summaries, &vacancyDB)
	}
	return summaries, nil
}

func (r *SummaryRepository) GetUserSendSummaries(userID uint64) (summaries baseModels.OrgSummaries, err error) {
	getSummary := `SELECT v.id, s.id, r.approved, r.rejected, r.interview_date
				   FROM vacancy v 
				   JOIN response r on v.id = r.vacancy_id
				       AND r.approved = false
				       AND r.rejected = false
				   JOIN summary s on r.summary_id = s.id
				   WHERE s.author = $1
				   order by r.date desc`

	rows, err := r.db.Query(getSummary, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summaries = make(baseModels.OrgSummaries, 0)

	for rows.Next() {
		var vacancyDB baseModels.VacancyResponse

		err = rows.Scan(&vacancyDB.VacancyID, &vacancyDB.SummaryID, &vacancyDB.Accepted, &vacancyDB.Denied,
			&vacancyDB.InterviewDate)
		if err != nil {
			return nil, err
		}

		summaries = append(summaries, &vacancyDB)
	}
	return summaries, nil
}

func (r *SummaryRepository) SendSummaryByMail(summaryID uint64, to string) (err error) {
	summary, err := r.GetSummary(summaryID)
	if err != nil {
		return err
	}

	htmlContent, err := mail.SummaryToHTML(*summary)
	if err != nil {
		return err
	}

	err = mail.SendMessage(htmlContent, to)
	if err != nil {
		return err
	}

	return err
}

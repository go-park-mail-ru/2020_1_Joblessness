package summaryRepoPostgres

import (
	"database/sql"
	"joblessness/haha/models"
	summaryInterfaces "joblessness/haha/summary/interfaces"
	"strings"
	"time"
)

type Summary struct {
	ID uint64
	AuthorID uint64
	Keywords string
	Name string
	SalaryFrom int
	SalaryTo int
}

type Education struct {
	SummaryID uint64
	Institution string
	Speciality string
	Graduated sql.NullTime
	Type string
}

type Experience struct {
	SummaryID uint64
	CompanyName string
	Role string
	Responsibilities string
	Start sql.NullTime
	Stop sql.NullTime
}

type User struct {
	ID             uint64
	Login          string
	Password       string
	OrganizationID uint64
	PersonID       uint64
	Tag            string
	Email          string
	Phone          string
	Registered     time.Time
	Avatar         string
}

type Person struct {
	ID       uint64
	Name     string
	Gender   string
	Birthday time.Time
}

func toPostgres(s *models.Summary) (summary *Summary, educations []*Education, experiences []*Experience) {
	summary = &Summary{
		ID:       s.ID,
		AuthorID: s.Author.ID,
		Keywords: s.Keywords,
		Name: s.Name,
		SalaryFrom: s.SalaryFrom,
		SalaryTo: s.SalaryTo,
	}

	for _, education := range s.Educations {
		educations = append(educations, &Education{
			SummaryID:   summary.ID,
			Institution: education.Institution,
			Speciality:  education.Speciality,
			Graduated:   sql.NullTime{education.Graduated, !education.Graduated.IsZero()},
			Type:        education.Type,
		})
	}

	for _, experience := range s.Experiences {
		experiences = append(experiences, &Experience{
			SummaryID:        summary.ID,
			CompanyName:      experience.CompanyName,
			Role:             experience.Role,
			Responsibilities: experience.Responsibilities,
			Start:            sql.NullTime{experience.Start, !experience.Start.IsZero()},
			Stop:             sql.NullTime{experience.Stop, !experience.Stop.IsZero()},
		})
	}

	return summary, educations, experiences
}

func toModel(s *Summary, eds []*Education, exs []*Experience, u *User, p *Person) *models.Summary {
	var educations []models.Education

	for _, ed := range eds {
		educations = append(educations, models.Education{
			Institution: ed.Institution,
			Speciality:  ed.Speciality,
			Graduated:   ed.Graduated.Time,
			Type:        ed.Type,
		})
	}

	var experiences []models.Experience

	for _, ex := range exs {
		experiences = append(experiences, models.Experience{
			CompanyName:      ex.CompanyName,
			Role:             ex.Role,
			Responsibilities: ex.Responsibilities,
			Start:            ex.Start.Time,
			Stop:             ex.Stop.Time,
		})
	}

	nameArr := strings.Split(p.Name, " ")
	firstName := nameArr[0]
	var lastName string
	if len(nameArr) > 1 {
		lastName = nameArr[1]
	}

	author := models.Author{
		ID:        u.ID,
		Tag:       u.Tag,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		FirstName: firstName,
		LastName:  lastName,
		Gender:    p.Gender,
		Birthday:  p.Birthday,
	}

	return &models.Summary{
		ID:          s.ID,
		Author:      author,
		Keywords:    s.Keywords,
		Name: s.Name,
		SalaryFrom: s.SalaryFrom,
		SalaryTo: s.SalaryTo,
		Educations:  educations,
		Experiences: experiences,
	}
}

type GetOptions struct {
	userID uint64
	page int
}

type SummaryRepository struct {
	db *sql.DB
}

func NewSummaryRepository(db *sql.DB) *SummaryRepository {
	return &SummaryRepository{db}
}

func (r *SummaryRepository) CreateSummary(summary *models.Summary) (summaryID uint64, err error) {
	summaryDB, educationDBs, experienceDBs := toPostgres(summary)

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

func (r *SummaryRepository) GetEducationsBySummaryID(summaryID uint64) ([]*Education, error) {
	getEducations := `SELECT institution, speciality, graduated, type
					  FROM education WHERE summary_id = $1`

	rows, err := r.db.Query(getEducations, summaryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	educationDBs := make([]*Education, 0)

	for rows.Next() {
		educationDB := Education{SummaryID: summaryID}

		err = rows.Scan(&educationDB.Institution, &educationDB.Speciality, &educationDB.Graduated,
			&educationDB.Type)
		if err != nil {
			return nil, err
		}

		educationDBs = append(educationDBs, &educationDB)
	}

	return educationDBs, nil
}

func (r *SummaryRepository) GetExperiencesBySummaryID(summaryID uint64) ([]*Experience, error) {
	getExperience := `SELECT company_name, role, responsibilities, start, stop
					  FROM experience WHERE summary_id = $1`

	rows, err := r.db.Query(getExperience, summaryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	experienceDBs := make([]*Experience, 0)

	for rows.Next() {
		experienceDB := Experience{SummaryID: summaryID}

		err = rows.Scan(&experienceDB.CompanyName, &experienceDB.Role, &experienceDB.Responsibilities,
			&experienceDB.Start, &experienceDB.Stop)
		if err != nil {
			return nil, err
		}

		experienceDBs = append(experienceDBs, &experienceDB)
	}

	return experienceDBs, nil
}

func (r *SummaryRepository) GetSummaryAuthor(authorID uint64) (*User, *Person, error) {
	user := User{ID: authorID}
	var person Person

	getUser := `SELECT tag, email, phone, avatar, name, gender, birthday
				FROM users 
				JOIN person p on users.person_id = p.id
				WHERE users.id = $1`
	err := r.db.QueryRow(getUser, user.ID).Scan(&user.Tag, &user.Email, &user.Phone, &user.Avatar,
		&person.Name, &person.Gender, &person.Birthday)

	return &user, &person, err
}

func (r *SummaryRepository) GetSummaries(opt *GetOptions) (models.Summaries, error) {
	var (
		rows *sql.Rows
		err error
	)

	if opt.userID == 0 {
		getSummaries := `SELECT id, author, keywords, name, salary_from, salary_to
					 	 FROM summary
					 	 LIMIT $1 OFFSET $2;`
		rows, err = r.db.Query(getSummaries, 10,  opt.page*10)
		if err != nil {
			return nil, err
		}
	} else {
		getSummaries := `SELECT id, author, keywords, name, salary_from, salary_to
					 	 FROM summary WHERE author = $1
						LIMIT $2 OFFSET $3;`
		rows, err = r.db.Query(getSummaries, opt.userID, 10,  opt.page*10)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	summaries := make(models.Summaries, 0)

	for rows.Next() {
		var summaryDB Summary

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

		summaries = append(summaries, toModel(&summaryDB, educationDBs, experienceDBs, userDB, personDB))
	}

	return summaries, nil
}

func (r *SummaryRepository) GetAllSummaries(page int) (summaries models.Summaries, err error) {
	return r.GetSummaries(&GetOptions{0, page})
}

func (r *SummaryRepository) GetUserSummaries(page int, userID uint64) (summaries models.Summaries, err error) {
	return r.GetSummaries(&GetOptions{userID, page})
}

func (r *SummaryRepository) GetSummary(summaryID uint64) (*models.Summary, error) {
	summaryDB := Summary{ID: summaryID}

	getSummary := `SELECT author, keywords, name, salary_from, salary_to
				   FROM summary WHERE id = $1`
	err := r.db.QueryRow(getSummary, summaryID).Scan(&summaryDB.AuthorID, &summaryDB.Keywords, &summaryDB.Name,
		&summaryDB.SalaryFrom, &summaryDB.SalaryTo)
	if err != nil {
		return &models.Summary{}, err
	}

	educationDBs, err := r.GetEducationsBySummaryID(summaryDB.ID)
	if err != nil {
		return &models.Summary{}, err
	}

	experienceDBs, err := r.GetExperiencesBySummaryID(summaryDB.ID)
	if err != nil {
		return &models.Summary{}, err
	}

	userDB, personDB, err := r.GetSummaryAuthor(summaryDB.AuthorID)

	return toModel(&summaryDB, educationDBs, experienceDBs, userDB, personDB), err
}

func (r *SummaryRepository) CheckAuthor(summaryID uint64, authorID uint64) (err error) {
	var isAuthor bool

	checkAuthor := `SELECT author = $1 FROM summary WHERE id = $2`
	if err = r.db.QueryRow(checkAuthor, authorID, summaryID).Scan(&isAuthor); err != nil {
		return err
	}

	if !isAuthor {
		return summaryInterfaces.NewErrorPersonIsNotOwner(authorID, summaryID)
	}

	return err
}

func (r *SummaryRepository) ChangeSummary(summary *models.Summary) (err error) {
	summaryDB, educationDBs, experienceDBs := toPostgres(summary)

	changeSummary := `UPDATE summary
					  SET keywords = COALESCE(NULLIF($1, ''), keywords)
					  WHERE id = $2`
	_, err = r.db.Exec(changeSummary, summaryDB.Keywords, summaryDB.ID)
	if err != nil {
		return summaryInterfaces.NewErrorSummaryNotFound(summaryDB.ID)
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
		return summaryInterfaces.NewErrorSummaryNotFound(summaryID)
	}

	//TODO Убрал удаление связанных строк CASCADE есть в бд
	return nil
}

func (r *SummaryRepository) SendSummary(sendSummary *models.SendSummary) (err error) {
	setLike := `INSERT INTO response (summary_id, vacancy_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING;`
	rows, err := r.db.Exec(setLike, sendSummary.SummaryID, sendSummary.VacancyID)
	if err != nil {
		return err
	}
	if rowsAf, _ := rows.RowsAffected(); rowsAf == 0 {
		return summaryInterfaces.NewErrorSummaryAlreadySent()
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
		return summaryInterfaces.NewErrorNoSummaryToRefresh()
	}

	return nil
}

func (r *SummaryRepository) IsOrganizationVacancy(vacancyID, userID uint64) (err error) {
	var isAuthor bool
	checkAuthor := `SELECT v.organization_id = $2
				FROM vacancy v 
				WHERE v.id = $1`
	if err = r.db.QueryRow(checkAuthor, vacancyID, userID).Scan(&isAuthor); err != nil {
		return err
	}

	if !isAuthor {
		return summaryInterfaces.NewErrorOrganizationIsNotOwner()
	}

	return err
}

func (r *SummaryRepository) ResponseSummary(sendSummary *models.SendSummary)  (err error) {
	response := `UPDATE response 
				SET date = CURRENT_TIMESTAMP,
				    approved = $1,
				    rejected = $2
				WHERE summary_id = $3
				AND vacancy_id = $4;`
	rows, err := r.db.Exec(response, sendSummary.Accepted, sendSummary.Denied, sendSummary.SummaryID, sendSummary.VacancyID)
	if err != nil {
		return err
	}
	if rowsAf, _ := rows.RowsAffected(); rowsAf == 0 {
		return summaryInterfaces.NewErrorNoSummaryToRefresh()
	}

	return nil
}

func (r *SummaryRepository) GetOrgSendSummaries(userID uint64) (summaries models.OrgSummaries, err error) {
	getSummary := `SELECT u.id, u.tag, v.id, s.id, s.keywords, s.name, v.name, r.approved, r.rejected
				   FROM vacancy v 
				   JOIN response r on v.id = r.vacancy_id
				       AND r.approved = false
				       AND r.rejected = false
				   JOIN summary s on r.summary_id = s.id
				   JOIN users u on s.author = u.id
				   WHERE v.organization_id = $1
				   order by r.date desc`

	rows, err := r.db.Query(getSummary, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summaries = make(models.OrgSummaries, 0)

	for rows.Next() {
		var vacancyDB models.VacancyResponse

		err = rows.Scan(&vacancyDB.UserID, &vacancyDB.Tag, &vacancyDB.VacancyID, &vacancyDB.SummaryID,
			&vacancyDB.Keywords, &vacancyDB.SummaryName, &vacancyDB.VacancyName, &vacancyDB.Accepted, &vacancyDB.Denied)
		if err != nil {
			return nil, err
		}

		summaries = append(summaries, &vacancyDB)
	}
	return summaries, nil
}

func (r *SummaryRepository) GetUserSendSummaries(userID uint64) (summaries models.OrgSummaries, err error) {
	getSummary := `SELECT v.id, s.id, s.keywords, s.name, v.name, r.approved, r.rejected
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

	summaries = make(models.OrgSummaries, 0)

	for rows.Next() {
		var vacancyDB models.VacancyResponse

		err = rows.Scan(&vacancyDB.VacancyID, &vacancyDB.SummaryID, &vacancyDB.Keywords, &vacancyDB.SummaryName,
			&vacancyDB.VacancyName, &vacancyDB.Accepted, &vacancyDB.Denied)
		if err != nil {
			return nil, err
		}

		summaries = append(summaries, &vacancyDB)
	}
	return summaries, nil
}

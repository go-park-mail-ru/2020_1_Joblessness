package postgres

import (
	"database/sql"
	"joblessness/haha/models"
)

type Summary struct {
	ID uint64
	AuthorID uint64
	Keywords string
}

type Education struct {
	SummaryID uint64
	Institution string
	Speciality string
	Graduated string
	Type string
}

type Experience struct {
	SummaryID uint64
	CompanyName string
	Role string
	Responsibilities string
	Start string
	Stop string
}

func toPostgres(s *models.Summary) (*Summary, *[]Education, *[]Experience) {
	return &Summary{
		ID:       s.ID,
		AuthorID: s.UserID,
		Keywords: "",
	},
	&[]Education{
		{
			SummaryID:   s.ID,
			Institution: "",
			Speciality:  "",
			Graduated:   "",
			Type:        "",
		},
	},
	&[]Experience{
		{
			SummaryID:        s.ID,
			CompanyName:      "",
			Role:             "",
			Responsibilities: "",
			Start:            "",
			Stop:             "",
		},
	}
}

func toModel(s *Summary, ed *[]Education, ex *[]Experience) *models.Summary {
	return &models.Summary{
		ID:          s.ID,
		UserID:      s.AuthorID,
		FirstName:   "",
		LastName:    "",
		PhoneNumber: "",
		Email:       "",
		BirthDate:   "",
		Gender:      "",
		Experience:  "",
		Education:   "",
	}
}

type GetOptions struct {
	userID uint64
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) CreateSummary(summary *models.Summary) (summaryID uint64, err error) {
	summaryDB, educationDBs, experienceDBs := toPostgres(summary)

	createSummary := `INSERT INTO summary (author, keywords)
					  VALUES ($1, $2) RETURNING id;`
	err = r.db.QueryRow(createSummary, summaryDB.AuthorID, summaryDB.Keywords).Scan(&summaryDB.ID)
	if err != nil {
		return summaryID, err
	}

	createEducation := `INSERT INTO education (summary_id, institution, speciality, graduated, type)
						VALUES ($1, $2, $3, $4, $5)`

	for _, educationDB := range *educationDBs {
		_, err = r.db.Exec(createEducation, educationDB.SummaryID, educationDB.Institution, educationDB.Speciality,
						   educationDB.Graduated, educationDB.Type)
		if err != nil {
			return summaryID, err
		}
	}

	createExperience := `INSERT INTO experience (summary_id, company_name, role, responsibilities, start, stop)
						 VALUES ($1, $2, $3, $4, $5, $6)`

	for _, experienceDB := range *experienceDBs {
		_, err = r.db.Exec(createExperience, experienceDB.SummaryID, experienceDB.CompanyName, experienceDB.Role,
						   experienceDB.Responsibilities, experienceDB.Start, experienceDB.Stop)
		if err != nil {
			return summaryID, err
		}
	}

	return summaryID, nil
}

func (r *Repository) GetEducationsBySummaryID(summaryID uint64) (*[]Education, error) {
	getEducations := `SELECT (institution, speciality, graduated, type)
					  FROM education WHERE summary_id = $1`

	rows, err := r.db.Query(getEducations, summaryID)
	if err != nil {
		return &[]Education{}, err
	}

	var educationDBs []Education

	for rows.Next() {
		educationDB := Education{SummaryID: summaryID}

		err = rows.Scan(&educationDB.Institution, &educationDB.Speciality, &educationDB.Graduated,
			&educationDB.Type)
		if err != nil {
			return &[]Education{}, err
		}

		educationDBs = append(educationDBs, educationDB)
	}

	return &educationDBs, nil
}

func (r *Repository) GetExperiencesBySummaryID(summaryID uint64) (*[]Experience, error) {
	getExperience := `SELECT (company_name, role, responsibilities, start, stop)
					  FROM experience WHERE summary_id = $1`

	rows, err := r.db.Query(getExperience, summaryID)
	if err != nil {
		return &[]Experience{}, err
	}

	var experienceDBs []Experience

	for rows.Next() {
		experienceDB := Experience{SummaryID: summaryID}

		err = rows.Scan(&experienceDB.CompanyName, &experienceDB.Role, &experienceDB.Responsibilities,
			&experienceDB.Start, &experienceDB.Stop)
		if err != nil {
			return &[]Experience{}, err
		}

		experienceDBs = append(experienceDBs, experienceDB)
	}

	return &experienceDBs, nil
}

func (r *Repository) GetSummaries(opt *GetOptions) (*[]models.Summary, error) {
	var rows *sql.Rows
	var err error

	if opt.userID == 0 {
		getSummaries := `SELECT (id, author, keywords)
					 	FROM summary;`
		rows, err = r.db.Query(getSummaries)
		if err != nil {
			return &[]models.Summary{}, err
		}
	} else {
		getSummaries := `SELECT (id, author, keywords)
					 	FROM summary WHERE author = $1;`
		rows, err = r.db.Query(getSummaries, opt.userID)
		if err != nil {
			return &[]models.Summary{}, err
		}
	}

	var summaries []models.Summary

	for rows.Next() {
		var summaryDB Summary

		err = rows.Scan(&summaryDB.ID, &summaryDB.AuthorID, &summaryDB.Keywords)
		if err != nil {
			return &[]models.Summary{}, err
		}

		educationDBs, err := r.GetEducationsBySummaryID(summaryDB.ID)
		if err != nil {
			return &[]models.Summary{}, err
		}

		experienceDBs, err := r.GetExperiencesBySummaryID(summaryDB.ID)
		if err != nil {
			return &[]models.Summary{}, err
		}

		summaries = append(summaries, *toModel(&summaryDB, educationDBs, experienceDBs))
	}

	return &summaries, nil
}

func (r *Repository) GetAllSummaries() (summaries *[]models.Summary, err error) {
	return r.GetSummaries(&GetOptions{})
}

func (r *Repository) GetUserSummaries(userID uint64) (summaries *[]models.Summary, err error) {
	return r.GetSummaries(&GetOptions{userID})
}

func (r *Repository) GetSummary(summaryID uint64) (*models.Summary, error) {
	var summaryDB Summary

	getSummary := `SELECT (author, keywords)
				   FROM summary WHERE id = $1`
	err := r.db.QueryRow(getSummary, summaryID).Scan(&summaryDB)
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

	return toModel(&summaryDB, educationDBs, experienceDBs), nil
}

func (r *Repository) ChangeSummary(summary *models.Summary) (err error) {
	summaryDB, educationDBs, experienceDBs := toPostgres(summary)

	changeSummary := `UPDATE summary
					  SET author = $1, keywords = $2
					  WHERE id = $3`
	_, err = r.db.Exec(changeSummary, summaryDB.AuthorID, summaryDB.Keywords, summaryDB.ID)
	if err != nil {
		return err
	}

	changeEducation := `UPDATE education
						SET institution = $1, speciality = $2, graduated = $3, type = $4
						WHERE summary_id = $5`

	for _, educationDB := range *educationDBs {
		_, err = r.db.Exec(changeEducation, &educationDB.Institution, &educationDB.Speciality, &educationDB.Graduated,
						   &educationDB.Type, educationDB.SummaryID)
		if err != nil {
			return err
		}
	}

	changeExperience := `UPDATE experience
						SET company_name = $1, role = $2, responsibilities = $3, start = $4, stop = $5
						WHERE summary_id = $6`

	for _, experienceDB := range *experienceDBs {
		_, err = r.db.Exec(changeExperience, &experienceDB.CompanyName, &experienceDB.Role,
						   &experienceDB.Responsibilities, &experienceDB.Start, &experienceDB.Stop,
						   &experienceDB.SummaryID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) DeleteSummary(summaryID uint64) (err error) {
	deleteSummary := `DELETE FROM summary
					  WHERE id = $1`
	_, err = r.db.Exec(deleteSummary, summaryID)
	if err != nil {
		return err
	}

	deleteEducations := `DELETE FROM education
						 WHERE summary_id = $1`
	_, err = r.db.Exec(deleteEducations, summaryID)
	if err != nil {
		return err
	}

	deleteExperiences := `DELETE FROM experience
						  WHERE summary_id = $1`
	_, err = r.db.Exec(deleteExperiences, summaryID)
	if err != nil {
		return err
	}

	return nil
}

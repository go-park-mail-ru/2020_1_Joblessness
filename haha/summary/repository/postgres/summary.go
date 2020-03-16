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

func (r *Repository) GetSummaries() (summaries *[]models.Summary, err error) {
	getSummaries := `SELECT (id, author, keywords)
					 FROM summary;`

	summaryRows, err := r.db.Query(getSummaries)
	if err != nil {
		return summaries, err
	}

	getEducations := `SELECT (institution, speciality, graduated, type)
					  FROM education WHERE summary_id = $1`

	getExperience := `SELECT (company_name, role, responsibilities, start, stop)
					  FROM experience WHERE summary_id = $1`

	for summaryRows.Next() {
		var summaryDB Summary

		err = summaryRows.Scan(&summaryDB.ID, &summaryDB.AuthorID, &summaryDB.Keywords)
		if err != nil {
			return summaries, err
		}

		educationRows, err := r.db.Query(getEducations)
		if err != nil {
			return summaries, err
		}

		var educationDBs []Education

		for educationRows.Next() {
			educationDB := Education{SummaryID: summaryDB.ID}

			err = educationRows.Scan(&educationDB.Institution, &educationDB.Speciality, &educationDB.Graduated,
									 &educationDB.Type)
			if err != nil {
				return summaries, err
			}

			educationDBs = append(educationDBs, educationDB)
		}

		experienceRows, err := r.db.Query(getExperience)
		if err != nil {
			return summaries, err
		}

		var experienceDBs []Experience

		for experienceRows.Next() {
			experienceDB := Experience{SummaryID: summaryDB.ID}

			err = experienceRows.Scan(&experienceDB.CompanyName, &experienceDB.Role, &experienceDB.Responsibilities,
									  &experienceDB.Start, &experienceDB.Stop)
			if err != nil {
				return summaries, err
			}

			experienceDBs = append(experienceDBs, experienceDB)
		}

		*summaries = append(*summaries, *toModel(&summaryDB, &educationDBs, &experienceDBs))
	}

	return summaries, nil
}

func (r *Repository) GetSummary(summaryID uint64) (summary *models.Summary, err error) {
	return &models.Summary{}, nil
}

func (r *Repository) ChangeSummary(summary *models.Summary) (err error) {
	return nil
}

func (r *Repository) DeleteSummary(summaryID uint64) (err error) {
	return nil
}

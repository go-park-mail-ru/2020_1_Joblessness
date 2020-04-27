package userPostgres

import (
	"database/sql"
	"errors"
	"joblessness/haha/models/base"
	pgModels "joblessness/haha/models/postgres"
	"joblessness/haha/user/interfaces"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) SaveAvatarLink(link string, userID uint64) (err error) {
	if link == "" {
		return errors.New("no link to save")
	}

	insertUser := `UPDATE users SET avatar = $1 WHERE id = $2;`
	_, err = r.db.Exec(insertUser, link, userID)

	return err
}

func (r *UserRepository) GetPerson(userID uint64) (*baseModels.Person, error) {
	user := pgModels.User{ID: userID}

	getUser := "SELECT login, COALESCE(person_id, 0), email, phone, avatar, tag FROM users WHERE id = $1;"
	err := r.db.QueryRow(getUser, userID).
		Scan(&user.Login, &user.PersonID, &user.Email, &user.Phone, &user.Avatar, &user.Tag)
	if err != nil {
		return nil, err
	}

	if user.PersonID == 0 {
		return nil, userInterfaces.NewErrorUserNotPerson(userID)
	}

	var person pgModels.Person

	getPerson := "SELECT name, surname, gender, birthday FROM person WHERE id = $1;"
	err = r.db.QueryRow(getPerson, user.PersonID).Scan(&person.Name, &person.LastName, &person.Gender, &person.Birthday)
	if err != nil {
		return nil, err
	}

	return pgModels.ToBasePerson(&user, &person), nil
}

func (r *UserRepository) changeUser(user *pgModels.User) error {
	changeUser := `UPDATE users 
					SET password = COALESCE(NULLIF($1, ''), password), 
					    tag = COALESCE(NULLIF($2, ''), tag), 
					    email = COALESCE(NULLIF($3, ''), email),
					    phone = COALESCE(NULLIF($4, ''), phone)
					WHERE id = $5;`
	_, err := r.db.Exec(changeUser, user.Password, user.Tag, user.Email, user.Phone, user.ID)

	return err
}

func (r *UserRepository) ChangePerson(p *baseModels.Person) error {
	user, dbPerson := pgModels.ToPgPerson(p)

	getUser := "SELECT person_id FROM users WHERE id = $1;"
	err := r.db.QueryRow(getUser, user.ID).Scan(&user.PersonID)
	if err != nil {
		return userInterfaces.NewErrorUserNotPerson(user.ID)
	}

	var birthday sql.NullTime
	if p.Birthday.IsZero() {
		birthday.Valid = false
	} else {
		birthday.Time = p.Birthday
		birthday.Valid = true
	}

	changePerson := `UPDATE person 
					SET name = COALESCE(NULLIF($1, ''), name), 
					    surname = COALESCE(NULLIF($2, ''), surname), 
					    gender = COALESCE(NULLIF($3, ''), gender), 
					    birthday = COALESCE($4, birthday) 
					WHERE id = $5;`
	_, err = r.db.Exec(changePerson, dbPerson.Name, dbPerson.LastName, dbPerson.Gender, birthday, user.PersonID)
	if err != nil {
		return err
	}
	err = r.changeUser(user)

	return err
}

func (r *UserRepository) GetOrganization(userID uint64) (*baseModels.Organization, error) {
	user := pgModels.User{ID: userID}

	getUser := "SELECT login, COALESCE(organization_id, 0), email, phone, avatar, tag FROM users WHERE id = $1;"
	err := r.db.QueryRow(getUser, userID).
		Scan(&user.Login, &user.OrganizationID, &user.Email, &user.Phone, &user.Avatar, &user.Tag)
	if err != nil {
		return nil, err
	}

	if user.OrganizationID == 0 {
		return nil, userInterfaces.NewErrorUserNotOrganization(userID)
	}

	var org pgModels.Organization

	getOrg := "SELECT name, site, about FROM organization WHERE id = $1;"
	err = r.db.QueryRow(getOrg, user.OrganizationID).Scan(&org.Name, &org.Site, &org.About)
	if err != nil {
		return nil, err
	}

	return pgModels.ToBaseOrganization(&user, &org), nil
}

func (r *UserRepository) ChangeOrganization(o *baseModels.Organization) error {
	user, dbOrg := pgModels.ToPgOrganization(o)

	getUser := "SELECT organization_id FROM users WHERE id = $1;"
	err := r.db.QueryRow(getUser, user.ID).Scan(&user.OrganizationID)
	if err != nil {
		return userInterfaces.NewErrorUserNotOrganization(user.ID)
	}

	changeOrg := `UPDATE organization 
					SET name = COALESCE(NULLIF($1, ''), name),
					    site = COALESCE(NULLIF($2, ''), site),
					    about = COALESCE(NULLIF($3, ''), about)
					WHERE id = $4;`
	_, err = r.db.Exec(changeOrg, dbOrg.Name, dbOrg.Site, dbOrg.About, user.OrganizationID)
	if err != nil {
		return err
	}
	err = r.changeUser(user)

	return err
}

func (r *UserRepository) GetListOfOrgs(page int) (result baseModels.Organizations, err error) {
	getOrgs := `SELECT users.id as userId, name, site
				FROM users, organization
				WHERE users.organization_id = organization.id
				ORDER BY registered desc
				LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(getOrgs, 10, page*10)

	if err != nil {
		return result, err
	}
	defer rows.Close()

	var (
		userId     uint64
		name, site string
	)
	result = make(baseModels.Organizations, 0)

	for rows.Next() {
		err := rows.Scan(&userId, &name, &site)
		if err != nil {
			return result, err
		}

		result = append(result, &baseModels.Organization{
			ID:         userId,
			Login:      "",
			Password:   "",
			Name:       name,
			Site:       site,
			Email:      "",
			Phone:      "",
			Tag:        "",
			Registered: time.Time{},
		})
	}

	return result, rows.Err()
}

func (r *UserRepository) SetOrDeleteLike(userID, favoriteID uint64) (res bool, err error) {
	setLike := `INSERT INTO favorite (user_id, favorite_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING;`
	rows, err := r.db.Exec(setLike, userID, favoriteID)
	if err != nil {
		return false, userInterfaces.NewErrorUserNotFound(userID)
	}
	if rowsAf, err := rows.RowsAffected(); rowsAf != 0 {
		return true, err
	}

	deleteLike := `DELETE FROM favorite 
				WHERE user_id = $1 AND favorite_id = $2;`
	_, err = r.db.Exec(deleteLike, userID, favoriteID)

	return false, err
}

func (r *UserRepository) LikeExists(userID, favoriteID uint64) (res bool, err error) {
	setLike := `SELECT count(*)
				FROM favorite f
				WHERE f.user_id = $1
				AND f.favorite_id = $2;`
	rows, err := r.db.Query(setLike, userID, favoriteID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (r *UserRepository) GetUserFavorite(userID uint64) (res baseModels.Favorites, err error) {
	getFavorite := `SELECT u.id, u.tag, u.avatar, u.person_id, 'o.name', 'p.name', 'p.surname'
				FROM favorite f, users u 
				WHERE f.favorite_id = u.id
				AND f.user_id = $1;`
	rows, err := r.db.Query(getFavorite, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		personID sql.NullInt64
		personName sql.NullString
		personSurname sql.NullString
		orgName sql.NullString
	)
	res = make(baseModels.Favorites, 0)

	for rows.Next() {
		var favorite baseModels.Favorite
		err := rows.Scan(&favorite.ID, &favorite.Tag, &favorite.Avatar, &personID, &orgName, &personName, &personSurname)
		if err != nil {
			return nil, err
		}

		if personID.Valid {
			favorite.IsPerson = true
			favorite.Name = personName.String
			favorite.Surname = personSurname.String
		} else {
			favorite.IsPerson = false
			favorite.Name = orgName.String
		}
		res = append(res, &favorite)
	}

	return res, err
}

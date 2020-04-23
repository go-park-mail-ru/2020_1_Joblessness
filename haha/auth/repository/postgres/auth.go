package authPostgres

import (
	"database/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"joblessness/haha/utils/salt"
	"time"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r AuthRepository) CreateUser(login, password string, personID, organizationID uint64) (err error) {
	hashedPassword, err := salt.HashAndSalt(password)

	insertUser := `INSERT INTO users (login, password, organization_id, person_id) 
					VALUES(NULLIF($1, ''), NULLIF($2, ''), NULLIF($3, 0), NULLIF($4, 0))`
	_, err = r.db.Exec(insertUser, login, hashedPassword, organizationID, personID)

	return err
}

func (r *AuthRepository) RegisterPerson(login, password, name string) (err error) {
	var personID uint64

	err = r.db.QueryRow("INSERT INTO person (name) VALUES($1) RETURNING id", name).
		Scan(&personID)
	if err != nil {
		return err
	}
	//TODO исполнять как единая транзация
	err = r.CreateUser(login, password, personID, 0)

	return err
}

func (r *AuthRepository) RegisterOrganization(login, password, name string) (err error) {
	var organizationID uint64

	err = r.db.QueryRow("INSERT INTO organization (name) VALUES($1) RETURNING id", name).
		Scan(&organizationID)
	if err != nil {
		return err
	}
	//TODO исполнять как единая транзация
	err = r.CreateUser(login, password, 0, organizationID)

	return err
}

func (r *AuthRepository) Login(login, password, SID string) (userId uint64, err error) {
	//TODO user_id, session_id уникальные

	checkUser := "SELECT id, password FROM users WHERE login = $1"
	var hashedPwd string
	rows := r.db.QueryRow(checkUser, login)
	err = rows.Scan(&userId, &hashedPwd)
	if err != nil || !salt.ComparePasswords(hashedPwd, password) {
		return 0, status.Error(codes.InvalidArgument, "wrong login or password")
	}

	insertSession := `INSERT INTO session (user_id, session_id, expires) 
					VALUES($1, $2, $3)`
	_, err = r.db.Exec(insertSession, userId, SID, time.Now().Add(time.Hour))

	return userId, err
}

func (r *AuthRepository) Logout(sessionId string) (err error) {
	//TODO user_id, session_id уникальные

	deleteRow := "DELETE FROM session WHERE session_id = $1;"
	_, err = r.db.Exec(deleteRow, sessionId)

	return err
}

func (r *AuthRepository) SessionExists(sessionId string) (userId uint64, err error) {
	//TODO session_id - pk, возвращать тип сессии

	checkUser := "SELECT user_id, expires FROM session WHERE session_id = $1;"
	var expires time.Time
	err = r.db.QueryRow(checkUser, sessionId).Scan(&userId, &expires)
	if err != nil {
		return 0, status.Error(500, "wrong sid")
	}

	if expires.Before(time.Now()) {
		deleteRow := "DELETE FROM session WHERE session_id = $1;"
		_, err = r.db.Exec(deleteRow, sessionId)
		if err != nil {
			return 0, err
		}
		userId = 0
		return userId, status.Error(500, "wrong sid")
	}

	return userId, err
}

func (r *AuthRepository) DoesUserExists(login string) (err error) {
	var columnCount int
	checkUser := "SELECT count(*) FROM users WHERE login = $1"
	err = r.db.QueryRow(checkUser, login).Scan(&columnCount)

	if err != nil {
		return err
	}

	if columnCount != 0 {
		return status.Error(codes.AlreadyExists, "user already exists")
	}
	return nil
}

func (r *AuthRepository) GetRole(userID uint64) (string, error) {
	var personID, organizationID sql.NullInt64
	checkUser := "SELECT person_id, organization_id FROM users WHERE id = $1;"
	err := r.db.QueryRow(checkUser, userID).Scan(&personID, &organizationID)
	if err != nil {
		return "", err
	}

	if personID.Valid {
		return "person", nil
	} else if organizationID.Valid {
		return "organization", nil
	}
	return "", nil
}

package authHttp

import (
	"database/sql"
	"encoding/json"
	"github.com/juju/loggo"
	"github.com/kataras/golog"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"joblessness/haha/auth/interfaces"
	"joblessness/haha/models/base"
	"net/http"
	"time"
)

type Handler struct {
	useCase authInterfaces.AuthUseCase
	logger  *loggo.Logger
}

func NewHandler(useCase authInterfaces.AuthUseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) RegisterPerson(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	var user baseModels.Person
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = validator.New().Struct(user); err != nil {
		golog.Errorf("#%s: %s", rID, "Empty login or password")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.RegisterPerson(&user)
	switch true {
	case errors.Is(err, authInterfaces.ErrUserAlreadyExists):
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	case err == nil:
		golog.Infof("#%s: %s",  rID, "Успешно")
		w.WriteHeader(http.StatusCreated)
	default:
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	}
}

func (h *Handler) RegisterOrg(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	var org baseModels.Organization
	err := json.NewDecoder(r.Body).Decode(&org)
	if err != nil {
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = validator.New().Struct(org); err != nil {
		golog.Errorf("#%s: %s", rID, "Empty login or password")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.RegisterOrganization(&org)
	switch true {
	case errors.Is(err, authInterfaces.ErrUserAlreadyExists):
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	case err == nil:
		golog.Infof("#%s: %s",  rID, "success")
		w.WriteHeader(http.StatusCreated)
	default:
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	var user baseModels.UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	golog.Debugf("#%s: %w", rID, user)
	if err != nil {
		golog.Errorf("#%s: %w\n%w", rID, err, user)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = validator.New().Struct(user); err != nil {
		golog.Errorf("#%s: %s", rID, "Empty login or password")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, role, sessionId, err := h.useCase.Login(user.Login, user.Password)
	switch true {
	case errors.Is(err, authInterfaces.ErrWrongLoginOrPassword):
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	case err == nil:
		cookie := &http.Cookie {
			Name: "session_id",
			Value: sessionId,
			Expires: time.Now().Add(time.Hour),
			MaxAge: 100000,
			Path: "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, cookie)

		jsonData, _ := json.Marshal(baseModels.ResponseRole{ID: userId, Role: role})
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonData)
	default:
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = h.useCase.Logout(session.Value)
	if err != nil {
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	session.Path = "/"
	http.SetCookie(w, session)

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, _ := r.Context().Value("userID").(uint64)

	role, err := h.useCase.GetRole(userID)
	switch true {
	case errors.Is(err, sql.ErrNoRows):
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
		json, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	case err == nil:
		jsonData, _ := json.Marshal(baseModels.ResponseRole{ID: userID, Role: role})
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonData)
	default:
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(json)
	}
}

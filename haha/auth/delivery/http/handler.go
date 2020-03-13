package httpAuth

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"joblessness/haha/auth"
	"joblessness/haha/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	useCase auth.UseCase
}

func NewHandler(useCase auth.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type inputPerson struct {
	ID uint `json:"id,omitempty"`
	Login string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`

	FirstName string `json:"first-name,omitempty"`
	LastName string `json:"last-name,omitempty"`
	Email string `json:"email,omitempty"`
	PhoneNumber string `json:"phone-number,omitempty"`
}

type UserLogin struct {
	Login string `json:"login"`
	Password string `json:"password"`
}

type ResponseId struct {
	ID int `json:"id"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user models.Person
	err = json.Unmarshal(body, &user)
	log.Println("user recieved: ", user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.Login == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.RegisterPerson(user.Login, user.Password, user.FirstName, user.LastName, user.Email, user.PhoneNumber)
	switch err {
	case auth.ErrUserAlreadyExists:
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	case nil:
		log.Println("Success")
	default:
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user models.UserLogin
	err = json.Unmarshal(body, &user)
	log.Println("user recieved: ", user)
	if err != nil {
		log.Println("Unmarshal went wrong")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.Login == "" || user.Password == "" {
		log.Println("login or password empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, sessionId, err := h.useCase.Login(user.Login, user.Password)
	switch err {
	case auth.ErrWrongLogPas:
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	case nil:
		log.Println("Success")
	default:
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	jsonData, _ := json.Marshal(models.Response{userId})
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = h.useCase.Logout(session.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	session.Path = "/"
	http.SetCookie(w, session)

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId, err := h.useCase.SessionExists(session.Value)
	switch err {
	case auth.ErrWrongSID:
		log.Println(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	case nil:
		log.Println("Success")
	default:
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, _ := json.Marshal(models.Response{userId})
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (h *Handler) GetPerson(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)

	user, err := h.useCase.GetPerson(userID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type Response struct {
		User models.UserInfo `json:"user"`
		Summaries []models.UserSummary `json:"summaries"`
	}

	jsonData, _ := json.Marshal(Response{
		models.UserInfo{user.FirstName, user.LastName, "", ""},
		[]models.UserSummary{},
	})
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) ChangePerson(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	log.Println("session cookie: ", session)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := h.useCase.SessionExists(session.Value)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if reqID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64); reqID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var person models.Person

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &person)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	person.ID = userID

	err = h.useCase.ChangePerson(person)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

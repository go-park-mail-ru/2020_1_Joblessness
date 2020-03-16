package httpAuth

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
	"github.com/kataras/golog"
	"io/ioutil"
	"joblessness/haha/auth"
	"joblessness/haha/models"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	useCase auth.AuthUseCase
	logger  *loggo.Logger
}

func NewHandler(useCase auth.AuthUseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type UserLogin struct {
	Login string `json:"login"`
	Password string `json:"password"`
}

type ResponseId struct {
	ID int `json:"id"`
}

func (h *Handler) SetAvatar(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		golog.Errorf("#%s: %s",  rID, "no cookie")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if reqID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64); reqID != userID {
		golog.Errorf("#%s: %s",  rID, "user requested and session user doesnt match")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := r.ParseMultipartForm(1024 * 1024 * 5) //5mb
	if err != nil {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	form := r.MultipartForm

	//TODO перенести в юзкейс
	err = h.useCase.SetAvatar(form, userID)

	switch err {
	case auth.ErrUploadAvatar:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusFailedDependency)
		return
	case nil:
		golog.Infof("#%s: %s",  rID, "Успешно")
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) RegisterPerson(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	body, err := ioutil.ReadAll(r.Body)
	golog.Debugf("#%s: %s",  rID, body)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user models.Person
	err = json.Unmarshal(body, &user)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.Login == "" || user.Password == "" {
		golog.Errorf("#%s: %s",  rID, "Empty login or password")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.RegisterPerson(&user)
	switch err {
	case auth.ErrUserAlreadyExists:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	case nil:
		golog.Infof("#%s: %s",  rID, "Успешно")
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) RegisterOrg(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	golog.Debugf("#%s: %s",  rID, body)

	var org models.Organization
	err = json.Unmarshal(body, &org)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if org.Login == "" || org.Password == "" {
		golog.Errorf("#%s: %s",  rID, "Empty login or password")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.RegisterOrganization(&org)
	switch err {
	case auth.ErrUserAlreadyExists:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	case nil:
		golog.Infof("#%s: %s",  rID, "success")
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user models.UserLogin
	err = json.Unmarshal(body, &user)
	golog.Debugf("#%s: %s",  rID, user)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.Login == "" || user.Password == "" {
		golog.Errorf("#%s: %s",  rID, "login or password in empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, sessionId, err := h.useCase.Login(user.Login, user.Password)
	switch err {
	case auth.ErrWrongLogPas:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	case nil:
		golog.Infof("#%s: %s",  rID, "success")
	default:
		golog.Errorf("#%s: %w",  rID, err)
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
	rID := r.Context().Value("rID").(string)

	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = h.useCase.Logout(session.Value)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
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

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		golog.Errorf("#%s: %s",  rID, "no cookie")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jsonData, _ := json.Marshal(models.Response{userID})
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (h *Handler) GetPerson(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	userID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)
	user, err := h.useCase.GetPerson(userID)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
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
	rID := r.Context().Value("rID").(string)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		golog.Errorf("#%s: %s",  rID, "no cookie")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if reqID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64); reqID != userID {
		golog.Errorf("#%s: %s",  rID, "user requested and session user doesnt match")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var person models.Person

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &person)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	person.ID = userID

	err = h.useCase.ChangePerson(person)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	userID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)

	user, err := h.useCase.GetOrganization(userID)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type Response struct {
		User models.OrganizationInfo `json:"user"`
	}

	jsonData, _ := json.Marshal(Response{
		models.OrganizationInfo{user.Name, user.Site, "", ""},
	})
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) ChangeOrganization(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		golog.Errorf("#%s: %s",  rID, "no cookie")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if reqID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64); reqID != userID {
		golog.Errorf("#%s: %s",  rID, "user requested and session user doesnt match")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var org models.Organization

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &org)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	org.ID = userID

	err = h.useCase.ChangeOrganization(org)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetListOfOrgs(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	//TODO проверять существование контекста

	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	listOrgs, err := h.useCase.GetListOfOrgs(page)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type Response struct {
		Organizations []models.Organization `json:"organizations"`
	}

	jsonData, _ := json.Marshal(Response{
		listOrgs,
	})

	w.WriteHeader(http.StatusNoContent)
	w.Write(jsonData)
}
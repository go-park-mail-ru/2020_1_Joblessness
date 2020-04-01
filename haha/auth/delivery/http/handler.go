package httpAuth

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
	"github.com/kataras/golog"
	"joblessness/haha/auth/interfaces"
	"joblessness/haha/models"
	"net/http"
	"strconv"
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
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	form := r.MultipartForm

	//TODO перенести в юзкейс
	err = h.useCase.SetAvatar(form, userID)

	switch err {
	case authInterfaces.ErrUploadAvatar:
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

	var user models.Person
	err := json.NewDecoder(r.Body).Decode(&user)
	golog.Debugf("#%s: %w", rID, user)
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
	case authInterfaces.ErrUserAlreadyExists:
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

	var org models.Organization
	err := json.NewDecoder(r.Body).Decode(&org)
	golog.Debugf("#%s: %w", rID, org)
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
	case authInterfaces.ErrUserAlreadyExists:
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

	var user models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	golog.Debugf("#%s: %w", rID, user)
	if err != nil {
		golog.Errorf("#%s: %w\n%w",  rID, err, user)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.Login == "" || user.Password == "" {
		golog.Errorf("#%s: %s",  rID, "login or password in empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, role, sessionId, err := h.useCase.Login(user.Login, user.Password)
	switch err {
	case authInterfaces.ErrWrongLogPas:
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

	jsonData, _ := json.Marshal(models.Response{ID: userId, Role: role})
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

	role, err := h.useCase.GetRole(userID)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, _ := json.Marshal(models.Response{ID: userID, Role: role})
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

	jsonData, _ := json.Marshal(user)
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
	err := json.NewDecoder(r.Body).Decode(&person)
	golog.Debugf("#%s: %w", rID, person)
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
		User models.Organization `json:"user"`
	}

	jsonData, _ := json.Marshal(Response{*user})
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
	err := json.NewDecoder(r.Body).Decode(&org)
	golog.Debugf("#%s: %w", rID, org)
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

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) LikeUser(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		golog.Errorf("#%s: %s",  rID, "no cookie")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var favoriteID uint64
	favoriteID, _ = strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)

	likeSet, err := h.useCase.LikeUser(userID, favoriteID)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type Response struct {
		Like bool `json:"like"`
	}
	jsonData, _ := json.Marshal(Response{
		likeSet,
	})
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) GetUserFavorite(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		golog.Errorf("#%s: %s",  rID, "no cookie")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if favoriteID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64); favoriteID != userID {
		golog.Errorf("#%s: %s",  rID, "user requested and session user doesnt match")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	favorites, err := h.useCase.GetUserFavorite(userID)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, _ := json.Marshal(favorites)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
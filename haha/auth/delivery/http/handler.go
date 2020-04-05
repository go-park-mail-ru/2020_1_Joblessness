package httpAuth

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
	"github.com/kataras/golog"
	"gopkg.in/go-playground/validator.v9"
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

func (h *Handler) SetAvatar(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, _ := r.Context().Value("userID").(uint64)

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

	err = h.useCase.SetAvatar(form, userID)

	switch err {
	case authInterfaces.ErrUploadAvatar:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusFailedDependency)
	case nil:
		golog.Infof("#%s: %s",  rID, "Успешно")
		w.WriteHeader(http.StatusCreated)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Handler) RegisterPerson(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	var user models.Person
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = validator.New().Struct(user); err != nil {
		golog.Errorf("#%s: %s",  rID, "Empty login or password")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.RegisterPerson(&user)
	switch err {
	case authInterfaces.ErrUserAlreadyExists:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
	case nil:
		golog.Infof("#%s: %s",  rID, "Успешно")
		w.WriteHeader(http.StatusCreated)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) RegisterOrg(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	var org models.Organization
	err := json.NewDecoder(r.Body).Decode(&org)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = validator.New().Struct(org); err != nil {
		golog.Errorf("#%s: %s",  rID, "Empty login or password")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.RegisterOrganization(&org)
	switch err {
	case authInterfaces.ErrUserAlreadyExists:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
	case nil:
		golog.Infof("#%s: %s",  rID, "success")
		w.WriteHeader(http.StatusCreated)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	var user models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	golog.Debugf("#%s: %w", rID, user)
	if err != nil {
		golog.Errorf("#%s: %w\n%w",  rID, err, user)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = validator.New().Struct(user); err != nil {
		golog.Errorf("#%s: %s",  rID, "Empty login or password")
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

	jsonData, _ := json.Marshal(models.ResponseRole{ID: userId, Role: role})
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
	userID, _ := r.Context().Value("userID").(uint64)

	role, err := h.useCase.GetRole(userID)
	switch err {
	case sql.ErrNoRows :
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		jsonData, _ := json.Marshal(models.ResponseRole{ID: userID, Role: role})
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonData)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) GetPerson(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)

	user, err := h.useCase.GetPerson(userID)
	switch err {
	case authInterfaces.ErrUserNotPerson :
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		jsonData, _ := json.Marshal(user)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) ChangePerson(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, _ := r.Context().Value("userID").(uint64)

	if reqID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64); reqID != userID {
		golog.Errorf("#%s: %s",  rID, "user requested and session user doesnt match")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var person models.Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	person.ID = userID
	err = h.useCase.ChangePerson(person)
	switch err {
	case authInterfaces.ErrUserNotPerson :
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)

	user, err := h.useCase.GetOrganization(userID)
	switch err {
	case authInterfaces.ErrUserNotOrg :
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		jsonData, _ := json.Marshal(user)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) ChangeOrganization(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, _ := r.Context().Value("userID").(uint64)

	if reqID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64); reqID != userID {
		golog.Errorf("#%s: %s",  rID, "user requested and session user doesnt match")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var org models.Organization
	err := json.NewDecoder(r.Body).Decode(&org)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org.ID = userID
	err = h.useCase.ChangeOrganization(org)
	switch err {
	case authInterfaces.ErrUserNotOrg :
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) GetListOfOrgs(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	page, err := strconv.Atoi(r.FormValue("page"))

	listOrgs, err := h.useCase.GetListOfOrgs(page)
	switch err {
	case authInterfaces.ErrUserNotFound :
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		jsonData, _ := json.Marshal(listOrgs)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) LikeUser(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, _ := r.Context().Value("userID").(uint64)

	var (
		favoriteID uint64
		response models.ResponseBool
		err error
	)
	favoriteID, _ = strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)

	response.Like, err = h.useCase.LikeUser(userID, favoriteID)
	switch err {
	case authInterfaces.ErrUserNotFound :
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		jsonData, _ := json.Marshal(response)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) LikeExists(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, _ := r.Context().Value("userID").(uint64)

	var (
		favoriteID uint64
		response models.ResponseBool
		err error
	)
	favoriteID, _ = strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)

	response.Like, err = h.useCase.LikeExists(userID, favoriteID)
	switch err {
	case nil:
		jsonData, _ := json.Marshal(response)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) GetUserFavorite(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, _ := r.Context().Value("userID").(uint64)

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
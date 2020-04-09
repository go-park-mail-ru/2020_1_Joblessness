package httpUser

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
	"github.com/kataras/golog"
	"joblessness/haha/models"
	"joblessness/haha/user/interfaces"
	"net/http"
	"strconv"
)

type Handler struct {
	useCase userInterfaces.UserUseCase
	logger  *loggo.Logger
}

func NewHandler(useCase userInterfaces.UserUseCase) *Handler {
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

	switch err.(type) {
	case *userInterfaces.ErrorUploadAvatar:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusFailedDependency)
		json, _ := json.Marshal(models.Error{Message: err.Error()})
		w.Write(json)
	case nil:
		golog.Infof("#%s: %s",  rID, "Успешно")
		w.WriteHeader(http.StatusCreated)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(models.Error{Message: err.Error()})
		w.Write(json)
	}
}

func (h *Handler) GetPerson(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)

	user, err := h.useCase.GetPerson(userID)
	switch err.(type) {
	case *userInterfaces.ErrorUserNotPerson :
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
		json, _ := json.Marshal(models.Error{Message: err.Error()})
		w.Write(json)
	case nil:
		jsonData, _ := json.Marshal(user)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(models.Error{Message: err.Error()})
		w.Write(json)
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

	var person *models.Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	person.ID = userID
	err = h.useCase.ChangePerson(person)
	switch err.(type) {
	case *userInterfaces.ErrorUserNotPerson :
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
		json, _ := json.Marshal(models.Error{Message: err.Error()})
		w.Write(json)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(models.Error{Message: err.Error()})
		w.Write(json)
	}
}

func (h *Handler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, _ := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)

	user, err := h.useCase.GetOrganization(userID)
	switch err.(type) {
	case *userInterfaces.ErrorUserNotOrganization :
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
		json, _ := json.Marshal(models.Error{Message: err.Error()})
		w.Write(json)
	case nil:
		jsonData, _ := json.Marshal(user)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(models.Error{Message: err.Error()})
		w.Write(json)
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

	var org *models.Organization
	err := json.NewDecoder(r.Body).Decode(&org)
	if err != nil {
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org.ID = userID
	err = h.useCase.ChangeOrganization(org)
	switch err.(type) {
	case *userInterfaces.ErrorUserNotOrganization:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
		json, _ := json.Marshal(models.Error{Message: err.Error()})
		w.Write(json)
	case nil:
		w.WriteHeader(http.StatusNoContent)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(models.Error{Message: err.Error()})
		w.Write(json)
	}
}

func (h *Handler) GetListOfOrgs(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	page, err := strconv.Atoi(r.FormValue("page"))

	listOrgs, err := h.useCase.GetListOfOrgs(page)
	switch err.(type) {
	case *userInterfaces.ErrorUserNotFound:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
		json, _ := json.Marshal(models.Error{Message: err.Error()})
		w.Write(json)
	case nil:
		jsonData, _ := json.Marshal(listOrgs)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(models.Error{Message: err.Error()})
		w.Write(json)
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
	switch err.(type) {
	case *userInterfaces.ErrorUserNotFound:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusNotFound)
		json, _ := json.Marshal(models.Error{Message: err.Error()})
		w.Write(json)
	case nil:
		jsonData, _ := json.Marshal(response)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(models.Error{Message: err.Error()})
		w.Write(json)
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
	switch err.(type) {
	case nil:
		jsonData, _ := json.Marshal(response)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	default:
		golog.Errorf("#%s: %w",  rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json, _ := json.Marshal(models.Error{Message: err.Error()})
		w.Write(json)
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
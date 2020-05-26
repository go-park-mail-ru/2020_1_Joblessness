package vacancyHttp

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"github.com/mailru/easyjson"
	"gopkg.in/go-playground/validator.v9"
	"joblessness/haha/models/base"
	"joblessness/haha/vacancy/interfaces"
	"net/http"
	"strconv"
)

type Handler struct {
	useCase vacancyInterfaces.VacancyUseCase
}

func NewHandler(useCase vacancyInterfaces.VacancyUseCase) *Handler {
	return &Handler{useCase}
}

func (h *Handler) CreateVacancy(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	var newVacancy baseModels.Vacancy
	newVacancy.Organization.ID = r.Context().Value("userID").(uint64)

	err := easyjson.UnmarshalFromReader(r.Body, &newVacancy)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = validator.New().Struct(newVacancy); err != nil {
		golog.Errorf("#%s: %s", rID, "Empty vacancy name")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newVacancy.ID, err = h.useCase.CreateVacancy(&newVacancy)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, _ := easyjson.Marshal(baseModels.ResponseID{ID: newVacancy.ID})
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (h *Handler) GetVacancy(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	vacancyID, _ := strconv.ParseUint(mux.Vars(r)["vacancy_id"], 10, 64)

	getVacancy, err := h.useCase.GetVacancy(vacancyID)
	switch err {
	case sql.ErrNoRows:
		golog.Errorf("#%s: %s", rID, err)
		w.WriteHeader(http.StatusNotFound)
	case nil:
		jsonData, _ := easyjson.Marshal(getVacancy)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	default:
		golog.Errorf("#%s: %s", rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) GetVacancies(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	page := r.FormValue("page")

	vacancies, err := h.useCase.GetVacancies(page)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, _ := easyjson.Marshal(vacancies)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *Handler) ChangeVacancy(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	vacancyID, _ := strconv.ParseUint(mux.Vars(r)["vacancy_id"], 10, 64)

	var newVacancy baseModels.Vacancy
	err := easyjson.UnmarshalFromReader(r.Body, &newVacancy)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newVacancy.ID = vacancyID
	newVacancy.Organization.ID = r.Context().Value("userID").(uint64)

	err = h.useCase.ChangeVacancy(&newVacancy)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	vacancyID, _ := strconv.ParseUint(mux.Vars(r)["vacancy_id"], 10, 64)

	authorID := r.Context().Value("userID").(uint64)

	err := h.useCase.DeleteVacancy(vacancyID, authorID)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetOrgVacancies(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	orgID, _ := strconv.ParseUint(mux.Vars(r)["organization_id"], 10, 64)

	vacancies, err := h.useCase.GetOrgVacancies(orgID)
	if err != nil {
		golog.Errorf("#%s: %s", rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, _ := easyjson.Marshal(vacancies)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

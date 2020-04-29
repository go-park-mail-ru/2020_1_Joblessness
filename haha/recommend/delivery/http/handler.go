package recommendHttp

import (
	"encoding/json"
	"errors"
	baseModels "joblessness/haha/models/base"
	"joblessness/haha/recommend/interfaces"
	"net/http"
	"strconv"
)

type Handler struct {
	useCase recommendInterfaces.RecommendUseCase
}

func NewHandler(useCase recommendInterfaces.RecommendUseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) GetRecommendedVacancies(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.ParseUint(r.FormValue("page"), 10, 64)

	vacancies, err := h.useCase.GetRecommendedVacancies(
		r.Context().Value("userID").(uint64),
		page - 1,
	)
	switch true {
	case errors.Is(err, recommendInterfaces.ErrNoUser):
		w.WriteHeader(http.StatusNotFound)
		jsonData, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(jsonData)
	case errors.Is(err, recommendInterfaces.ErrNoRecommendation):
		w.WriteHeader(http.StatusOK)
		jsonData, _ := json.Marshal([]baseModels.Vacancy{})
		w.Write(jsonData)
	case err == nil:
		w.WriteHeader(http.StatusOK)
		jsonData, _ := json.Marshal(vacancies)
		w.Write(jsonData)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		jsonData, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(jsonData)
	}
}

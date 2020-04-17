package recommendationHttp

import (
	"encoding/json"
	"errors"
	baseModels "joblessness/haha/models/base"
	"joblessness/haha/recommendation/interfaces"
	"net/http"
)

type Handler struct {
	useCase recommendationInterfaces.UseCase
}

func NewHandler(useCase recommendationInterfaces.UseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) GetRecommendedVacancies(w http.ResponseWriter, r *http.Request) {
	vacancies, err := h.useCase.GetRecommendedVacancies(
		r.Context().Value("userID").(uint64),
	)
	switch true {
	case errors.Is(err, recommendationInterfaces.ErrNoRecommendation):
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

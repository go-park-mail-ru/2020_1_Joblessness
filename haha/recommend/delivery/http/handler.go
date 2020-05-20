package recommendHttp

import (
	"encoding/json"
	"errors"
	"github.com/kataras/golog"
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
	rID := r.Context().Value("rID").(string)

	page, err := strconv.ParseUint(r.FormValue("page"), 10, 64)
	if page == 0 {
		page = 1
	}

	vacancies, err := h.useCase.GetRecommendedVacancies(
		r.Context().Value("userID").(uint64),
		page,
	)
	switch true {
	case errors.Is(err, recommendInterfaces.ErrNoUser):
		w.WriteHeader(http.StatusNotFound)
		jsonData, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(jsonData)
		golog.Errorf("#%s: %w", rID, err)
	case errors.Is(err, recommendInterfaces.ErrNoRecommendation):
		w.WriteHeader(http.StatusOK)
		jsonData, _ := json.Marshal([]baseModels.Vacancy{})
		w.Write(jsonData)
		golog.Errorf("#%s: %w", rID, err)
	case err == nil:
		w.WriteHeader(http.StatusOK)
		jsonData, _ := json.Marshal(vacancies)
		w.Write(jsonData)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		jsonData, _ := json.Marshal(baseModels.Error{Message: err.Error()})
		w.Write(jsonData)
		golog.Errorf("#%s: %w", rID, err)
	}
}

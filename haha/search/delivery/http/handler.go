package httpSearch

import (
	"encoding/json"
	"github.com/juju/loggo"
	"github.com/kataras/golog"
	searchInterfaces "joblessness/haha/search/interfaces"
	"net/http"
)

type Handler struct {
	useCase searchInterfaces.SearchUseCase
	logger  *loggo.Logger
}

func NewHandler(useCase searchInterfaces.SearchUseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	searchType := r.FormValue("type")
	request := r.FormValue("request")
	since := r.FormValue("since")
	desc := r.FormValue("desc")

	resultForum, err := h.useCase.Search(searchType, request, since, desc)

	switch err {
	case searchInterfaces.ErrUnknownRequest:
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusBadRequest)
	case nil:
		resultJSON, _ := json.Marshal(resultForum)
		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	default:
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

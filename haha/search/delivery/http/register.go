package searchHttp

import (
	"github.com/gorilla/mux"
	searchInterfaces "joblessness/haha/search/interfaces"
)

func RegisterHTTPEndpoints(router *mux.Router, uc searchInterfaces.SearchUseCase) {
	h := NewHandler(uc)

	router.HandleFunc("/search", h.Search).Methods("GET")
}

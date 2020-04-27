package recommendHttp

import (
	"github.com/gorilla/mux"
	"joblessness/haha/middleware"
	"joblessness/haha/recommend/interfaces"
)

func RegisterHTTPEndpoints(router *mux.Router, m *middleware.SessionHandler, u recommendInterfaces.RecommendUseCase) {
	h := NewHandler(u)

	router.HandleFunc("/recommend", m.PersonRequired(h.GetRecommendedVacancies)).Methods("GET")
}

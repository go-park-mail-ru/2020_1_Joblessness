package recommendationHttp

import (
	"github.com/gorilla/mux"
	"joblessness/haha/middleware"
	"joblessness/haha/recommendation/interfaces"
)

func RegisterHTTPEndpoints(router *mux.Router, m *middleware.SessionHandler, u recommendationInterfaces.UseCase) {
	h := NewHandler(u)

	router.HandleFunc("/recommendation", m.PersonRequired(h.GetRecommendedVacancies)).Methods("GET")
}

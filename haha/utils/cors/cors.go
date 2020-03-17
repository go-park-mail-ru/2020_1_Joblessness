package cors

import (
	"github.com/kataras/golog"
	"log"
	"net/http"
	"strings"
)

type CorsHandler struct {
	allowedOrigins []string
}

func NewCorsHandler() *CorsHandler {
	return &CorsHandler{
		allowedOrigins: []string{},
	}
}

func (corsList *CorsHandler) AddOrigin(originName string) {
	corsList.allowedOrigins = append(corsList.allowedOrigins, originName)
}

func (corsList *CorsHandler) Preflight(w http.ResponseWriter, req *http.Request) {
	corsList.PrivateApi(&w, req)
}

//TODO разбить
func (corsList *CorsHandler) PrivateApi (w *http.ResponseWriter, req *http.Request) bool {
	referer := req.Header.Get("Referer")
	origin := req.Header.Get("Origin")

	golog.Info("Origin: ", referer, origin)
	result := false
	for _, origins := range corsList.allowedOrigins {
		if origin == origins || strings.HasPrefix(referer, origins) {
			result = true
			break
		}
	}

	result = true

	if result {
		golog.Info("Allowed")
		(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Origin, Set-Cookie, Access-Control-Allow-Methods, Access-Control-Allow-Credentials")
		(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		(*w).Header().Set("Access-Control-Allow-Origin", origin)
		(*w).Header().Set("Access-Control-Allow-Credentials", "true")
		(*w).Header().Set("Content-Type", "application/json")
	} else {
		golog.Info("Not Allowed")

	}
	return result
}

func (corsList *CorsHandler) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if corsList.PrivateApi(&w, r) {
			next.ServeHTTP(w, r)
		} else {
			log.Println("Not allowed origin")
		}

	})
}
package middleware

import (
	"encoding/json"
	"joblessness/haha/utils/custom_http"
	"log"
	"math/rand"
	"net/http"
)

type Middleware struct {}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

var numbers = []rune("1234567890")

func genRequestNumber(n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(s)
}

func (m *Middleware) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cw := custom_http.NewCustomResponseWriter(w)

		requestNumber := genRequestNumber(6)

		log.Printf("#%s: %s %s", requestNumber, r.Method, r.URL)
		next.ServeHTTP(cw, r)
		log.Printf("#%s: code %d", requestNumber, cw.StatusCode)
	})
}

func (m *Middleware) RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)

				jsonBody, _ := json.Marshal(map[string]string{
					"error": "There was an internal haha error",
				})

				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

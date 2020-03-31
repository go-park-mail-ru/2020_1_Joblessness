package middleware

import (
	"context"
	"encoding/json"
	"github.com/kataras/golog"
	"joblessness/haha/utils/custom_http"
	"math/rand"
	"net/http"
)

type RecoveryHandler struct {}

func NewMiddleware() *RecoveryHandler {
	return &RecoveryHandler{}
}

var numbers = []rune("1234567890")

func genRequestNumber(n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(s)
}

func (m *RecoveryHandler) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := custom_http.NewStatusResponseWriter(w)

		requestNumber := genRequestNumber(6)
		r = r.WithContext(context.WithValue(r.Context(), "rID", requestNumber))

		golog.Infof("#%s: %s %s", requestNumber, r.Method, r.URL)
		next.ServeHTTP(sw, r)
		golog.Infof("#%s: code %d", requestNumber, sw.StatusCode)
	})
}

func (m *RecoveryHandler) RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rID, ok := r.Context().Value("rID").(string)

			err := recover()
			if err != nil {
				if ok {
					golog.Errorf("#%s: %w",  rID, err)
				} else {
					golog.Errorf("Panic with no id: %w", err)
				}

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
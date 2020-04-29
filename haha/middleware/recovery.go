package middleware

import (
	"context"
	"encoding/json"
	"github.com/kataras/golog"
	"github.com/prometheus/client_golang/prometheus"
	prom "joblessness/haha/prometheus"
	"joblessness/haha/utils/custom_http"
	"math/rand"
	"net/http"
	"time"
)

type RecoveryHandler struct{}

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

		labels := prometheus.Labels{
			"method": r.Method,
			"path":   r.URL.Path,
		}

		golog.Infof("#%s: %s %s", requestNumber, r.Method, r.URL)
		prom.RequestCurrent.With(labels).Inc()
		start := time.Now()

		next.ServeHTTP(w, r)

		prom.RequestDuration.With(labels).Observe(time.Since(start).Seconds())
		prom.RequestCurrent.With(labels).Dec()
		golog.Infof("#%s: code %d", requestNumber, sw.StatusCode)

		prom.RequestCount.With(labels).Inc()
	})
}

func (m *RecoveryHandler) RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rID, ok := r.Context().Value("rID").(string)

			err := recover()
			if err != nil {
				if ok {
					golog.Errorf("#%s Panic: %w", rID, err)
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

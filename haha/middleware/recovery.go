package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kataras/golog"
	"github.com/prometheus/client_golang/prometheus"
	prom "joblessness/haha/prometheus"
	"joblessness/haha/utils/custom_http"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
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

var rIDKey = "rID"

func formatPath(path string) string {
	pathArray := strings.Split(path[1:], "/")
	for i, _ := range pathArray {
		if _, err := strconv.Atoi(pathArray[i]); err == nil {
			pathArray[i] = "*"
		}
	}
	return "/" + strings.Join(pathArray, "/")
}

func (m *RecoveryHandler) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := custom_http.NewStatusResponseWriter(w)

		requestNumber := genRequestNumber(6)

		r = r.WithContext(context.WithValue(r.Context(), rIDKey, requestNumber))

		labels := prometheus.Labels{
			"method": r.Method,
			"path":   formatPath(r.URL.Path),
		}

		golog.Infof("#%s: %s %s", requestNumber, r.Method, r.URL)
		prom.RequestCurrent.With(labels).Inc()
		start := time.Now()

		next.ServeHTTP(sw, r)

		if r.URL.Path != "/api/metrics" {
			prom.RequestDuration.With(labels).Observe(float64(int(time.Since(start).Milliseconds())))
		}
		prom.RequestCurrent.With(labels).Dec()
		golog.Infof("#%s: code %d", requestNumber, sw.StatusCode)

		statusLabels := prometheus.Labels{
			"method": r.Method,
			"path":   r.URL.Path,
			"status": fmt.Sprintf("%d", sw.StatusCode),
		}
		prom.RequestCount.With(statusLabels).Inc()
	})
}

func (m *RecoveryHandler) RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rID, ok := r.Context().Value("rID").(string)

			err := recover()
			if err != nil {
				if ok {
					golog.Errorf("#%s Panic: %s", rID, err.(error).Error())
				} else {
					golog.Errorf("Panic with no id: %s", err.(error).Error())
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

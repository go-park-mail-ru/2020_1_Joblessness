package middleware

import (
	"context"
	"encoding/json"
	"github.com/kataras/golog"
	"joblessness/haha/auth"
	"joblessness/haha/utils/custom_http"
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
		sw := custom_http.NewStatusResponseWriter(w)

		requestNumber := genRequestNumber(6)
		r = r.WithContext(context.WithValue(r.Context(), "rID", requestNumber))

		golog.Infof("#%s: %s %s", requestNumber, r.Method, r.URL)
		next.ServeHTTP(sw, r)
		golog.Infof("#%s: code %d", requestNumber, sw.StatusCode)
	})
}

func (m *Middleware) RecoveryMiddleware(next http.Handler) http.Handler {
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

type AuthMiddleware struct {
	auth auth.AuthUseCase
}

func NewAuthMiddleware(authUseCase auth.AuthUseCase) *AuthMiddleware {
	return &AuthMiddleware{auth: authUseCase}
}

func (m *AuthMiddleware) CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rID, ok := r.Context().Value("rID").(string)
		if !ok {
			rID = "no request id"
		}

		session, err := r.Cookie("session_id")
		if err != nil {
			golog.Infof("#%s: %s",  rID, "No cookie")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		golog.Infof("#%s: %s",  rID, session.Value)
		userID, err := m.auth.SessionExists(session.Value)
		switch err {
		case auth.ErrWrongSID:
			golog.Errorf("#%s: %w",  rID, err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		case nil:
			golog.Infof("#%s: %s",  rID, "success")
		default:
			golog.Errorf("#%s: %w",  rID, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "userID", userID))
		next.ServeHTTP(w, r)
	}
}
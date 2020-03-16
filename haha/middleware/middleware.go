package middleware

import (
	"context"
	"encoding/json"
	"github.com/kataras/golog"
	"joblessness/haha/auth"
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
		sw := custom_http.NewStatusResponseWriter(w)

		requestNumber := genRequestNumber(6)
		r = r.WithContext(context.WithValue(r.Context(), "rID", requestNumber))

		log.Printf("#%s: %s %s", requestNumber, r.Method, r.URL)
		next.ServeHTTP(sw, r)
		log.Printf("#%s: code %d", requestNumber, sw.StatusCode)
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
	auth auth.UseCase
}

func NewAuthMiddleware(authUseCase auth.UseCase) *AuthMiddleware {
	return &AuthMiddleware{auth: authUseCase}
}

func (m *AuthMiddleware) CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rID := r.Context().Value("rID").(string)

		session, err := r.Cookie("session_id")
		log.Println("session cookie: ", session)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
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
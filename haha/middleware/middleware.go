package middleware

import (
	"context"
	"encoding/json"
	"joblessness/haha/auth"
	"joblessness/haha/utils/custom_http"
	"log"
	"math/rand"
	"net/http"
)

type Middleware struct {
}

func NewMiddleware(authUseCase auth.UseCase) *Middleware {
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
		r = r.WithContext(context.WithValue(r.Context(), "requestNumber", requestNumber))
		next.ServeHTTP(w, r)

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

type AuthMiddleware struct {
	auth auth.UseCase
}

func NewAuthMiddleware(authUseCase auth.UseCase) *AuthMiddleware {
	return &AuthMiddleware{auth: authUseCase}
}

func (m *AuthMiddleware) CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_id")
		log.Println("session cookie: ", session)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID, err := m.auth.SessionExists(session.Value)
		switch err {
		case auth.ErrWrongSID:
			log.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		case nil:
			log.Println("Success")
		default:
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "userID", userID))
		next.ServeHTTP(w, r)
	}
}
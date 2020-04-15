package middleware

import (
	"context"
	"github.com/kataras/golog"
	"joblessness/haha/auth/interfaces"
	"net/http"
)

type SessionHandler struct {
	auth authInterfaces.AuthUseCase
}

func NewAuthMiddleware(authUseCase authInterfaces.AuthUseCase) *SessionHandler {
	return &SessionHandler{auth: authUseCase}
}

func (m *SessionHandler) UserRequired(next http.HandlerFunc) http.HandlerFunc {
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
		case authInterfaces.ErrWrongSID:
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

func (m *SessionHandler) PersonRequired (next http.HandlerFunc) http.HandlerFunc {
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
		userID, err := m.auth.PersonSession(session.Value)
		switch err {
		case authInterfaces.ErrUserNotPerson:
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

func (m *SessionHandler) OrganizationRequired (next http.HandlerFunc) http.HandlerFunc {
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
		userID, err := m.auth.OrganizationSession(session.Value)
		switch err {
		case authInterfaces.ErrUserNotOrganization:
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

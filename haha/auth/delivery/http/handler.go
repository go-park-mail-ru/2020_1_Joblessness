package authHttp

import (
	"encoding/json"
	"github.com/juju/loggo"
	"github.com/kataras/golog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/go-playground/validator.v9"
	"joblessness/haha/auth/interfaces"
	"joblessness/haha/models/base"
	"net/http"
	"time"
)

type Handler struct {
	useCase authInterfaces.AuthUseCase
	logger  *loggo.Logger
}

func NewHandler(useCase authInterfaces.AuthUseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) RegisterPerson(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	var user baseModels.Person
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = validator.New().Struct(user); err != nil {
		golog.Errorf("#%s: %s", rID, "Empty login or password")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.RegisterPerson(user.Login, user.Password, user.FirstName)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				golog.Errorf("#%s: %w", rID, err)
				w.WriteHeader(http.StatusBadRequest)
				json, _ := json.Marshal(baseModels.Error{Message: e.Message()})
				w.Write(json)
			default:
				golog.Errorf("#%s: %w", rID, err)
				w.WriteHeader(http.StatusInternalServerError)
				json, _ := json.Marshal(baseModels.Error{Message: e.Message()})
				w.Write(json)
			}
		} else {
			golog.Error(authInterfaces.ErrParseGrpcError.Error())
		}
	} else {
		golog.Infof("#%s: %s", rID, "Успешно")
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) RegisterOrg(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	var org baseModels.Organization
	err := json.NewDecoder(r.Body).Decode(&org)
	if err != nil {
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = validator.New().Struct(org); err != nil {
		golog.Errorf("#%s: %s", rID, "Empty login or password")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.RegisterOrganization(org.Login, org.Password, org.Name)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				golog.Errorf("#%s: %s", rID, e.Message())
				w.WriteHeader(http.StatusBadRequest)
				json, _ := json.Marshal(baseModels.Error{Message: e.Message()})
				w.Write(json)
			default:
				golog.Errorf("#%s: %w", rID, err)
				w.WriteHeader(http.StatusInternalServerError)
				json, _ := json.Marshal(baseModels.Error{Message: e.Message()})
				w.Write(json)
			}
		} else {
			golog.Error(authInterfaces.ErrParseGrpcError.Error())
		}
	} else {
		golog.Infof("#%s: %s", rID, "success")
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)

	var user baseModels.UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	golog.Debugf("#%s: %w", rID, user)
	if err != nil {
		golog.Errorf("#%s: %w\n%w", rID, err, user)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = validator.New().Struct(user); err != nil {
		golog.Errorf("#%s: %s", rID, "Empty login or password")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, role, sessionId, err := h.useCase.Login(user.Login, user.Password)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				golog.Errorf("#%s: %w", rID, err)
				w.WriteHeader(http.StatusBadRequest)
				json, _ := json.Marshal(baseModels.Error{Message: e.Message()})
				w.Write(json)
			default:
				golog.Errorf("#%s: %w", rID, err)
				w.WriteHeader(http.StatusInternalServerError)
				json, _ := json.Marshal(baseModels.Error{Message: e.Message()})
				w.Write(json)
			}
		} else {
			golog.Error(authInterfaces.ErrParseGrpcError.Error())
		}
	} else {
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    sessionId,
			Expires:  time.Now().Add(time.Hour),
			MaxAge:   100000,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, cookie)

		jsonData, _ := json.Marshal(baseModels.ResponseRole{ID: userId, Role: role})
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonData)
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = h.useCase.Logout(session.Value)
	if err != nil {
		golog.Errorf("#%s: %w", rID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	session.Path = "/"
	http.SetCookie(w, session)

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value("rID").(string)
	userID, _ := r.Context().Value("userID").(uint64)

	role, err := h.useCase.GetRole(userID)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				golog.Errorf("#%s: %w", rID, err)
				w.WriteHeader(http.StatusNotFound)
				json, _ := json.Marshal(baseModels.Error{Message: e.Message()})
				w.Write(json)
			default:
				golog.Errorf("#%s: %w", rID, err)
				w.WriteHeader(http.StatusInternalServerError)
				json, _ := json.Marshal(baseModels.Error{Message: e.Message()})
				w.Write(json)
			}
		} else {
			golog.Error(authInterfaces.ErrParseGrpcError.Error())
		}
	} else {
		jsonData, _ := json.Marshal(baseModels.ResponseRole{ID: userID, Role: role})
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonData)
	}
}

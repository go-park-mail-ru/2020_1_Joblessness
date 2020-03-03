package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"io"
	"joblessness/haha/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var sess = session.Must(session.NewSession(&aws.Config{
	Region: aws.String("ru-msk"),
	Credentials: credentials.NewStaticCredentials("orFNtcQG9pi8NvqcFhLAj4",
		"33CiuS769M4u1wHAk42HhdtCrCb795MGuez3biaE3CeK", ""),
	Endpoint: aws.String("https://hb.bizmrg.com"),
}))
var svc = s3.New(sess)

func (api *AuthHandler) AuthRequiredMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, found := api.Sessions[session.Value]
		if !found {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *AuthHandler) GetUserPage(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /user/{user_id}")

	var currentUser *models.User
	userId, _ := strconv.Atoi(mux.Vars(r)["user_id"])

	for _, user := range api.Users {
		if (*user).ID == uint(userId) {
			currentUser = user
		}
	}

	if currentUser == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userAvatar := "" //deprecated

	type Response struct {
		User models.UserInfo `json:"user"`
		Summaries []models.UserSummary `json:"summaries"`
	}

	jsonData, _ := json.Marshal(Response{
		models.UserInfo{currentUser.FirstName, currentUser.LastName, "", userAvatar},
		[]models.UserSummary{},
	})
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (api *AuthHandler) ChangeUserInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("PUT /user/{user_id}")

	session, err := r.Cookie("session_id")
	log.Println("session cookie: ", session)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userId, found := api.Sessions[session.Value]
	if !found {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if reqId, _ := strconv.Atoi(mux.Vars(r)["user_id"]); uint(reqId) != userId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var currentUser *models.User

	log.Println("Users counter", len(api.Users))
	for _, user := range api.Users {
		if (*user).ID == userId {
			currentUser = user
		}
	}

	if currentUser == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var data map[string]string
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&data)
//TODO Проверять есть ли все поля
	(*currentUser).LastName = data["last-name"]
	(*currentUser).FirstName = data["first-name"]

	w.WriteHeader(http.StatusNoContent)
}

func (api *AuthHandler) SetAvatar(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /users/{user_id}/avatar")

	session, err := r.Cookie("session_id")
	log.Println("session cookie: ", session)
	//get sessionn id from cookie
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// get cookie by session id
	userId, found := api.Sessions[session.Value]

	if !found {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// compare id from request and cookie
	if reqId, _ := strconv.Atoi(mux.Vars(r)["user_id"]); uint(reqId) != userId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var currentUser *models.User
	// get user data
	for _, user := range api.Users {
		// binary search could work a lot quicker
		// but there's no need to implement it
		// since DB is responsible for this
		if (*user).ID == userId {
			currentUser = user
		}
	}

	if currentUser == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(1024 * 1024 * 5) //5mb
	if err != nil {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer file.Close()

	var buf bytes.Buffer
	io.Copy(&buf, file)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("imgs-hh"),
		Key:    aws.String(currentUser.Login + "-avatar-" + header.Filename),
		Body:   strings.NewReader(buf.String()),
	})

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusFailedDependency)
		return
	}

	api.UserAvatars[1] = currentUser.Login + "-avatar-" + header.Filename

	w.WriteHeader(http.StatusCreated)
}

func (api *AuthHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /users/{user_id}/avatar")

	var currentUser *models.User
	userId, _ := strconv.Atoi(mux.Vars(r)["user_id"])

	for _, user := range api.Users {
		if (*user).ID == uint(userId) {
			currentUser = user
		}
	}
	// Check if user exists
	if currentUser == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// try to load user's avatar
	body, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("imgs-hh"),
		Key:    aws.String(api.UserAvatars[uint(userId)]),
	})
	// Not found or no response
	if err != nil {
		body, err = svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String("imgs-hh"),
			Key:    aws.String("default-avatar.png"),
		})
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
	defer body.Body.Close()

	var buf bytes.Buffer
	io.Copy(&buf, body.Body)
	content := strings.NewReader(buf.String())

	http.ServeContent(w, r, api.UserAvatars[1], *body.LastModified, content)
	buf.Reset()
}


package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/Teriton/chemistryBack/internal/models"
	"github.com/Teriton/chemistryBack/pkg/authmngr"
	"github.com/Teriton/chemistryBack/pkg/dbrepo"
)

type UserWithCompletedLessonsCount struct {
	models.User
	CompletedLessonsCount int `json:"completed_lessons"`
}

type UserWithPasswordToEdit struct {
	models.AddUser
	CurrentPassword string `json:"current_password"`
}

type UserHandler struct {
	authMngr authmngr.AuthorizationMngr
	dbRepo   dbrepo.DBRepo
}

func NewUserHandler(authMngr authmngr.AuthorizationMngr, dbRepo dbrepo.DBRepo) (*UserHandler, error) {
	return &UserHandler{authMngr, dbRepo}, nil
}

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] /user")
	w.Header().Set("Content-Type", "application/json")
	cookies := r.CookiesNamed("token")

	var err error
	if len(cookies) < 1 {
		err = errors.New("cookie is not set")
	}
	if checkError(w, err, http.StatusForbidden) {
		return
	}

	jwtToken := cookies[0].Value
	jwtContent, err := uh.authMngr.VerifyToken(jwtToken)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	user, err := uh.dbRepo.GetUserByUserName(jwtContent.Username)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) GetUserWithCopletedLessosnCount(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] /user/completedLessonsCount")
	w.Header().Set("Content-Type", "application/json")
	cookies := r.CookiesNamed("token")

	var err error
	if len(cookies) < 1 {
		err = errors.New("cookie is not set")
	}
	if checkError(w, err, http.StatusForbidden) {
		return
	}

	jwtToken := cookies[0].Value
	jwtContent, err := uh.authMngr.VerifyToken(jwtToken)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	user, err := uh.dbRepo.GetUserByUserName(jwtContent.Username)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	completedLessons, err := uh.dbRepo.GetCompletedLessonsLenForUser(user.ID)
	if checkError(w, err, http.StatusForbidden) {
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(UserWithCompletedLessonsCount{user, completedLessons})
}

func (uh *UserHandler) EditUserInfo(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] /user/edit")
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if checkError(w, err, http.StatusBadRequest) {
		return
	}
	var userData UserWithPasswordToEdit
	err = json.Unmarshal(body, &userData)
	if checkError(w, err, http.StatusBadRequest) {
		return
	}
	cookies := r.CookiesNamed("token")

	if len(cookies) < 1 {
		err = errors.New("cookie is not set")
	}
	if checkError(w, err, http.StatusForbidden) {
		return
	}

	jwtToken := cookies[0].Value
	jwtContent, err := uh.authMngr.VerifyPasswordAndToken(jwtToken, userData.CurrentPassword)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	jwt, err := uh.authMngr.EditUserInfo(userData.AddUser, jwtContent.Username)
	if checkError(w, err, http.StatusForbidden) {
		return
	}

	cookie := createCookieJWT(jwt)
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

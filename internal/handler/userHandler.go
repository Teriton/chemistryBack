package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

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

	if userData.Password == "" {
		userData.Password = userData.CurrentPassword
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

type practiceXPRequest struct {
	CorrectCount int `json:"correct_count"`
}

// AddPracticeXP начисляет XP за практический модуль (без записи в lessons_completed).
func (uh *UserHandler) AddPracticeXP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] /user/practice-xp")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		checkError(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if checkError(w, err, http.StatusBadRequest) {
		return
	}
	var req practiceXPRequest
	err = json.Unmarshal(body, &req)
	if checkError(w, err, http.StatusBadRequest) {
		return
	}
	const maxCorrect = 20
	const xpPerCorrect = 50
	if req.CorrectCount < 0 || req.CorrectCount > maxCorrect {
		checkError(w, errors.New("correct_count out of range"), http.StatusBadRequest)
		return
	}
	xp := req.CorrectCount * xpPerCorrect
	if xp == 0 {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]any{"status": "ok", "xp_added": 0})
		return
	}
	cookies := r.CookiesNamed("token")
	if len(cookies) < 1 {
		checkError(w, errors.New("cookie is not set"), http.StatusForbidden)
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
	err = uh.dbRepo.AddXPToUser(user.ID, xp)
	if checkError(w, err, http.StatusInternalServerError) {
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]any{"status": "ok", "xp_added": xp})
}

type avatarBody struct {
	Avatar string `json:"avatar"`
}

// SetAvatar сохраняет data URL изображения (или пустую строку для сброса на значение по умолчанию).
func (uh *UserHandler) SetAvatar(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] /user/avatar")
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		checkError(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if checkError(w, err, http.StatusBadRequest) {
		return
	}
	var req avatarBody
	err = json.Unmarshal(body, &req)
	if checkError(w, err, http.StatusBadRequest) {
		return
	}
	const maxLen = 400000
	if len(req.Avatar) > maxLen {
		checkError(w, errors.New("avatar too large"), http.StatusBadRequest)
		return
	}
	if req.Avatar != "" {
		if !strings.HasPrefix(req.Avatar, "data:image/") || !strings.Contains(req.Avatar, ";base64,") {
			checkError(w, errors.New("invalid avatar format"), http.StatusBadRequest)
			return
		}
	}
	cookies := r.CookiesNamed("token")
	if len(cookies) < 1 {
		checkError(w, errors.New("cookie is not set"), http.StatusForbidden)
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
	err = uh.dbRepo.UpdateUserAvatar(user.ID, req.Avatar)
	if checkError(w, err, http.StatusInternalServerError) {
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (uh *UserHandler) GetAchievements(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] /user/achievements")
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
	completedAchievements, err := uh.dbRepo.GetCompletedAchievementsForUser(user.ID)
	if checkError(w, err, http.StatusForbidden) {
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(completedAchievements)
}

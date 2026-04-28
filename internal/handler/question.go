package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/Teriton/chemistryBack/pkg/articlemngr"
	"github.com/Teriton/chemistryBack/pkg/authmngr"
	"github.com/Teriton/chemistryBack/pkg/dbrepo"
)

type CompleteDataRequest struct {
	LessonTitle string `json:"lesson_title"`
	Xp          int    `json:"xp"`
}

type QuestionHandler struct {
	articleMngr articlemngr.ArticleMngr
	authMngr    authmngr.AuthorizationMngr
	dbRepo      dbrepo.DBRepo
}

func NewQuestionHandler(
	articleMngr articlemngr.ArticleMngr,
	authMngr authmngr.AuthorizationMngr,
	dbRepo dbrepo.DBRepo,
) (*QuestionHandler, error) {
	return &QuestionHandler{
		articleMngr,
		authMngr,
		dbRepo,
	}, nil
}

func (qh *QuestionHandler) CompleteArticle(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] /complete")
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
	jwtContent, err := qh.authMngr.VerifyToken(jwtToken)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	body, err := io.ReadAll(r.Body)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	var completeDataRequest CompleteDataRequest
	err = json.Unmarshal(body, &completeDataRequest)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	if completeDataRequest.LessonTitle == "" {
		checkError(w, errors.New("can't unmarshal request body'"), http.StatusForbidden)
		return
	}

	err = qh.articleMngr.CompleteLesson(
		jwtContent.Username,
		completeDataRequest.LessonTitle,
		completeDataRequest.Xp,
	)
	if err != nil && err.Error() == "lesson already completed" {
		infoMessage(w, err.Error(), http.StatusAlreadyReported)
		return
	} else if checkError(w, err, http.StatusForbidden) {
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

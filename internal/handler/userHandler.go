package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Teriton/chemistryBack/pkg/authmngr"
	"github.com/Teriton/chemistryBack/pkg/dbrepo"
)

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
	username, err := uh.authMngr.Verify(jwtToken)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	user, err := uh.dbRepo.GetUserByUserName(username)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(user)
}

package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Teriton/chemistryBack/internal/models"
	"github.com/Teriton/chemistryBack/pkg/authmngr"
)

type LoginDataRequst struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func createCookieJWT(jwt string) http.Cookie {
	return http.Cookie{
		Name:     "token",
		Value:    jwt,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}
func deleteCookieJWT() http.Cookie {
	return http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}

type AuthHandler struct {
	authMngr authmngr.AuthorizationMngr
}

func NewAuthHandler(authMngr authmngr.AuthorizationMngr) (*AuthHandler, error) {
	return &AuthHandler{authMngr}, nil
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] /login")
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	var loginData LoginDataRequst
	err = json.Unmarshal(body, &loginData)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	jwt, err := h.authMngr.Login(loginData.Login, loginData.Password)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	cookie := createCookieJWT(jwt)
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] /signup")
	body, err := io.ReadAll(r.Body)
	if checkError(w, err, http.StatusBadRequest) {
		return
	}
	var signupData models.AddUser
	err = json.Unmarshal(body, &signupData)
	if checkError(w, err, http.StatusBadRequest) {
		return
	}
	jwt, err := h.authMngr.Signup(signupData)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	cookie := createCookieJWT(jwt)
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] /logout")
	w.Header().Set("Content-Type", "application/json")
	cookie := deleteCookieJWT()
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

//func (h *AuthHandler) signupTest(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("Signup")
//	body, err := io.ReadAll(r.Body)
//	var tempUser struct {
//		Username string `json:"username"`
//	}
//	if err != nil {
//		fmt.Println(err)
//	}
//	err = json.Unmarshal(body, &tempUser)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Printf("Raw: %s\nResult: %#v ", body, tempUser)
//}

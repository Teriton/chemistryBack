// Package authmngr to manage users
package authmngr

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Teriton/chemistryBack/internal/models"
	"github.com/Teriton/chemistryBack/pkg/dbrepo"
)

type AuthorizationMngr interface {
	Login(username string, password string) (string, error) // Returns JWT token
	Signup(models.AddUser) (string, error)                  // Returns JWT token

}

func checkSignupData(user models.AddUser) error {
	switch {
	case user.Username == "":
		return errors.New("incorrect username")
	case len(user.Password) < 4:
		return errors.New("password must be at least 4 characters")
	case user.Email == "":
		return errors.New("incorrect email")
	}
	return nil
}

type AuthenticationMngr interface {
	verify(jwt string) (models.User, error)
}

type Mngr struct {
	dbrepo    dbrepo.DBRepo
	pswHasher Hasher
}

func NewAuthMngr(dbrepo dbrepo.DBRepo, pswHaser Hasher) (*Mngr, error) {
	return &Mngr{dbrepo, pswHaser}, nil
}

func generateJWT(username string) (string, error) {
	key := os.Getenv("JWT_SECRET_TOKEN")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	signed, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (m Mngr) Login(username string, password string) (string, error) {
	user, err := m.dbrepo.GetUserByUserName(username)
	if err != nil {
		return "", err
	}
	if !m.pswHasher.check([]byte(user.Password), []byte(password)) {
		return "", errors.New("incorect password")
	}
	token, err := generateJWT(user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (m Mngr) Signup(userToAdd models.AddUser) (string, error) {
	err := checkSignupData(userToAdd)
	if err != nil {
		return "", err
	}
	hashedPassword, err := m.pswHasher.encode([]byte(userToAdd.Password))
	if err != nil {
		return "", nil
	}
	userToAdd.Password = string(hashedPassword)
	err = m.dbrepo.CreateUser(userToAdd)
	if err != nil {
		return "", err
	}
	token, err := generateJWT(userToAdd.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

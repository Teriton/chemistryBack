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
	login(username string, password string) (string, error) // Returns JWT token
	signup(models.AddUser) (bool, error)
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

func (m Mngr) login(username string, password string) (string, error) {
	user, err := m.dbrepo.GetUserByUserName(username)
	if err != nil {
		return "", err
	}
	if !m.pswHasher.check([]byte(user.Password), []byte(password)) {
		return "", errors.New("incorect password")
	}
	key := os.Getenv("JWT_SECRET_TOKEN")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	signed, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (m Mngr) signup(models.AddUser) (bool, error) {
	return false, nil
}

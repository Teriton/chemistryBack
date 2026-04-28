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

type JWTContent struct {
	Username string
	UserID   int
}

type AuthorizationMngr interface {
	Login(username string, password string) (string, error) // Returns JWT token
	Signup(models.AddUser) (string, error)                  // Returns JWT token
	VerifyToken(jwtToken string) (JWTContent, error)
	EditUserInfo(models.AddUser, string) (string, error)
	VerifyPasswordAndToken(string, string) (JWTContent, error)
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

type JWTClaims struct {
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
	jwt.RegisteredClaims
}

type Mngr struct {
	dbrepo    dbrepo.DBRepo
	pswHasher Hasher
}

func NewAuthMngr(dbrepo dbrepo.DBRepo, pswHaser Hasher) (*Mngr, error) {
	return &Mngr{dbrepo, pswHaser}, nil
}

func generateJWT(username string, userID int) (string, error) {
	key := os.Getenv("JWT_SECRET_TOKEN")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"user_id":  userID,
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
	token, err := generateJWT(user.Username, user.ID)
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
	user, err := m.dbrepo.GetUserByUserName(userToAdd.Username)
	if err != nil {
		return "", err
	}
	token, err := generateJWT(userToAdd.Username, user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (m Mngr) VerifyToken(jwtToken string) (JWTContent, error) {
	key := os.Getenv("JWT_SECRET_TOKEN")
	parseToken, err := jwt.ParseWithClaims(
		jwtToken,
		&JWTClaims{},
		func(token *jwt.Token) (any, error) {
			return []byte(key), nil
		},
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return JWTContent{}, err
	}
	if claims, ok := parseToken.Claims.(*JWTClaims); ok {
		if err != nil {
			return JWTContent{}, err
		}
		return JWTContent{Username: claims.Username, UserID: claims.UserID}, nil
	}

	return JWTContent{}, errors.New("error while parsing token")
}

func (m Mngr) EditUserInfo(userToAdd models.AddUser, currentUsername string) (string, error) {
	err := checkSignupData(userToAdd)
	if err != nil {
		return "", err
	}
	hashedPassword, err := m.pswHasher.encode([]byte(userToAdd.Password))
	if err != nil {
		return "", nil
	}
	userToAdd.Password = string(hashedPassword)
	err = m.dbrepo.EditUser(userToAdd, currentUsername)
	if err != nil {
		return "", err
	}
	user, err := m.dbrepo.GetUserByUserName(userToAdd.Username)
	if err != nil {
		return "", err
	}
	token, err := generateJWT(userToAdd.Username, user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (m Mngr) VerifyPasswordAndToken(jwtToken string, password string) (JWTContent, error) {
	key := os.Getenv("JWT_SECRET_TOKEN")
	parseToken, err := jwt.ParseWithClaims(
		jwtToken,
		&JWTClaims{},
		func(token *jwt.Token) (any, error) {
			return []byte(key), nil
		},
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return JWTContent{}, err
	}

	claims, ok := parseToken.Claims.(*JWTClaims)
	if !ok {
		return JWTContent{}, errors.New("error while parsing token")
	}

	user, err := m.dbrepo.GetUserByUserName(claims.Username)
	if err != nil {
		return JWTContent{}, err
	}

	if !m.pswHasher.check([]byte(user.Password), []byte(password)) {
		return JWTContent{}, errors.New("incorect password")
	}

	return JWTContent{Username: claims.Username, UserID: claims.UserID}, nil
}

package authmngr

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Teriton/chemistryBack/pkg/dbrepo"
)

type TestJWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func createAuthMngr(t *testing.T) *Mngr {
	var psqlRepo dbrepo.DBRepo
	psqlRepo, err := dbrepo.NewPsqlRepo(os.Getenv("POSTGRESQL_URL"))
	check(err, t)
	var pswHasher Hasher
	pswHasher, err = NewPasswordHasher()
	check(err, t)
	authMngr, err := NewAuthMngr(psqlRepo, pswHasher)
	check(err, t)
	return authMngr
}

func check(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

func TestCreateAndVerifyJWT(t *testing.T) {
	key := os.Getenv("JWT_SECRET_TOKEN")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "Shpack",
		"exp":      time.Now().Add(time.Second * 2).Unix(),
	})
	signed, err := token.SignedString([]byte(key))
	check(err, t)
	fmt.Println(signed)

	parseToken, err := jwt.ParseWithClaims(signed, &TestJWTClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	}, jwt.WithExpirationRequired())
	check(err, t)
	if claims, ok := parseToken.Claims.(*TestJWTClaims); ok {
		fmt.Println(claims.Username, claims.Issuer)
	} else {
		t.Error("Error while parsing token")
	}
}

func TestLogin(t *testing.T) {
	authMngr := createAuthMngr(t)
	jwt, err := authMngr.login("Shpack", "654321")
	check(err, t)
	fmt.Println(jwt)
}

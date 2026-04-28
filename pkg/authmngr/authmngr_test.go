package authmngr

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Teriton/chemistryBack/internal/models"
	"github.com/Teriton/chemistryBack/pkg/dbrepo"
)

func createAuthMngr(t *testing.T) (*Mngr, *dbrepo.PsqlRepo) {
	var psqlRepo dbrepo.DBRepo
	psqlRepo, err := dbrepo.NewPsqlRepo(os.Getenv("POSTGRESQL_URL"))
	psqlRepoAsserted, ok := psqlRepo.(*dbrepo.PsqlRepo)
	if !ok {
		t.Error("Cant assert dbrepo to psqlrepo")
		return nil, nil
	}
	check(err, t)
	var pswHasher Hasher
	pswHasher, err = NewPasswordHasher()
	check(err, t)
	authMngr, err := NewAuthMngr(psqlRepo, pswHasher)
	check(err, t)
	return authMngr, psqlRepoAsserted
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

	parseToken, err := jwt.ParseWithClaims(signed, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	}, jwt.WithExpirationRequired())
	check(err, t)
	if claims, ok := parseToken.Claims.(*JWTClaims); ok {
		fmt.Println(claims.Username, claims.Issuer)
	} else {
		t.Error("Error while parsing token")
	}
}

func TestLogin(t *testing.T) {
	authMngr, _ := createAuthMngr(t)
	jwt, err := authMngr.Login("Shpack", "654321")
	check(err, t)
	fmt.Println(jwt)
}

func TestSignup(t *testing.T) {
	authMngr, dbRepo := createAuthMngr(t)
	userToAdd := models.AddUser{
		Email:    "test@test.test",
		Password: "1234",
		Username: "Tester",
	}
	jwt, err := authMngr.Signup(userToAdd)
	check(err, t)
	err = dbRepo.DeleteUserByUserName(userToAdd.Username)
	check(err, t)
	fmt.Println(jwt)
}

func TestEdit(t *testing.T) {
	authMngr, dbRepo := createAuthMngr(t)
	userToAdd := models.AddUser{
		Email:    "test@test.test",
		Password: "1234",
		Username: "Tester",
	}
	jwt, err := authMngr.Signup(userToAdd)
	check(err, t)
	println(jwt)
	currentUsername := userToAdd.Username
	userToAdd.Username = "Tester2"
	jwt, err = authMngr.EditUserInfo(userToAdd, currentUsername)
	check(err, t)
	println(jwt)
	err = dbRepo.DeleteUserByUserName(userToAdd.Username)
	check(err, t)
}
func TestVerify(t *testing.T) {
	authMngr, _ := createAuthMngr(t)
	userToLogin := models.AddUser{
		Password: "654321",
		Username: "Shpack",
	}
	jwt, err := authMngr.Login(
		userToLogin.Username,
		userToLogin.Password,
	)

	check(err, t)

	fmt.Println("Shpack: ", jwt)
	user, err := authMngr.VerifyToken(jwt)
	check(err, t)
	fmt.Println(user)
}

func TestVerifyTokenAndPassword(t *testing.T) {
	authMngr, _ := createAuthMngr(t)
	userToLogin := models.AddUser{
		Password: "654321",
		Username: "Shpack",
	}
	jwt, err := authMngr.Login(
		userToLogin.Username,
		userToLogin.Password,
	)

	check(err, t)

	fmt.Println("Shpack: ", jwt)
	user, err := authMngr.VerifyPasswordAndToken(jwt, userToLogin.Password)
	check(err, t)
	fmt.Println(user)
}

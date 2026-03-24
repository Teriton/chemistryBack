package dbrepo

import (
	"fmt"
	"os"
	"testing"

	"github.com/Teriton/chemistryBack/internal/models"
)

func CreatePsql(t *testing.T) DBRepo {
	var dbRepo DBRepo
	dbRepo, err := NewPsqlRepo(os.Getenv("POSTGRESQL_URL"))
	check(err, t)
	return dbRepo
}

func TestPsqlDBCreate(t *testing.T) {
	dbRepo := CreatePsql(t)
	defer dbRepo.CloseDB()
}

func TestGetUserByUserName(t *testing.T) {
	dbRepo := CreatePsql(t)
	defer dbRepo.CloseDB()

	user, err := dbRepo.GetUserByUserName("Vitaly")
	check(err, t)
	fmt.Println(user)
}

func TestCreateUserAndDelteUser(t *testing.T) {
	dbRepo := CreatePsql(t)
	defer dbRepo.CloseDB()

	userToAdd := models.AddUser{Email: "sosi@gmail.com", Password: "12345", Username: "Shpack"}
	err := dbRepo.CreateUser(userToAdd)
	check(err, t)
	err = dbRepo.DeleteUserByUserName(userToAdd.Username)
	check(err, t)
}

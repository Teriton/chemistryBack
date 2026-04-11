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

	user, err := dbRepo.GetUserByUserName("Shpack")
	check(err, t)
	fmt.Println(user)
}

func TestCreateUserAndDelteUser(t *testing.T) {
	dbRepo := CreatePsql(t)
	defer dbRepo.CloseDB()

	userToAdd := models.AddUser{Email: "sosi@gmail.com", Password: "12345", Username: "Shpachok"}
	err := dbRepo.CreateUser(userToAdd)
	check(err, t)
	err = dbRepo.DeleteUserByUserName(userToAdd.Username)
	check(err, t)
}

func TestLessons(t *testing.T) {
	dbRepo := CreatePsql(t)
	defer dbRepo.CloseDB()
	title := "ShpackTest"

	_, err := dbRepo.GetLessonByTitle(title)
	if err == nil {
		t.Error("an error should be here")
		return
	}
	err = dbRepo.CreateLessonWithTitle(title)
	check(err, t)
	lesson, err := dbRepo.GetLessonByTitle(title)
	check(err, t)
	fmt.Printf("Title: %s\n ID: %v\n", lesson.Title, lesson.ID)
	err = dbRepo.DeleteLessonByTitle(title)
	check(err, t)
}

func TestCompletedLessons(t *testing.T) {
	dbRepo := CreatePsql(t)
	defer dbRepo.CloseDB()
	title := "ShpackTest"
	user, err := dbRepo.GetUserByUserName("Shpack")
	check(err, t)
	err = dbRepo.CreateLessonWithTitle(title)
	check(err, t)
	lesson, err := dbRepo.GetLessonByTitle(title)
	check(err, t)

	err = dbRepo.CreateCompletedLesson(user.ID, lesson.ID)
	check(err, t)
	count, err := dbRepo.GetCompletedLessonsLenForUser(user.ID)
	check(err, t)
	fmt.Printf("Lessons complited by user Shpack: %v\n", count)
	if count < 1 {
		t.Error("lessons should be more than 0")
	}
	err = dbRepo.DeleteCompletedLesson(user.ID, lesson.ID)
	check(err, t)
	err = dbRepo.DeleteLessonByTitle(title)
	check(err, t)
}

func TestChangingXP(t *testing.T) {
	dbRepo := CreatePsql(t)
	defer dbRepo.CloseDB()

	user, err := dbRepo.GetUserByUserName("Shpack")
	check(err, t)

	err = dbRepo.AddXPToUser(user.ID, 100)
	check(err, t)
	err = dbRepo.AddXPToUser(user.ID, -100)
	check(err, t)
}

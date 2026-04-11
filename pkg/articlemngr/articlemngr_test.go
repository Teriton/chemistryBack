package articlemngr

import (
	"fmt"
	"os"
	"testing"

	"github.com/Teriton/chemistryBack/pkg/dbrepo"
)

func check(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

func createArticleMngr(t *testing.T) (*ArticleMngr, *dbrepo.PsqlRepo) {
	var psqlRepo dbrepo.DBRepo
	psqlRepo, err := dbrepo.NewPsqlRepo(os.Getenv("POSTGRESQL_URL"))
	psqlRepoAsserted, ok := psqlRepo.(*dbrepo.PsqlRepo)
	if !ok {
		t.Error("Cant assert dbrepo to psqlrepo")
		return nil, nil
	}
	check(err, t)
	articleMngr, err := NewArticleMngr(psqlRepo)
	check(err, t)
	return articleMngr, psqlRepoAsserted
}

func TestArticleMngr(t *testing.T) {
	articleMngr, dbRepo := createArticleMngr(t)
	defer dbRepo.CloseDB()
	title := "ShpackTest"
	user, err := dbRepo.GetUserByUserName("Shpack")
	check(err, t)

	err = articleMngr.CompleteLesson(user.Username, title, 100)
	check(err, t)
	count, err := dbRepo.GetCompletedLessonsLenForUser(user.ID)
	check(err, t)
	fmt.Printf("Shpack completed %v lessons\n", count)
	dbRepo.DeleteLessonByTitle(title)
}

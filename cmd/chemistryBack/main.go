package main

import (
	"os"

	"github.com/Teriton/chemistryBack/internal/app"
	"github.com/Teriton/chemistryBack/pkg/articlemngr"
	"github.com/Teriton/chemistryBack/pkg/articlereader"
	"github.com/Teriton/chemistryBack/pkg/authmngr"
	"github.com/Teriton/chemistryBack/pkg/dbrepo"
)

func checkForError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	articleReader := articlereader.NewDirReader("articles")
	dbRepo, err := dbrepo.NewPsqlRepo(os.Getenv("POSTGRESQL_URL"))
	defer dbRepo.CloseDB()
	checkForError(err)
	pswHasher, err := authmngr.NewPasswordHasher()
	checkForError(err)
	authMngr, err := authmngr.NewAuthMngr(dbRepo, pswHasher)
	checkForError(err)
	articleMngr, err := articlemngr.NewArticleMngr(dbRepo)
	checkForError(err)
	app := app.NewApp(articleReader, authMngr, dbRepo, *articleMngr, ":8080")
	app.Run()
}

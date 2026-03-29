package main

import (
	"os"

	"github.com/Teriton/chemistryBack/internal/app"
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
	checkForError(err)
	pswHasher, err := authmngr.NewPasswordHasher()
	checkForError(err)
	authMngr, err := authmngr.NewAuthMngr(dbRepo, pswHasher)
	checkForError(err)
	app := app.NewApp(articleReader, authMngr, ":8080")
	app.Run()
}

// Package app creates and runs chemistryBack
package app

import (
	"log"
	"net/http"

	"github.com/Teriton/chemistryBack/internal/handler"
	"github.com/Teriton/chemistryBack/pkg/articlemngr"
	"github.com/Teriton/chemistryBack/pkg/articlereader"
	"github.com/Teriton/chemistryBack/pkg/authmngr"
	"github.com/Teriton/chemistryBack/pkg/dbrepo"
)

type App struct {
	Server        *http.Server
	ArticleReader articlereader.ArticleReader
	AuthMngr      authmngr.AuthorizationMngr
	ArticleMngr   articlemngr.ArticleMngr
}

func NewApp(
	articleReader articlereader.ArticleReader,
	authMngr authmngr.AuthorizationMngr,
	dbRepo dbrepo.DBRepo,
	articleMngr articlemngr.ArticleMngr,
	addr string,
) *App {
	mux := http.NewServeMux()

	articleHandler := handler.NewArticlesHandler(authMngr, dbRepo, articleReader)
	authHandler, err := handler.NewAuthHandler(authMngr)
	if err != nil {
		panic("can't create auth handler")
	}
	userHandler, err := handler.NewUserHandler(authMngr, dbRepo)
	if err != nil {
		panic("can't create user handler")
	}
	questionHandler, err := handler.NewQuestionHandler(
		articleMngr,
		authMngr,
		dbRepo,
	)
	if err != nil {
		panic("can't create question handler")
	}

	mux.HandleFunc("GET /articles/list", articleHandler.ListArticlesWithCompletion)
	mux.HandleFunc("GET /articles/byPath/{path...}", articleHandler.GetArticle)

	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /signup", authHandler.Signup)
	mux.HandleFunc("POST /logout", authHandler.Logout)

	mux.HandleFunc("GET /user", userHandler.GetUser)
	mux.HandleFunc("GET /user/completedLessons", userHandler.GetUserWithCopletedLessosnCount)

	mux.HandleFunc("POST /complete", questionHandler.CompleteArticle)

	handler := CORS(mux)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return &App{server, articleReader, authMngr, articleMngr}
}

func (a *App) Run() error {
	log.Printf("Starting server on %s", a.Server.Addr)
	return a.Server.ListenAndServe()
}

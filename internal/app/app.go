// Package app creates and runs chemistryBack
package app

import (
	"log"
	"net/http"

	"github.com/Teriton/chemistryBack/internal/handler"
	"github.com/Teriton/chemistryBack/pkg/articlereader"
	"github.com/Teriton/chemistryBack/pkg/authmngr"
	"github.com/Teriton/chemistryBack/pkg/dbrepo"
)

type App struct {
	Server        *http.Server
	ArticleReader articlereader.ArticleReader
	AuthMngr      authmngr.AuthorizationMngr
}

func NewApp(articleReader articlereader.ArticleReader, authMngr authmngr.AuthorizationMngr, dbRepo dbrepo.DBRepo, addr string) *App {
	mux := http.NewServeMux()

	articleHandler := handler.NewArticlesHandler(articleReader)
	authHandler, err := handler.NewAuthHandler(authMngr)
	if err != nil {
		panic("can't create auth handler")
	}
	userHandler, err := handler.NewUserHandler(authMngr, dbRepo)
	if err != nil {
		panic("can't create user handler")
	}

	mux.HandleFunc("GET /articles/list", articleHandler.ListArticles)
	mux.HandleFunc("GET /articles/byPath/{path...}", articleHandler.GetArticle)

	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /signup", authHandler.Signup)
	mux.HandleFunc("POST /logout", authHandler.Logout)

	mux.HandleFunc("GET /user", userHandler.GetUser)

	handler := CORS(mux)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return &App{server, articleReader, authMngr}
}

func (a *App) Run() error {
	log.Printf("Starting server on %s", a.Server.Addr)
	return a.Server.ListenAndServe()
}

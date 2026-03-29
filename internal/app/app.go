// Package app creates and runs chemistryBack
package app

import (
	"log"
	"net/http"

	"github.com/Teriton/chemistryBack/internal/handler"
	"github.com/Teriton/chemistryBack/pkg/articlereader"
	"github.com/Teriton/chemistryBack/pkg/authmngr"
)

type App struct {
	Server        *http.Server
	ArticleReader articlereader.ArticleReader
	AuthMngr      authmngr.AuthorizationMngr
}

func NewApp(articleReader articlereader.ArticleReader, authMngr authmngr.AuthorizationMngr, addr string) *App {
	mux := http.NewServeMux()

	articleHandler := handler.NewArticlesHandler(articleReader)
	authHandler, err := handler.NewAuthHandler(authMngr)
	if err != nil {
		panic("cant create auth handler")
	}

	mux.HandleFunc("GET /articles/list", articleHandler.ListArticles)
	mux.HandleFunc("GET /articles/byPath/{path...}", articleHandler.GetArticle)

	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /signup", authHandler.Signup)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &App{server, articleReader, authMngr}
}

func (a *App) Run() error {
	log.Printf("Starting server on %s", a.Server.Addr)
	return a.Server.ListenAndServe()
}

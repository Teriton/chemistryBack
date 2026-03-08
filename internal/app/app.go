package app

import (
	"log"
	"net/http"

	"github.com/Teriton/chemistryBack/internal/handler"
	"github.com/Teriton/chemistryBack/pkg/articlereader"
)

type App struct {
	Server        *http.Server
	ArticleReader articlereader.ArticleReader
}

func NewApp(articlereader articlereader.ArticleReader, addr string) *App {
	mux := http.NewServeMux()

	articleHandler := handler.NewArticlesHandler(articlereader)

	mux.HandleFunc("GET /articles", articleHandler.ListArticles)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &App{server, articlereader}
}

func (a *App) Run() error {
	log.Printf("Starting server on %s", a.Server.Addr)
	return a.Server.ListenAndServe()
}

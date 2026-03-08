// Package handler handl all routes
package handler

import (
	"log"
	"net/http"

	"github.com/Teriton/chemistryBack/pkg/articlereader"
)

type ArticlesHandler struct {
	articleReader articlereader.ArticleReader
}

func NewArticlesHandler(articleReader articlereader.ArticleReader) *ArticlesHandler {
	return &ArticlesHandler{articleReader: articleReader}
}

func (h *ArticlesHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	chapter, err := h.articleReader.GetRootChapter()
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	jsonString, err := chapter.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonString)
}

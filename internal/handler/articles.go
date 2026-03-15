// Package handler handl all routes
package handler

import (
	"encoding/json"
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

func (h *ArticlesHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	pathToArticle := r.PathValue("path")
	article, err := h.articleReader.GetArticle(pathToArticle)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		jsonString, err := json.Marshal(map[string]string{"error": "Can't find article with this name"})
		if err != nil {
			log.Fatal(err)
		}
		w.Write(jsonString)
		return
	}
	jsonString, err := article.MarshalJSONWithContent()
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonString)
}

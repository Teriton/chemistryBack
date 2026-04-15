// Package handler handl all routes
package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Teriton/chemistryBack/pkg/articlereader"
	"github.com/Teriton/chemistryBack/pkg/authmngr"
	"github.com/Teriton/chemistryBack/pkg/dbrepo"
)

type ArticlesHandler struct {
	authMngr      authmngr.AuthorizationMngr
	dbRepo        dbrepo.DBRepo
	articleReader articlereader.ArticleReader
}

func NewArticlesHandler(
	authMngr authmngr.AuthorizationMngr,
	dbRepo dbrepo.DBRepo,
	articleReader articlereader.ArticleReader,
) *ArticlesHandler {
	return &ArticlesHandler{
		articleReader: articleReader,
		dbRepo:        dbRepo,
		authMngr:      authMngr,
	}
}

func (h *ArticlesHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] /articles/list")
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

func (h *ArticlesHandler) LessonsCompleted(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] /articles/lessonsCompleted")
	w.Header().Set("Content-Type", "application/json")
	cookies := r.CookiesNamed("token")

	var err error
	if len(cookies) < 1 {
		err = errors.New("cookie is not set")
	}
	if checkError(w, err, http.StatusForbidden) {
		return
	}

	jwtToken := cookies[0].Value
	jwtContent, err := h.authMngr.Verify(jwtToken)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	titles, err := h.dbRepo.GetCompletedLessonsForUser(jwtContent.UserID)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	response := map[string]any{
		"titles": titles,
	}

	jsonData, err := json.Marshal(response)
	if checkError(w, err, http.StatusForbidden) {
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonData)
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

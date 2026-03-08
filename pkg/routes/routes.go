// Package routes defines all routes in app
package routes

import "net/http"

func MakeChemistryMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /articles", listArticles)

	return mux
}

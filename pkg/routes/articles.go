package routes

import (
	"fmt"
	"net/http"
)

func listArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there %s", r.URL.Path[1:])
}

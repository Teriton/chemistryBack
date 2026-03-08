package main

import (
	"log"
	"net/http"

	"github.com/Teriton/chemistryBack/pkg/routes"
)

func main() {
	mux := routes.MakeChemistryMux()

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

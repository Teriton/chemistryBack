package main

import (
	"log"

	"github.com/Teriton/chemistryBack/pkg/articlereader"
)

func main() {
	dirReadr, err := articlereader.NewDirDeader("articles/test")
	if err != nil {
		log.Fatal(err)
	}
	dirReadr.GetArticles()
}

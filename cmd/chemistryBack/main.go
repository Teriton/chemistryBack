package main

import (
	"github.com/Teriton/chemistryBack/internal/app"
	"github.com/Teriton/chemistryBack/pkg/articlereader"
)

func main() {

	articleReader := articlereader.NewDirReader("articles/main")
	app := app.NewApp(articleReader, ":8080")
	app.Run()
}

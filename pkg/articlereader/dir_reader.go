// Package articlereader reads all articles
package articlereader

import (
	"log"
	"os"
	"strings"
)

type ArticleReader interface {
	GetRootChapter() (Chapter, error)
	GetArticle(string) (Article, error)
}

type DirReader struct {
	dir         string
	rootChapter Chapter
}

func NewDirReader(dir string) *DirReader {
	dirReader := DirReader{dir, readInternalDir(dir)}
	return &dirReader
}

func getArticleTitleFromPath(path string) string {
	splitedPath := strings.Split(path, "/")
	return strings.Split(splitedPath[len(splitedPath)-1], ".")[0]
}

func artticleFromFile(filePath string) (Article, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Article{}, err
	}
	defer file.Close()

	currentArticle := Article{}
	currentArticle.title = getArticleTitleFromPath(file.Name())
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return Article{}, err
	}
	currentArticle.content = string(fileContent)

	return currentArticle, nil
}

func readInternalDir(dir string) Chapter {
	currentChapter := Chapter{title: getArticleTitleFromPath(dir), articleObjects: []ArticleObject{}}
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.Type().IsDir() {
			currentChapter.articleObjects = append(currentChapter.articleObjects, readInternalDir(dir+"/"+file.Name()))
			continue
		}
		article, err := artticleFromFile(dir + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		currentChapter.articleObjects = append(currentChapter.articleObjects, article)
	}
	return currentChapter
}

func (dr DirReader) GetRootChapter() (Chapter, error) {
	return dr.rootChapter, nil
}

func (dr DirReader) GetArticle(path string) (Article, error) {
	article, err := dr.rootChapter.GetArticle(path)
	if err != nil {
		return Article{}, err
	}
	return article, nil
}

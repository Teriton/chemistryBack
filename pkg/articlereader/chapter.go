package articlereader

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Chapter struct {
	Title          string
	articleObjects []ArticleObject
}

func (c Chapter) ArticleObjects() []ArticleObject {
	return c.articleObjects
}

func printChapter(c Chapter) {
	fmt.Printf("\t%s\n", c.Title)
	for _, articleObject := range c.articleObjects {
		article, ok := articleObject.(Article)
		if ok {
			fmt.Println(article.Title)
		}
		if chapter, ok := articleObject.(Chapter); ok {
			printChapter(chapter)
		}
	}
}

func (c Chapter) PrintChapter() {
	printChapter(c)
}

func (c Chapter) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Title    string          `json:"title"`
		Articles []ArticleObject `json:"articles,omitempty"`
	}{
		Title:    c.Title,
		Articles: c.articleObjects,
	})
}

func getArticle(articleToFind []string, c Chapter) (Article, error) {
	for _, articleObject := range c.articleObjects {
		article, ok := articleObject.(Article)
		if ok && len(articleToFind) == 1 && article.Title == articleToFind[0] {
			return article, nil
		}
		if chapter, ok := articleObject.(Chapter); ok && chapter.Title == articleToFind[0] {
			if len(articleToFind) == 1 {
				return Article{Title: chapter.Title, Content: "CHAPTER"}, nil
			} else {
				article, err := getArticle(articleToFind[1:], chapter)
				if err != nil {
					return Article{}, err
				}
				return article, nil
			}
		}
	}
	return Article{}, errors.New("can't finde article")
}

func (c Chapter) GetArticle(path string) (Article, error) {
	splitPath := strings.Split(path, "/")
	article, err := getArticle(splitPath, c)
	if err != nil {
		return Article{}, err
	}
	return article, nil
}

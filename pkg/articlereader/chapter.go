package articlereader

import (
	"encoding/json"
	"fmt"
)

type Chapter struct {
	title          string
	articleObjects []ArticleObject
}

func (c Chapter) Title() string {
	return c.title
}

func (c Chapter) ArticleObjects() []ArticleObject {
	return c.articleObjects
}

func printChapter(c Chapter) {
	fmt.Printf("\t%s\n", c.title)
	for _, articleObject := range c.articleObjects {
		article, ok := articleObject.(Article)
		if ok {
			fmt.Println(article.Title())
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
	type Alias Chapter
	return json.Marshal(&struct {
		Title    string          `json:"title"`
		Articles []ArticleObject `json:"articles,omitempty"`
	}{
		Title:    c.title,
		Articles: c.articleObjects,
	})
}

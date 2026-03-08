package articlereader

import "fmt"

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

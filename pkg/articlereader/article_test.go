package articlereader

import (
	"fmt"
	"testing"
)

func createArticle() Article {
	title := "test"
	content := "content of test page"
	article := Article{title, content}
	return article
}

func TestArticle(t *testing.T) {
	title := "test"
	content := "content of test page"
	article := Article{title, content}
	if article.Title() != title {
		t.Errorf("Title from an article: %s, not equal to %s", article.Title(), title)
	}
	if article.Content() != content {
		t.Errorf("Title from an content: %s, not equal to %s", article.Content(), content)
	}
}

func TestJsonEncodeArticle(t *testing.T) {
	article := createArticle()
	jsonByte, err := article.MarshalJSON()
	check(err, t)
	jsonString := string(jsonByte)
	fmt.Println(jsonString)
}

func TestJsonEncodeArticleWithContext(t *testing.T) {
	article := createArticle()
	jsonByte, err := article.MarshalJSONWithContent()
	check(err, t)
	jsonString := string(jsonByte)
	fmt.Println(jsonString)
}

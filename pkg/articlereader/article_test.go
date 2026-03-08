package articlereader

import (
	"testing"
)

func TestArticle(t *testing.T) {
	title := "Test"
	content := "Content of test page"
	article := Article{title, content}
	if article.Title() != title {
		t.Errorf("Title from an article: %s, not equal to %s", article.Title(), title)
	}
	if article.Content() != content {
		t.Errorf("Title from an content: %s, not equal to %s", article.Content(), content)
	}
}

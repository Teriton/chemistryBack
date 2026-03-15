package articlereader

import "encoding/json"

type ArticleObject interface {
	Title() string
}

type Article struct {
	title   string
	content string
}

func (a Article) Title() string {
	return a.title
}

func (a *Article) Content() string {
	return a.content
}

func (a Article) MarshalJSONWithContent() ([]byte, error) {
	type Alias struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	return json.Marshal(&Alias{
		Title:   a.title,
		Content: a.content,
	})
}

func (a Article) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Title string `json:"title"`
	}

	return json.Marshal(&Alias{
		Title: a.title,
	})
}

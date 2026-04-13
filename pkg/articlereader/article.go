package articlereader

import "encoding/json"

type ArticleObject interface {
}

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (a Article) MarshalJSONWithContent() ([]byte, error) {
	type Alias struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	return json.Marshal(&Alias{
		Title:   a.Title,
		Content: a.Content,
	})
}

func (a Article) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Title string `json:"title"`
	}

	return json.Marshal(&Alias{
		Title: a.Title,
	})
}

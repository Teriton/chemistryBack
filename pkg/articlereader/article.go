package articlereader

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

func NewArticle(title string, content string) (*Article, error) {
	return &Article{title: title, content: content}, nil
}

package articlereader

import "testing"

func TestPrintChapter(t *testing.T) {
	chapter := Chapter{"Root", []ArticleObject{
		Article{"Hello World", ""},
		Article{"Main", ""},
		Chapter{"Chapter 1", []ArticleObject{
			Article{"Hello World2", ""},
			Article{"Hello World3", ""},
		}},
	},
	}
	chapter.PrintChapter()
}

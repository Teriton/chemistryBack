package articlereader

import (
	"fmt"
	"testing"
)

func createChapter() Chapter {
	chapter := Chapter{"Root", []ArticleObject{
		Article{"Hello World", ""},
		Article{"Main", ""},
		Chapter{"Chapter 1", []ArticleObject{
			Article{"Hello World2", ""},
			Article{"Hello World3", ""},
		}},
	},
	}
	return chapter
}

func TestPrintChapter(t *testing.T) {
	chapter := createChapter()
	chapter.PrintChapter()
}

func TestJsonEncodeChapter(t *testing.T) {
	chapter := createChapter()
	jsonByte, err := chapter.MarshalJSON()
	check(err, t)
	jsonString := string(jsonByte)
	fmt.Println(jsonString)
}

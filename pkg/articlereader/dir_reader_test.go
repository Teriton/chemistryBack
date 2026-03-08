package articlereader

import (
	"testing"
)

func check(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}
func newDirReaderForTest() *DirReader {
	dirReader := NewDirReader("../../articles/test")
	return dirReader
}

func TestArticleFromFile(t *testing.T) {
	arcticle, err := artticleFromFile("../../articles/test/Test.md")
	check(err, t)
	if arcticle.Title() != "Test" {
		t.Errorf("ArticleFromFile result: %#v\n", arcticle)
	}
}

func TestGetRootChapter(t *testing.T) {
	dirReader := newDirReaderForTest()
	chapter, err := dirReader.GetRootChapter()
	check(err, t)
	if len(chapter.articleObjects) == 0 {
		t.Error("There are no articles")
		t.Skip()
	}
	chapter.PrintChapter()
}

package articlereader

import (
	"errors"
	"testing"
)

func check(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

func newDirReaderForTest() (*DirReader, error) {
	dirReader, err := NewDirReader("../../articles/test")
	if err != nil {
		return nil, errors.New("Can't create new DirReader")
	}
	return dirReader, nil
}

func TestArticleFromFile(t *testing.T) {
	arcticle, err := artticleFromFile("../../articles/test/Test.md")
	check(err, t)
	if arcticle.Title() != "Test" {
		t.Errorf("ArticleFromFile result: %#v\n", arcticle)
	}
}

func TestGetRootChapter(t *testing.T) {
	dirReader, err := newDirReaderForTest()
	check(err, t)
	chapter := dirReader.GetRootChapter()
	check(err, t)
	if len(chapter.articleObjects) == 0 {
		t.Error("There are no articles")
		t.Skip()
	}
	chapter.PrintChapter()
}

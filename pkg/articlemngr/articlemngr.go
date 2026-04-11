// Package articlemngr manipulate questions
package articlemngr

import (
	"github.com/Teriton/chemistryBack/pkg/dbrepo"
)

type ArticleMngr struct {
	dbrepo dbrepo.DBRepo
}

func NewArticleMngr(dbRepo dbrepo.DBRepo) (*ArticleMngr, error) {
	return &ArticleMngr{dbRepo}, nil
}

func (at ArticleMngr) CompleteLesson(username string, lessonTitle string, xp int) error {
	user, err := at.dbrepo.GetUserByUserName(username)
	if err != nil {
		return err
	}
	lesson, err := at.dbrepo.GetLessonByTitle(lessonTitle)
	if err != nil {
		if err.Error() == "no rows in result set" {
			err = at.dbrepo.CreateLessonWithTitle(lessonTitle)
			if err != nil {
				return err
			}
			lesson, err = at.dbrepo.GetLessonByTitle(lessonTitle)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	err = at.dbrepo.CreateCompletedLesson(user.ID, lesson.ID)
	if err != nil {
		return err
	}
	err = at.dbrepo.AddXPToUser(user.ID, xp)
	if err != nil {
		return err
	}
	return nil
}

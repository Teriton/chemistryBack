// Package dbrepo managets db connections
package dbrepo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Teriton/chemistryBack/internal/models"
)

type DBRepo interface {
	CreateUser(models.AddUser) error
	DeleteUserByUserName(string) error
	GetUserByUserName(string) (models.User, error)
	GetUserByID(int) (models.User, error)
	AddXPToUser(int, int) error

	GetLessonByTitle(string) (models.Lesson, error)
	CreateLessonWithTitle(string) error
	DeleteLessonByTitle(string) error

	CreateCompletedLesson(int, int) error
	DeleteCompletedLesson(int, int) error
	GetCompletedLessonsLenForUser(int) (int, error)

	CloseDB() error
}

type PsqlRepo struct {
	dbURL  string
	dbPool *pgxpool.Pool
}

func NewPsqlRepo(URL string) (*PsqlRepo, error) {
	var psqlRepo PsqlRepo
	psqlRepo.dbURL = URL
	dbpool, err := pgxpool.New(context.Background(), URL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}
	psqlRepo.dbPool = dbpool
	return &psqlRepo, nil
}

func (pr PsqlRepo) CreateUser(user models.AddUser) error {
	tx, err := pr.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO users(email, password, username) VALUES ($1, $2, $3)",
		user.Email, user.Password, user.Username,
	)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (pr PsqlRepo) DeleteUserByUserName(username string) error {
	tx, err := pr.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		context.Background(),
		"delete from users where username=$1",
		username,
	)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (pr PsqlRepo) GetUserByUserName(username string) (models.User, error) {
	var user models.User
	err := pr.dbPool.QueryRow(
		context.Background(),
		"select * from users where username = $1", username).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Username,
		&user.Xp,
		&user.Streak,
		&user.CreationDate,
	)

	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (pr *PsqlRepo) CloseDB() error {
	if pr.dbPool != nil {
		pr.dbPool.Close()
	}
	return nil
}
func (pr PsqlRepo) GetLessonByTitle(title string) (models.Lesson, error) {
	var lesson models.Lesson
	err := pr.dbPool.QueryRow(
		context.Background(),
		"select * from lessons where title = $1", title).Scan(
		&lesson.ID,
		&lesson.Title,
	)

	if err != nil {
		return models.Lesson{}, err
	}
	return lesson, nil
}
func (pr PsqlRepo) CreateLessonWithTitle(title string) error {
	tx, err := pr.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO lessons(title) VALUES ($1)",
		title,
	)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (pr PsqlRepo) DeleteLessonByTitle(title string) error {
	tx, err := pr.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		context.Background(),
		"delete from lessons where title=$1",
		title,
	)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (pr PsqlRepo) CreateCompletedLesson(userID int, lessonID int) error {
	tx, err := pr.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO lessons_completed(user_id, lesson_id) VALUES ($1, $2)",
		userID, lessonID,
	)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil

}

func (pr PsqlRepo) DeleteCompletedLesson(userID int, lessonID int) error {
	tx, err := pr.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		context.Background(),
		"delete from lessons_completed where user_id=$1 AND lesson_id=$2",
		userID, lessonID,
	)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (pr PsqlRepo) GetCompletedLessonsLenForUser(userID int) (int, error) {
	var complitedLessonsCount int
	err := pr.dbPool.QueryRow(
		context.Background(),
		"select COUNT(*) from lessons_completed where user_id = $1", userID).Scan(
		&complitedLessonsCount,
	)

	if err != nil {
		return -1, err
	}
	return complitedLessonsCount, nil
}

func (pr PsqlRepo) GetUserByID(userID int) (models.User, error) {
	var user models.User
	err := pr.dbPool.QueryRow(
		context.Background(),
		"select * from users where id = $1", userID).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Username,
		&user.Xp,
		&user.Streak,
		&user.CreationDate,
	)

	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (pr PsqlRepo) AddXPToUser(userID int, xp int) error {
	tx, err := pr.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		context.Background(),
		"UPDATE users SET xp = xp + $1 WHERE id = $2",
		xp, userID,
	)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil

}

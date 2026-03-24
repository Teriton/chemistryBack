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
		&user.Xp)
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

package dbrepo

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Teriton/chemistryBack/internal/models"
)

func check(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

func getConnection() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		return nil, fmt.Errorf("Unable to create connection pool: %v\n", err)
	}
	return dbpool, nil
}

func TestConnectionToDB(t *testing.T) {
	dbpool, err := getConnection()
	if err != nil {
		t.Error(err)
	}
	defer dbpool.Close()

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		t.Errorf("QueryRow failed: %v\n", err)
	}

	fmt.Println(greeting)
}
func TestSelectionData(t *testing.T) {
	dbpool, err := getConnection()
	if err != nil {
		t.Error(err)
	}
	defer dbpool.Close()

	var users []models.User
	var user models.User
	rows, err := dbpool.Query(context.Background(), "select * from users")
	check(err, t)
	_, err = pgx.ForEachRow(rows, []any{&user.ID, &user.Email, &user.Password, &user.Username, &user.Xp}, func() error {
		users = append(users, user)
		return nil
	})
	check(err, t)

	fmt.Println(users)
}

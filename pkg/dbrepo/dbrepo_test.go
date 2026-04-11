package dbrepo

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
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

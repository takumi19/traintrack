package database

import (
	"context"
	"errors"

	_ "github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

var ErrNotFound = errors.New("not found")

// Returns a DB connected to the provided URL. On failure returns nil and an error.
func New(dbUrl string, automigrate bool) (*DB, error) {
	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	if automigrate {
	}

	return &DB{
		Pool: conn,
	}, nil
}

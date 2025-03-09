package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	db *pgxpool.Pool
}

// Returns a DB connected to the provided URL. On failure returns nil and an error.
func New(dbUrl string) (*DB, error) {
	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &DB{
		db: conn,
	}, nil
}

func (s *DB) Close() error {
	s.db.Close()
	return nil
}

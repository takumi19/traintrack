package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type User struct {
	Id           int     `json:"id" db:"id"`
	FullName     *string `json:"full_name" db:"full_name"`
	Login        *string `json:"login" db:"login"`
	Email        *string `json:"email" db:"email"`
	PasswordHash *string `json:"password_hash" db:"password_hash"`
}

func (s *DB) CreateUser(user *User) (int64, error) {
	var id int64
	err := s.db.QueryRow(context.Background(), `
INSERT INTO users (
  full_name, login, email, password_hash
)
VALUES ($1, $2, $3, $4)
RETURNING id`, &user.FullName, &user.Login, &user.Email, &user.PasswordHash).Scan(&id)

	return id, err
}

// TODO: Get the user by login not by the id
func (s *DB) ReadUser(id int64) (*User, error) {
	rows, _ := s.db.Query(context.Background(), "SELECT * FROM users WHERE id=$1", id)

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil && err == pgx.ErrNoRows {
		return nil, err
	}

	return &user, nil
}

func (s *DB) UpdateUser(user *User) error {
	tx, err := s.db.Begin(context.Background())
	defer tx.Rollback(context.Background()) // Safe to call even after commit
	if err != nil {
		log.Default().Println("Failed to start transaction:", err)
		return err
	}

	// NOTE: COALESCE(NULLIF(value, '') will only update the cell if the value if not null
	_, err = tx.Exec(context.Background(), `
UPDATE users
SET full_name     = COALESCE(NULLIF($1, ''), full_name),
    login         = COALESCE(NULLIF($2, ''), login),
    email         = COALESCE(NULLIF($3, ''), email),
    password_hash = COALESCE(NULLIF($4, ''), password_hash)
WHERE
    id=$5`, user.FullName, user.Login, user.Email, user.PasswordHash, user.Id)
	if err != nil {
		return err
	}

	if err = tx.Commit(context.Background()); err != nil {
		log.Default().Println("Error when committing transaction:", err)
		return err
	}

	return nil
}

func (s *DB) DeleteUser(id int64) error {
	if cmd_tag, err := s.db.Exec(context.Background(), "DELETE FROM users WHERE id=$1", id); err != nil {
		log.Default().Println(cmd_tag, err)
		return err
	}

	return nil
}

func (s *DB) ListUsers() ([]User, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		log.Default().Println("Failed to retrieve all users")
		rows.Close()
		return nil, err
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		return nil, err
	}

	return users, nil
}

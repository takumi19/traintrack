package database

import "time"

type Chat struct {
	Id        int64     `json:"id" db:"id"`
	Name      *string    `json:"name" db:"name"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

type Message struct {
	Id        int64     `json:"id" db:"id"`
	Name      *string    `json:"name" db:"name"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

func (db *DB) GetChatsByUserId(id int64) ([]Chat, error) {
	return nil, nil
}

func (db *DB) GetUsersPersonalchat(userId1 int64, userId2 int64) ([]Chat, error) {
	return nil, nil
}

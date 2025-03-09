package database

import "time"

type Log struct {
	Id          int        `json:"id" db:"id"`
	UserId      int        `json:"user_id" db:"user_id"`
	WorkoutDate *string    `json:"workout_date" db:"workout_date"`
	Notes       *string    `json:"notes" db:"notes"`
	CreatedAt   *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at" db:"updated_at"`
}

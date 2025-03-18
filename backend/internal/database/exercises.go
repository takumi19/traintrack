package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ExerciseInfo struct {
	Id           int       `json:"id" db:"id"`
	Name         int       `json:"name" db:"name"`
	Notes        int       `json:"notes" db:"notes"`
	IsRepBased   bool      `json:"is_rep_based" db:"is_rep_based"`
	IsBodyweight bool      `json:"is_bodyweight" db:"is_bodyweight"`
	MuscleGroups []*string `json:"muscle_groups" db:"-"`
}

func (db *DB) ListExercises() ([]ExerciseInfo, error) {
	rows, err := db.Query(context.Background(), "SELECT * FROM exercises")
	if err != nil {
		return nil, err
	}

	exercises, err := pgx.CollectRows(rows, pgx.RowToStructByName[ExerciseInfo])
	if err != nil {
		return nil, err
	}
	return exercises, nil
}

func (db *DB) GetExerciseByID(id int64) (*ExerciseInfo, error) {
	rows, err := db.Query(context.Background(), "SELECT * FROM exercises WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	exercises, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[ExerciseInfo])
	if err != nil {
		return nil, err
	}
	return &exercises, nil
}

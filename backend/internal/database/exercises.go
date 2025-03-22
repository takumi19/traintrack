package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type ExerciseInfo struct {
	Id           int64    `json:"id" db:"id"`
	Name         string   `json:"name" db:"name"`
	Notes        *string  `json:"notes" db:"notes"`
	IsRepBased   bool     `json:"is_rep_based" db:"is_rep_based"`
	IsBodyweight bool     `json:"is_bodyweight" db:"is_bodyweight"`
	MuscleGroups []string `json:"muscle_groups" db:"-"`
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

	for i, exercise := range exercises {
		rows, err = db.Query(context.Background(), "SELECT (worked_muscle_group) FROM exercises_muscle_groups WHERE exercise_id=$1", exercise.Id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
        continue
			}
			return nil, err
		}

		exercises[i].MuscleGroups, err = pgx.CollectRows(rows, pgx.RowTo[string])
		if err != nil {
			return nil, err
		}
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

func (db *DB) GetExerciseByName(name string) (*ExerciseInfo, error) {
	rows, err := db.Query(context.Background(), "SELECT * FROM exercises WHERE name=$1", name)
	if err != nil {
		return nil, err
	}

	exercises, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[ExerciseInfo])
	if err != nil {
		return nil, err
	}
	return &exercises, nil
}

func (db *DB) AddExercise(e *ExerciseInfo) (int64, error) {
	err := db.QueryRow(context.Background(), `
INSERT INTO exercises (
  name, notes, is_rep_based, is_bodyweight
)
VALUES ($1, $2, $3, $4)
RETURNING id`, e.Name, e.Notes, e.IsRepBased, e.IsBodyweight).Scan(&e.Id)
	if err != nil {
		return 0, err
	}

	for _, muscleGroup := range e.MuscleGroups {
		_, err = db.Exec(context.Background(), `
INSERT INTO exercises_muscle_groups (
  exercise_id, worked_muscle_group
)
VALUES ($1, $2)`, e.Id, muscleGroup)
		if err != nil {
			break
		}
	}

	return e.Id, err
}

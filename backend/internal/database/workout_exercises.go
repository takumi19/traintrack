package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type WorkoutExercise struct {
	Id               int64      `json:"id" db:"id"`
	ProgramWorkoutId *int64     `json:"program_workout_id" db:"program_workout_id"`
	ExerciseId       *int64     `json:"exercise_id" db:"exercise_id"`
	OrderIndex       *int32     `json:"order_index" db:"order_index"`
	Notes            *string    `json:"notes" db:"notes"`
	CreatedAt        *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at" db:"updated_at"`

	Sets []WorkoutSet `json:"sets" db:"-"`
}

func (db *DB) ListWorkoutExercises(workoutId int64) ([]WorkoutExercise, error) {
	rows, err := db.Query(context.Background(), `
SELECT * FROM program_workout_exercises
WHERE program_workout_id=$1
ORDER BY order_index`, workoutId)
	if err != nil {
		return nil, err
	}

	exercises, err := pgx.CollectRows(rows, pgx.RowToStructByName[WorkoutExercise])
	if err != nil {
		return nil, err
	}

	for i := range exercises {
		sets, err := db.ListExerciseSets(int64(exercises[i].Id))
		if sets == nil {
			log.Default().Println("getWorkoutExercises:", err)
			continue
		}
		exercises[i].Sets = sets
	}

	return exercises, nil
}

func (db *DB) UpdateWorkoutExercise(exercise *WorkoutExercise) error {
	var err error = nil
	_, err = db.Exec(context.Background(), `
UPDATE program_workout_exercises
SET program_workout_id = COALESCE($1, program_workout_id),
    exercise_id        = COALESCE($2, exercise_id),
    order_index        = COALESCE($3, order_index),
    notes              = COALESCE($4, notes)
WHERE id=$5`, exercise.ProgramWorkoutId, exercise.ExerciseId, exercise.OrderIndex, exercise.Notes, exercise.Id)
	if err != nil {
		return nil
	}

	for i := range exercise.Sets {
		if err := db.UpdateWorkoutSet(&exercise.Sets[i]); err != nil {
			fmt.Println("Failed to update workout sets:", err)
			return err
		}
	}

	return err
}

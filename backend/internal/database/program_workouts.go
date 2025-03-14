package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type ProgramWorkout struct {
	Id            int64      `json:"id" db:"id"`
	ProgramWeekId int64      `json:"program_week_id" db:"program_week_id"`
	WorkoutIndex  *int32     `json:"workout_index" db:"workout_index"`
	Title         *string    `json:"title" db:"title"`
	Notes         *string    `json:"notes" db:"notes"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`

	Exercises []WorkoutExercise `json:"exercises" db:"-"`
}

func (db *DB) ListProgramWeekWorkouts(weekId int64) ([]ProgramWorkout, error) {
	rows, err := db.Query(context.Background(), `
SELECT * FROM program_workouts
WHERE program_week_id=$1
ORDER BY workout_index`, weekId)
	if err != nil {
		return nil, err
	}

	workouts, err := pgx.CollectRows(rows, pgx.RowToStructByName[ProgramWorkout])
	if err != nil {
		return nil, err
	}

	for i := range workouts {
		exercises, err := db.ListWorkoutExercises(int64(workouts[i].Id))
		if exercises == nil {
			log.Default().Println("getProgramWeekWorkouts:", err)
			continue
		}
		workouts[i].Exercises = exercises
	}

	return workouts, nil
}

func (db *DB) UpdateProgramWorkout(workout *ProgramWorkout) error {
	var err error = nil
	_, err = db.Exec(context.Background(), `
UPDATE program_workouts
SET workout_index = COALESCE($1, workout_index),
    title         = COALESCE($2, title),
    notes         = COALESCE($3, notes)
WHERE id=$4`, workout.WorkoutIndex, workout.Title, workout.Notes, workout.Id)
	if err != nil {
		return err
	}

	for i := range workout.Exercises {
		if err := db.UpdateWorkoutExercise(&workout.Exercises[i]); err != nil {
			fmt.Println("Failed to update workout exercises:", err)
			fmt.Printf("%+v\n", workout.Exercises[i])
			return err
		}
	}

	return err
}

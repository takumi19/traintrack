package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type WorkoutSet struct {
	Id                       int64  `json:"id" db:"id"`
	ProgramWorkoutExerciseId *int64 `json:"program_workout_exercise_id" db:"program_workout_exercise_id"`
	SetNumber                *int32 `json:"set_number" db:"set_number"`
	Rpe                      *int32 `json:"rpe" db:"rpe"`

	SuggestedRepsMin *int32 `json:"suggested_reps_min" db:"suggested_reps_min"`
	SuggestedRepsMax *int32 `json:"suggested_reps_max" db:"suggested_reps_max"`
	SuggestedReps    *int32 `json:"suggested_reps" db:"suggested_reps"`

	SuggestedWeightMin *float32 `json:"suggested_weight_min" db:"suggested_weight_min"`
	SuggestedWeightMax *float32 `json:"suggested_weight_max" db:"suggested_weight_max"`
	SuggestedWeight    *float32 `json:"suggested_weight" db:"suggested_weight"`

	SuggestedTimeMin *int32 `json:"suggested_time_min" db:"suggested_time_min"`
	SuggestedTimeMax *int32 `json:"suggested_time_max" db:"suggested_time_max"`
	SuggestedTime    *int32 `json:"suggested_time" db:"suggested_time"`

	Notes     *string    `json:"notes" db:"notes"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

func (db *DB) ListExerciseSets(exerciseId int64) ([]WorkoutSet, error) {
	rows, err := db.Query(context.Background(), `
SELECT * FROM program_workout_sets
WHERE program_workout_exercise_id=$1
ORDER BY set_number`, exerciseId)
	if err != nil {
		return nil, err
	}

	sets, err := pgx.CollectRows(rows, pgx.RowToStructByName[WorkoutSet])
	if err != nil {
		return nil, err
	}
	return sets, nil
}

func (db *DB) UpdateWorkoutSet(set *WorkoutSet) error {
	var err error = nil
	_, err = db.Exec(context.Background(), `
UPDATE program_workout_sets
SET program_workout_exercise_id = COALESCE($1, program_workout_exercise_id),
    set_number                  = COALESCE($2, set_number),
    rpe                         = COALESCE($3, rpe),
    suggested_reps_min          = COALESCE($4, suggested_reps_min),
    suggested_reps_max          = COALESCE($5, suggested_reps_max),
    suggested_reps              = COALESCE($6, suggested_reps),
    suggested_weight_min        = COALESCE($7, suggested_weight_min),
    suggested_weight_max        = COALESCE($8, suggested_weight_max),
    suggested_weight            = COALESCE($9, suggested_weight),
    suggested_time_min          = COALESCE($10, suggested_time_min),
    suggested_time_max          = COALESCE($11, suggested_time_max),
    suggested_time              = COALESCE($12, suggested_time),
    notes                       = COALESCE($13, notes)
WHERE id=$14`,
		set.ProgramWorkoutExerciseId,
		set.SetNumber,
		set.Rpe,
		set.SuggestedRepsMin,
		set.SuggestedRepsMax,
		set.SuggestedReps,
		set.SuggestedWeightMin,
		set.SuggestedWeightMax,
		set.SuggestedWeight,
		set.SuggestedTimeMin,
		set.SuggestedTimeMax,
		set.SuggestedTime,
		set.Notes,
		set.Id)

	return err
}

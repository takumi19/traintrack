package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type Program struct {
	Id        int64      `json:"id" db:"id"`
	AuthorId  int64      `json:"author_id" db:"author_id"`
	Name      *string    `json:"name" db:"name"`
	Notes     *string    `json:"notes" db:"notes"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	// Is it needed?
	Weeks []ProgramWeek `json:"program_weeks" db:"-"`
}

type ProgramWeek struct {
	Id                int64      `json:"id" db:"id"`
	ProgramTemplateId int64      `json:"program_template_id" db:"program_template_id"`
	WeekNumber        *int32     `json:"week_number" db:"week_number"`
	Notes             *string    `json:"notes" db:"notes"`
	CreatedAt         *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at" db:"updated_at"`

	Workouts []ProgramWorkout `json:"workouts" db:"-"`
}

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

// Should i use go context to pass the program id?
func (s *DB) GetProgramByID(id int64) (*Program, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM program_templates WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	program, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Program])
	if err != nil {
		return nil, err
	}

	return &program, nil
}

func (s *DB) ListPrograms() ([]Program, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM program_templates")
	if err != nil {
		return nil, err
	}

	programs, err := pgx.CollectRows(rows, pgx.RowToStructByName[Program])
	if err != nil {
		return nil, err
	}

	for i := range programs {
		weeks, err := s.ListProgramWeeks(int64(programs[i].Id))
		if weeks == nil {
			log.Default().Println("ListPrograms:", err)
			continue
		}
		programs[i].Weeks = weeks
	}

	return programs, nil
}

func (s *DB) ListProgramWeeks(programId int64) ([]ProgramWeek, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM program_weeks WHERE program_template_id=$1", programId)
	if err != nil {
		return nil, err
	}

	weeks, err := pgx.CollectRows(rows, pgx.RowToStructByName[ProgramWeek])
	if err != nil {
		return nil, err
	}

	for i := range weeks {
		workouts, err := s.ListProgramWeekWorkouts(int64(weeks[i].Id))
		if workouts == nil {
			log.Default().Println("getProgramWeeks:", err)
			continue
		}
		weeks[i].Workouts = workouts
	}

	return weeks, nil
}

func (s *DB) ListProgramWeekWorkouts(weekId int64) ([]ProgramWorkout, error) {
	rows, err := s.db.Query(context.Background(), `
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
		exercises, err := s.ListWorkoutExercises(int64(workouts[i].Id))
		if exercises == nil {
			log.Default().Println("getProgramWeekWorkouts:", err)
			continue
		}
		workouts[i].Exercises = exercises
	}

	return workouts, nil
}

func (s *DB) ListWorkoutExercises(workoutId int64) ([]WorkoutExercise, error) {
	rows, err := s.db.Query(context.Background(), `
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
		sets, err := s.ListExerciseSets(int64(exercises[i].Id))
		if sets == nil {
			log.Default().Println("getWorkoutExercises:", err)
			continue
		}
		exercises[i].Sets = sets
	}

	return exercises, nil
}

func (s *DB) ListExerciseSets(exerciseId int64) ([]WorkoutSet, error) {
	rows, err := s.db.Query(context.Background(), `
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

func (s *DB) UpdateProgram(program *Program) error {
	tx, err := s.db.Begin(context.Background())
	defer tx.Rollback(context.Background()) // Safe to call even after commit
	if err != nil {
		log.Default().Println("Failed to start transaction:", err)
		return err
	}

	_, err = tx.Exec(context.Background(), `
UPDATE program_templates
SET name     = COALESCE($1, name),
    notes    = COALESCE($2, notes)
WHERE
    id=$3`, program.Name, program.Notes, program.Id)
	if err != nil {
		log.Default().Println("Failed to update program notes or name: ", err)
		return err
	}

	for i := range program.Weeks {
		if err := s.UpdateProgramWeek(&program.Weeks[i], tx); err != nil {
			fmt.Println("Failed to update program weeks:", err)
			return err
		}
	}

	if err = tx.Commit(context.Background()); err != nil {
		log.Default().Println("Error when committing transaction:", err)
		return err
	}

	return nil
}

func (s *DB) UpdateProgramWeek(week *ProgramWeek, tx pgx.Tx) error {
	var err error = nil
	_, err = tx.Exec(context.Background(), `
UPDATE program_weeks
SET notes       = COALESCE($1, notes),
    week_number = COALESCE($2, week_number)
WHERE id=$3`, week.Notes, week.WeekNumber, week.Id)
	if err != nil {
		return err
	}

	for i := range week.Workouts {
		if err := s.UpdateProgramWorkout(&week.Workouts[i], tx); err != nil {
			fmt.Println("Failed to update program week workouts:", err)
			return err
		}
	}

	return err
}

func (s *DB) UpdateProgramWorkout(workout *ProgramWorkout, tx pgx.Tx) error {
	var err error = nil
	_, err = tx.Exec(context.Background(), `
UPDATE program_workouts
SET workout_index = COALESCE($1, workout_index),
    title         = COALESCE($2, title),
    notes         = COALESCE($3, notes)
WHERE id=$4`, workout.WorkoutIndex, workout.Title, workout.Notes, workout.Id)
	if err != nil {
		return err
	}

	for i := range workout.Exercises {
		if err := s.UpdateWorkoutExercise(&workout.Exercises[i], tx); err != nil {
			fmt.Println("Failed to update workout exercises:", err)
			fmt.Printf("%+v\n", workout.Exercises[i])
			return err
		}
	}

	return err
}

func (s *DB) UpdateWorkoutExercise(exercise *WorkoutExercise, tx pgx.Tx) error {
	var err error = nil
	_, err = tx.Exec(context.Background(), `
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
		if err := s.UpdateWorkoutSet(&exercise.Sets[i], tx); err != nil {
			fmt.Println("Failed to update workout sets:", err)
			return err
		}
	}

	return err
}

func (s *DB) UpdateWorkoutSet(set *WorkoutSet, tx pgx.Tx) error {
	var err error = nil
	_, err = tx.Exec(context.Background(), `
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

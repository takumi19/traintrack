package main

import "time"

type User struct {
	Id           int     `json:"id" db:"id"`
	FullName     *string `json:"full_name" db:"full_name"`
	Login        *string `json:"login" db:"login"`
	Email        *string `json:"email" db:"email"`
	PasswordHash *string `json:"password_hash" db:"password_hash"`
}

type ExerciseInfo struct {
	Id           int       `json:"id" db:"id"`
	Name         int       `json:"name" db:"name"`
	Notes        int       `json:"notes" db:"notes"`
	IsRepBased   bool      `json:"is_rep_based" db:"is_rep_based"`
	IsBodyweight bool      `json:"is_bodyweight" db:"is_bodyweight"`
	MuscleGroups []*string `json:"muscle_groups" db:"-"`
}

type Log struct {
	Id          int        `json:"id" db:"id"`
	UserId      int        `json:"user_id" db:"user_id"`
	WorkoutDate *string    `json:"workout_date" db:"workout_date"`
	Notes       *string    `json:"notes" db:"notes"`
	CreatedAt   *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at" db:"updated_at"`
}

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
	WeekNumber        int32      `json:"week_number" db:"week_number"`
	Notes             *string    `json:"notes" db:"notes"`
	CreatedAt         *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at" db:"updated_at"`

	Workouts []ProgramWorkout `json:"workouts" db:"-"`
}

type ProgramWorkout struct {
	Id            int64      `json:"id" db:"id"`
	ProgramWeekId int64      `json:"program_week_id" db:"program_week_id"`
	WorkoutIndex  int32      `json:"workout_index" db:"workout_index"`
	Title         *string    `json:"title" db:"title"`
	Notes         *string    `json:"notes" db:"notes"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`

	Exercises []WorkoutExercise `json:"exercises" db:"-"`
}

type WorkoutExercise struct {
	Id               int64      `json:"id" db:"id"`
	ProgramWorkoutId int64      `json:"program_workout_id" db:"program_workout_id"`
	ExerciseId       int64     `json:"exercise_id" db:"exercise_id"`
	OrderIndex       int32      `json:"order_index" db:"order_index"`
	Notes            *string    `json:"notes" db:"notes"`
	CreatedAt        *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at" db:"updated_at"`

	Sets []WorkoutSet `json:"sets" db:"-"`
}

type WorkoutSet struct {
	Id                       int64  `json:"id" db:"id"`
	ProgramWorkoutExerciseId int64  `json:"program_workout_exercise_id" db:"program_workout_exercise_id"`
	SetNumber                int32  `json:"set_number" db:"set_number"`
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

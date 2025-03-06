package main

import (
	"context"
	"log"

	// "log"
	// "net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Alternative implementation would be for the Storage
// to have methods for retrieving tables/schemas by name
// and have a table interface for working with the individual
// tables.
//
// Each table would have methods for working with its particular data.
//
// This would probably require each table to keeep a conection
// to the db, which is redundant.

type Storage interface {
	// User functions
	CreateUser(*User) (int64, error)
	ReadUser(int64) (*User, error)
	UpdateUser(*User) error
	DeleteUser(int64) error
	ListUsers() ([]User, error)

	// Program functions
	ListPrograms() ([]Program, error)
	ListProgramWeeks(int64) ([]ProgramWeek, error)
	ListProgramWeekWorkouts(int64) ([]ProgramWorkout, error)
	ListWorkoutExercises(int64) ([]WorkoutExercise, error)
	ListExerciseSets(int64) ([]WorkoutSet, error)
	// Closes connection to the database if needed
	Close() error
}

type PostgresStorage struct {
	db *pgxpool.Pool
}

// Returns a PostgresStorage connected to the provided URL. On failure returns nil and an error.
func MakePostgres(dbUrl string) (*PostgresStorage, error) {
	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db: conn,
	}, nil
}

func (s *PostgresStorage) Close() error {
	s.db.Close()
	return nil
}

func (s *PostgresStorage) CreateUser(user *User) (int64, error) {
	var id int64
	err := s.db.QueryRow(context.Background(), `
INSERT INTO users (
  full_name, login, email, password_hash
)
VALUES ($1, $2, $3, $4)
RETURNING id`, &user.FullName, &user.Login, &user.Email, &user.PasswordHash).Scan(&id)

	return id, err
}

func (s *PostgresStorage) ReadUser(id int64) (*User, error) {
	rows, _ := s.db.Query(context.Background(), "SELECT * FROM users WHERE id=$1", id)

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil && err == pgx.ErrNoRows {
		return nil, err
	}

	return &user, nil
}

func (s *PostgresStorage) UpdateUser(user *User) error {
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

func (s *PostgresStorage) DeleteUser(id int64) error {
	if cmd_tag, err := s.db.Exec(context.Background(), "DELETE FROM users WHERE id=$1", id); err != nil {
		log.Default().Println(cmd_tag, err)
		return err
	}

	return nil
}

func (s *PostgresStorage) ListUsers() ([]User, error) {
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

func (s *PostgresStorage) ListPrograms() ([]Program, error) {
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

func (s *PostgresStorage) ListProgramWeeks(programId int64) ([]ProgramWeek, error) {
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

func (s *PostgresStorage) ListProgramWeekWorkouts(weekId int64) ([]ProgramWorkout, error) {
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

func (s *PostgresStorage) ListWorkoutExercises(workoutId int64) ([]WorkoutExercise, error) {
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

func (s *PostgresStorage) ListExerciseSets(exerciseId int64) ([]WorkoutSet, error) {
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

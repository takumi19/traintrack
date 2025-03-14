package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type ProgramWeek struct {
	Id                int64      `json:"id" db:"id"`
	ProgramTemplateId int64      `json:"program_template_id" db:"program_template_id"`
	WeekNumber        *int32     `json:"week_number" db:"week_number"`
	Notes             *string    `json:"notes" db:"notes"`
	CreatedAt         *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at" db:"updated_at"`

	Workouts []ProgramWorkout `json:"workouts" db:"-"`
}

func (db *DB) ListProgramWeeks(programId int64) ([]ProgramWeek, error) {
	rows, err := db.Query(context.Background(), "SELECT * FROM program_weeks WHERE program_template_id=$1", programId)
	if err != nil {
		return nil, err
	}

	weeks, err := pgx.CollectRows(rows, pgx.RowToStructByName[ProgramWeek])
	if err != nil {
		return nil, err
	}

	for i := range weeks {
		workouts, err := db.ListProgramWeekWorkouts(int64(weeks[i].Id))
		if workouts == nil {
			log.Default().Println("getProgramWeeks:", err)
			continue
		}
		weeks[i].Workouts = workouts
	}

	return weeks, nil
}

func (db *DB) UpdateProgramWeek(week *ProgramWeek) error {
	var err error = nil
	_, err = db.Exec(context.Background(), `
UPDATE program_weeks
SET notes       = COALESCE($1, notes),
    week_number = COALESCE($2, week_number)
WHERE id=$3`, week.Notes, week.WeekNumber, week.Id)
	if err != nil {
		return err
	}

	for i := range week.Workouts {
		if err := db.UpdateProgramWorkout(&week.Workouts[i]); err != nil {
			fmt.Println("Failed to update program week workouts:", err)
			return err
		}
	}

	return err
}

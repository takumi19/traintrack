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

func (db *DB) GetProgramByID(id int64) (*Program, error) {
	rows, err := db.Query(context.Background(), "SELECT * FROM program_templates WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	program, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Program])
	if err != nil {
		return nil, err
	}

	return &program, nil
}

func (db *DB) ListPrograms() ([]Program, error) {
	rows, err := db.Query(context.Background(), "SELECT * FROM program_templates")
	if err != nil {
		return nil, err
	}

	programs, err := pgx.CollectRows(rows, pgx.RowToStructByName[Program])
	if err != nil {
		return nil, err
	}

	for i := range programs {
		weeks, err := db.ListProgramWeeks(int64(programs[i].Id))
		if weeks == nil {
			log.Default().Println("ListPrograms:", err)
			continue
		}
		programs[i].Weeks = weeks
	}

	return programs, nil
}

func (db *DB) UpdateProgram(program *Program) error {
	tx, err := db.Begin(context.Background())
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
		if err := db.UpdateProgramWeek(&program.Weeks[i]); err != nil {
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


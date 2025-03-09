package database

type ExerciseInfo struct {
	Id           int       `json:"id" db:"id"`
	Name         int       `json:"name" db:"name"`
	Notes        int       `json:"notes" db:"notes"`
	IsRepBased   bool      `json:"is_rep_based" db:"is_rep_based"`
	IsBodyweight bool      `json:"is_bodyweight" db:"is_bodyweight"`
	MuscleGroups []*string `json:"muscle_groups" db:"-"`
}

package database_test

import(
	"testing"
	"traintrack/internal/database"
)

func TestDB_GetFullLogByUserID(t *testing.T) {
	tests := []struct {
		name string
		dbUrl       string
		automigrate bool
		userId  int64
		want    []database.LoggedWorkout
		wantErr bool
	}{
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := database.New(tt.dbUrl, tt.automigrate)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := db.GetFullLogByUserID(tt.userId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetFullLogByUserID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetFullLogByUserID() succeeded unexpectedly")
			}
			if true {
				t.Errorf("GetFullLogByUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}


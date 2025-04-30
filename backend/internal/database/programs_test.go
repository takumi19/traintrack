package database_test

import(
	"testing"
	"traintrack/internal/database"
)

func TestDB_ListPrograms(t *testing.T) {
	tests := []struct {
		name string
		dbUrl       string
		automigrate bool
		want        []database.Program
		wantErr     bool
	}{
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := database.New(tt.dbUrl, tt.automigrate)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := db.ListPrograms()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ListPrograms() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ListPrograms() succeeded unexpectedly")
			}
			if true {
				t.Errorf("ListPrograms() = %v, want %v", got, tt.want)
			}
		})
	}
}


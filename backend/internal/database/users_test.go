package database_test

import(
	"testing"
	"traintrack/internal/database"
)

func TestDB_ReadUser(t *testing.T) {
	tests := []struct {
		name string
		dbUrl       string
		automigrate bool
		id      int64
		want    *database.User
		wantErr bool
	}{
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := database.New(tt.dbUrl, tt.automigrate)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := db.ReadUser(tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ReadUser() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ReadUser() succeeded unexpectedly")
			}
			if true {
				t.Errorf("ReadUser() = %v, want %v", got, tt.want)
			}
		})
	}
}


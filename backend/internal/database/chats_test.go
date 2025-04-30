package database_test

import(
	"testing"
	"traintrack/internal/database"
)

func TestDB_GetChatsByUserId(t *testing.T) {
	tests := []struct {
		name string
		dbUrl       string
		automigrate bool
		id      int64
		want    []database.Chat
		wantErr bool
	}{
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := database.New(tt.dbUrl, tt.automigrate)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := db.GetChatsByUserId(tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetChatsByUserId() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetChatsByUserId() succeeded unexpectedly")
			}
			if true {
				t.Errorf("GetChatsByUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}


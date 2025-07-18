package store

import (
	"database/sql"
	"testing"
)

func setupTestD(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost port=5433 user=workout_user password=workout_password dbname=workout_db sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to open test database: %w", err)
	}

	err = Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("Failed to migrate test database: %w", err)
	}

	_, err = db.Exec(`TRUNCATE workouts, exercises CASCADE`)
	if err != nil {
		t.Fatalf("Failed to truncate test database: %w", err)
	}
	
	return db
}

package store

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestD(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost port=5433 user=workout_user password=workout_password dbname=workout_db sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	_, err = db.Exec(`TRUNCATE workouts, exercises CASCADE`)
	if err != nil {
		t.Fatalf("Failed to truncate test database: %v", err)
	}

	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestD(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name          string
		workout       *Workout
		expectedError bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "Push Day",
				Description:     "Chest, triceps and shoulders workout",
				DurationMinutes: 60,
				ID:              1,
				CaloriesBurned:  400,
				Exercises: []Exercise{
					{
						Name:       "Bench Press",
						Reps:       intToPtr(8),
						Sets:       3,
						Weight:     floatToPtr(82.5),
						Notes:      "Focus on form",
						OrderIndex: 1,
					},
				},
			},
			expectedError: false,
		},
		{
			name: "invalid exercises workout",
			workout: &Workout{
				Title:           "Push Day",
				Description:     "Chest, triceps and shoulders workout",
				DurationMinutes: 60,
				CaloriesBurned:  400,
				Exercises: []Exercise{
					{
						Name:       "Chest Fly",
						Sets:       3,
						Reps:       intToPtr(10),
						Notes:      "Use light weights",
						OrderIndex: 1,
					},
					//Can't have set,reps AND duration_seconds set at the same time!
					{
						Name:            "Tricep Dips",
						Reps:            intToPtr(12),
						Sets:            3,
						DurationSeconds: intToPtr(20),
						Weight:          floatToPtr(180.5),
						Notes:           "Bodyweight exercise",
						OrderIndex:      2,
					},
				},
			},
			expectedError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(test.workout)
			if test.expectedError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.workout.Title, createdWorkout.Title)
			assert.Equal(t, test.workout.Description, createdWorkout.Description)
			assert.Equal(t, test.workout.DurationMinutes, createdWorkout.DurationMinutes)
			assert.Equal(t, test.workout.CaloriesBurned, createdWorkout.CaloriesBurned)
			assert.Equal(t, len(test.workout.Exercises), len(createdWorkout.Exercises))

			gotWorkout, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)
			assert.Equal(t, createdWorkout.ID, gotWorkout.ID)
			assert.Equal(t, len(test.workout.Exercises), len(gotWorkout.Exercises))

			for i, exercise := range test.workout.Exercises {
				assert.Equal(t, exercise.Name, gotWorkout.Exercises[i].Name)
				assert.Equal(t, exercise.Sets, gotWorkout.Exercises[i].Sets)
				assert.Equal(t, exercise.Reps, gotWorkout.Exercises[i].Reps)
				assert.Equal(t, exercise.DurationSeconds, gotWorkout.Exercises[i].DurationSeconds)
				assert.Equal(t, exercise.Weight, gotWorkout.Exercises[i].Weight)
				assert.Equal(t, exercise.Notes, gotWorkout.Exercises[i].Notes)
				assert.Equal(t, exercise.OrderIndex, gotWorkout.Exercises[i].OrderIndex)
				assert.NotZero(t, gotWorkout.Exercises[i].ID)

			}

		})
	}
}
func intToPtr(i int) *int {
	return &i
}

func floatToPtr(i float64) *float64 {
	return &i
}

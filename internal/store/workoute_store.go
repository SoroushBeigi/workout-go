package store

import (
	"database/sql"
)

type Workout struct {
	ID              int        `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	DurationMinutes int        `json:"duration_minutes"`
	CaloriesBurned  int        `json:"calories_burned"`
	Exercises       []Exercise `json:"exercises"`
}

type Exercise struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}
type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutByID(id int64) (*Workout, error)
	UpdateWorkout(*Workout) error
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

func (pg PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	query := `INSERT INTO workouts (title, description, duration_minutes, calories_burned)
		VALUES ($1, $2, $3, $4) RETURNING id`

	err = tx.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)
	if err != nil {
		return nil, err
	}
	for _, exercise := range workout.Exercises {
		query := `INSERT INTO exercises (workout_id, name, sets, reps, duration_seconds, weight, notes, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
		err = tx.QueryRow(query, workout.ID, exercise.Name, exercise.Sets, exercise.Reps, exercise.DurationSeconds, exercise.Weight, exercise.Notes, exercise.OrderIndex).Scan(&exercise.ID)
		if err != nil {
			return nil, err
		}

	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (pg PostgresWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {
	workout := &Workout{}
	query := `SELECT id, title, description, duration_minutes, calories_burned FROM workouts WHERE id = $1`
	err := pg.db.QueryRow(query, id).Scan(&workout.ID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	exerciseQuery := `SELECT id, name, sets, reps, duration_seconds, weight, notes, order_index FROM exercises WHERE workout_id = $1 ORDER BY order_index`

	rows, err := pg.db.Query(exerciseQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var exercise Exercise
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Sets,
			&exercise.Reps,
			&exercise.DurationSeconds,
			&exercise.Weight,
			&exercise.DurationSeconds,
			&exercise.Weight,
			&exercise.Notes,
			&exercise.OrderIndex,
		)
		if err != nil {
			return nil, err
		}

		workout.Exercises = append(workout.Exercises, exercise)
	}
	return workout, nil
}

func (pg *PostgresWorkoutStore) UpdateWorkout(workout *Workout) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `
	UPDATE workouts
	SET title = $1, description = $2, duration_minutes = $3, calories_burned = $4 
	WHERE id = $5
	`
	result, err := tx.Exec(
		query,
		workout.Title,
		workout.Description,
		workout.DurationMinutes,
		workout.CaloriesBurned,
		workout.ID,
	)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	_, err = tx.Exec(`DELETE FROM exercises WHERE workout_id = $1`, workout.ID)
	if err != nil {
		return err
	}

	for _, ex := range workout.Exercises {
		query := `
		INSERT INTO exercises (workout_id, name, sets, reps, duration_seconds, weight, notes, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
		_, err := tx.Exec(
			query,
			workout.ID,
			ex.Name,
			ex.Sets,
			ex.Reps,
			ex.DurationSeconds,
			ex.Weight,
			ex.Notes,
			ex.OrderIndex,
		)
		if err != nil {
			return err
		}

	}

	return tx.Commit()
}

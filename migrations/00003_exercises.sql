-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS exercises (
    id BIGSERIAL PRIMARY KEY,
    --user_id
    workout_id BIGINT NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    sets INTEGER NOT NULL,
    reps INTEGER,
    duration_seconds INTEGER,
    weight DECIMAL(5, 2),
    notes TEXT,
    order_index INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_exercise CHECK (
        (sets IS NOT NULL OR duration_seconds IS NOT NULL) AND 
        (reps is NULL OR duration_seconds IS NULL)
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE exercises;
-- +goose StatementEnd
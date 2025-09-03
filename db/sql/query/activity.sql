-- name: CreateActivity :one
INSERT INTO activities (
    user_id, activity_type, done_at, duration_minutes, calories_burned, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, NOW(), NOW()
)
RETURNING id, user_id, activity_type, done_at, duration_minutes, calories_burned, created_at, updated_at;

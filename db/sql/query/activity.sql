-- name: CreateActivity :one
INSERT INTO activities (
    user_id, activity_type, done_at, duration_minutes, calories_burned, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, NOW(), NOW()
)
RETURNING id, user_id, activity_type, done_at, duration_minutes, calories_burned, created_at, updated_at;


-- name: DeleteActivity :execrows
DELETE FROM activities WHERE id = $1;

-- name: PatchActivity :one
UPDATE activities
SET
    activity_type = COALESCE(sqlc.narg('activity_type'), activity_type),
    done_at = COALESCE(sqlc.narg('done_at'), done_at),
    duration_minutes = COALESCE(sqlc.narg('duration_minutes'), duration_minutes),
    calories_burned = COALESCE(sqlc.narg('calories_burned'), calories_burned),
    updated_at = now()
WHERE id = sqlc.arg('id')
RETURNING id, activity_type, done_at, duration_minutes, calories_burned, created_at, updated_at;

-- name: GetPaginatedActivity :many
SELECT id, activity_type, done_at, duration_minutes, calories_burned, created_at
FROM activities
WHERE user_id = $1
  AND (COALESCE($2::text, '') = '' OR activity_type = $2::text)
  AND (COALESCE($3::timestamptz, '1900-01-01'::timestamptz) = '1900-01-01'::timestamptz OR done_at >= $3::timestamptz)
  AND (COALESCE($4::timestamptz, '2100-01-01'::timestamptz) = '2100-01-01'::timestamptz OR done_at <= $4::timestamptz)
  AND (COALESCE($5::integer, -1) = -1 OR calories_burned >= $5::integer)
  AND (COALESCE($6::integer, -1) = -1 OR calories_burned <= $6::integer)
LIMIT $7 OFFSET $8;
-- name: SelectUserByEmail :one
SELECT id, password_hash FROM users where email = $1;

-- name: CreateUser :exec
INSERT INTO users (email, password_hash) VALUES ($1, $2);

-- name: SelectProfileById :one
SELECT preference, weight_unit as weightUnit, height_unit as heightUnit, weight, height, email, name, image_uri as imageUri FROM users where id = $1;


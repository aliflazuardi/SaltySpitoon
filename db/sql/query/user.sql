-- name: SelectUserByEmail :one
SELECT id, password_hash FROM users where email = $1;

-- name: CreateUser :exec
INSERT INTO users (email, password_hash) VALUES ($1, $2);

-- name: SelectProfileById :one
SELECT preference, weight_unit as weightUnit, height_unit as heightUnit, weight, height, email, name, image_uri as imageUri FROM users where id = $1;

-- name: PatchProfileById :exec
UPDATE users 
    SET 
        preference = COALESCE(sqlc.narg(preference), preference),
        weight_unit = COALESCE(sqlc.narg(weight_unit), weight_unit),
        height_unit = COALESCE(sqlc.narg(height_unit), height_unit),
        weight = COALESCE(sqlc.narg(weight), weight),
        height = COALESCE(sqlc.narg(height), height),
        name = COALESCE(sqlc.narg(name), name),
        image_uri = COALESCE(sqlc.narg(image_uri), image_uri)
WHERE id = $1;
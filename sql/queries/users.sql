-- name: CreateUser :exec
INSERT INTO users (id, created_at, updated_at, username, hashed_password, steam_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);
--

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;
--

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: DeleteUsers :exec
DELETE FROM users;
--
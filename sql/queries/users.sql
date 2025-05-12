-- name: CreateUser :exec
INSERT INTO users (id, created_at, updated_at, email, hashed_password, steam_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);
--

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
--

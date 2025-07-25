// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO users (id, created_at, updated_at, username, hashed_password, steam_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
`

type CreateUserParams struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Username       string
	HashedPassword string
	SteamID        string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Username,
		arg.HashedPassword,
		arg.SteamID,
	)
	return err
}

const deleteUsers = `-- name: DeleteUsers :exec

DELETE FROM users
`

func (q *Queries) DeleteUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteUsers)
	return err
}

const getUserByID = `-- name: GetUserByID :one

SELECT id, created_at, updated_at, username, hashed_password, steam_id FROM users WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.HashedPassword,
		&i.SteamID,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one

SELECT id, created_at, updated_at, username, hashed_password, steam_id FROM users WHERE username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.HashedPassword,
		&i.SteamID,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :exec

UPDATE users
SET
    username = COALESCE(NULLIF($1::text, ''), username),
    hashed_password = COALESCE(NULLIF($2::text, ''), hashed_password),
    steam_id = COALESCE(NULLIF($3::text, ''), steam_id),
    updated_at = $4
WHERE id = $5
`

type UpdateUserParams struct {
	Column1   string
	Column2   string
	Column3   string
	UpdatedAt time.Time
	ID        uuid.UUID
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.Column1,
		arg.Column2,
		arg.Column3,
		arg.UpdatedAt,
		arg.ID,
	)
	return err
}

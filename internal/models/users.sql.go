// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package models

import (
	"context"
)

const usersFindByEmail = `-- name: UsersFindByEmail :one
SELECT id, uid, first_name, last_name, email, password, avatar, role, created_at, updated_at, deleted_at, last_active_at FROM users WHERE email = $1
`

func (q *Queries) UsersFindByEmail(ctx context.Context, email string) (*User, error) {
	row := q.db.QueryRow(ctx, usersFindByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.Avatar,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.LastActiveAt,
	)
	return &i, err
}

const usersFindByUid = `-- name: UsersFindByUid :one
SELECT id, uid, first_name, last_name, email, password, avatar, role, created_at, updated_at, deleted_at, last_active_at FROM users WHERE uid = $1
`

func (q *Queries) UsersFindByUid(ctx context.Context, uid string) (*User, error) {
	row := q.db.QueryRow(ctx, usersFindByUid, uid)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.Avatar,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.LastActiveAt,
	)
	return &i, err
}

const usersInsert = `-- name: UsersInsert :one
INSERT INTO users (
  uid, first_name, last_name, email, password, avatar, role
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING id, uid, first_name, last_name, email, password, avatar, role, created_at, updated_at, deleted_at, last_active_at
`

type UsersInsertParams struct {
	Uid       string  `json:"uid"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     string  `json:"email"`
	Password  *string `json:"password"`
	Avatar    *string `json:"avatar"`
	Role      string  `json:"role"`
}

func (q *Queries) UsersInsert(ctx context.Context, arg UsersInsertParams) (*User, error) {
	row := q.db.QueryRow(ctx, usersInsert,
		arg.Uid,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Password,
		arg.Avatar,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.Avatar,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.LastActiveAt,
	)
	return &i, err
}

const usersUpdateLastActive = `-- name: UsersUpdateLastActive :exec
UPDATE users SET
  last_active_at = now(),
  updated_at = now()
WHERE id = $1
`

func (q *Queries) UsersUpdateLastActive(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, usersUpdateLastActive, id)
	return err
}

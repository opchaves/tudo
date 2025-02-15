// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: notes.sql

package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const notesDelete = `-- name: NotesDelete :exec
DELETE FROM notes WHERE id=$1 AND user_id=$2
`

type NotesDeleteParams struct {
	ID     int32 `json:"id"`
	UserID int32 `json:"user_id"`
}

func (q *Queries) NotesDelete(ctx context.Context, arg NotesDeleteParams) error {
	_, err := q.db.Exec(ctx, notesDelete, arg.ID, arg.UserID)
	return err
}

const notesGetByID = `-- name: NotesGetByID :one
SELECT id, user_id, title, content, created_at FROM notes WHERE id=$1 AND user_id=$2
`

type NotesGetByIDParams struct {
	ID     int32 `json:"id"`
	UserID int32 `json:"user_id"`
}

type NotesGetByIDRow struct {
	ID        int32            `json:"id"`
	UserID    int32            `json:"user_id"`
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

func (q *Queries) NotesGetByID(ctx context.Context, arg NotesGetByIDParams) (*NotesGetByIDRow, error) {
	row := q.db.QueryRow(ctx, notesGetByID, arg.ID, arg.UserID)
	var i NotesGetByIDRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Content,
		&i.CreatedAt,
	)
	return &i, err
}

const notesGetByUser = `-- name: NotesGetByUser :many
SELECT id, user_id, title, content, created_at FROM notes WHERE user_id=$1
`

type NotesGetByUserRow struct {
	ID        int32            `json:"id"`
	UserID    int32            `json:"user_id"`
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

func (q *Queries) NotesGetByUser(ctx context.Context, userID int32) ([]*NotesGetByUserRow, error) {
	rows, err := q.db.Query(ctx, notesGetByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*NotesGetByUserRow
	for rows.Next() {
		var i NotesGetByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Content,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const notesInsert = `-- name: NotesInsert :one
INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3) RETURNING id, uid, title, content, user_id, created_at, updated_at
`

type NotesInsertParams struct {
	UserID  int32  `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (q *Queries) NotesInsert(ctx context.Context, arg NotesInsertParams) (*Note, error) {
	row := q.db.QueryRow(ctx, notesInsert, arg.UserID, arg.Title, arg.Content)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.Title,
		&i.Content,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const notesUpdate = `-- name: NotesUpdate :one
UPDATE notes SET title=$1, content=$2 WHERE id=$3 AND user_id=$4 RETURNING id, uid, title, content, user_id, created_at, updated_at
`

type NotesUpdateParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	ID      int32  `json:"id"`
	UserID  int32  `json:"user_id"`
}

func (q *Queries) NotesUpdate(ctx context.Context, arg NotesUpdateParams) (*Note, error) {
	row := q.db.QueryRow(ctx, notesUpdate,
		arg.Title,
		arg.Content,
		arg.ID,
		arg.UserID,
	)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.Uid,
		&i.Title,
		&i.Content,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

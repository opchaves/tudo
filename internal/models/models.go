// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Note struct {
	ID        int32
	UserID    int32
	Title     string
	Content   string
	CreatedAt pgtype.Timestamp
}

type User struct {
	ID       int32
	Name     string
	Email    string
	Password string
}

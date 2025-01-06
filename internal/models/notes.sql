-- name: NotesGetByUser :many
SELECT id, user_id, title, content, created_at FROM notes WHERE user_id=$1;

-- name: NotesInsert :one
INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3) RETURNING *;

-- name: NotesGetByID :one
SELECT id, user_id, title, content, created_at FROM notes WHERE id=$1 AND user_id=$2;

-- name: NotesUpdate :one
UPDATE notes SET title=$1, content=$2 WHERE id=$3 AND user_id=$4 RETURNING *;

-- name: NotesDelete :exec
DELETE FROM notes WHERE id=$1 AND user_id=$2;

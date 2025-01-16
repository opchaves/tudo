-- name: UsersInsert :one
INSERT INTO users (
  uid, first_name, last_name, email, password, avatar, role
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: UsersUpdateLastActive :exec
UPDATE users SET
  last_active_at = now(),
  updated_at = now()
WHERE id = $1;

-- name: UsersFindByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UsersFindByUid :one
SELECT * FROM users WHERE uid = $1;


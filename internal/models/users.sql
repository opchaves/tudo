-- name: UsersInsert :one
INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING *;

-- name: UsersFindByEmail :one
SELECT id, name, email, password FROM users WHERE email = $1;

-- name: GetUser :one
SELECT * FROM users WHERE unid = ? LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (unid) VALUES (?);

-- name: DeleteUser :exec
DELETE FROM users WHERE unid = ?;

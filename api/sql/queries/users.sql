-- name: DoesUserExistByUnid :one
SELECT EXISTS(SELECT 1 FROM users WHERE unid = ?) AS user_exists;

-- name: DoesUserExistByDiscordID :one
SELECT EXISTS(SELECT 1 FROM discord WHERE discord_id = ?) AS user_exists;

-- name: DoesUserExistByMinecraftID :one
SELECT EXISTS(SELECT 1 FROM minecraft WHERE minecraft_id = ?) AS user_exists;

-- name: NewUser :exec
INSERT INTO users (unid) VALUES (?);

-- name: NewDiscordUser :exec
INSERT INTO discord (unid, discord_id) VALUES (?, ?);

-- name: NewMinecraftUser :exec
INSERT INTO minecraft (unid, minecraft_id) VALUES (?, ?);

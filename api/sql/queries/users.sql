-- name: UserExistsByUnid :one
SELECT EXISTS(SELECT 1 FROM users WHERE unid = ?) AS user_exists;

-- name: UserExistsByDiscordID :one
SELECT EXISTS(SELECT 1 FROM discord WHERE discord_id = ?) AS user_exists;

-- name: UserExistsByMinecraftID :one
SELECT EXISTS(SELECT 1 FROM minecraft WHERE minecraft_id = ?) AS user_exists;
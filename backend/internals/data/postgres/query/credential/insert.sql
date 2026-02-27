-- name: Credential_Insert :one
INSERT INTO "credential" ("token", "user_id", "app_id", "channel_id", "type", "permissions", "expires_at")
VALUES (@token, @user_id, @app_id, @channel_id, @type, @permissions, @expires_at)
RETURNING *;

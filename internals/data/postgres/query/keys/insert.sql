-- name: App_Key_Insert :one
INSERT INTO "app_api_key"
("appId", "key", "name")
VALUES
(@app_id, @key, @name)
RETURNING *;
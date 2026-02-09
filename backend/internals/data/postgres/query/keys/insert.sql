-- name: App_Key_Insert :one
INSERT INTO "app_api_key"
("app_id", "id", "name", "prefix", "suffix")
VALUES
(@app_id, @id, @name, @prefix, @suffix)
RETURNING *;
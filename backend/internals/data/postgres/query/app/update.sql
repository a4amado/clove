-- name: App_Update :one
UPDATE "app"
SET
    "app_slug" = @app_slug,
    "allowed_origins" = @allowed_origins,
    "app_type" = @app_type
WHERE "id" = @id
RETURNING *;

-- name: App_Insert :one
INSERT INTO "app"
("app_slug", "region", "app_type", "user_id", "allowed_origins")
values
(@app_slug, @regions, @app_type, @user_id, @allowed_origins)
RETURNING *;
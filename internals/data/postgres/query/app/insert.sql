-- name: App_Insert :one
INSERT INTO "app"
("appSlug", "region", "appType", "userId", "allowedOrigins")
values
(@app_slug, @regions, @app_type, @user_id, @allowed_origins)
RETURNING *;
-- name: App_Key_Select :one
SELECT *
FROM "app_api_key"
WHERE "id" = @key_id AND "app_id" = @app_id
limit 1;
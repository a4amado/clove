-- name: App_Key_Select :one
SELECT *
FROM "app_api_key"
WHERE "id" = @key_id AND "appId" = @app_id
limit 1;
-- name: App_Key_Delete :execrows
Delete from "app_api_key"
WHERE "id" = @id and "app_id" = @app_id;

-- name: App_Key_List :many
SELECT * from "app_api_key"
WHERE   "app_id" = @app_id
Limit 20
OFFSET 10 *  sqlc.arg(page_idx)::int;
 
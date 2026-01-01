-- name: ListAppApiKeys :many
SELECT * from "app_api_key"
WHERE   "appId" = @app_id
Limit 20
OFFSET 10 *  sqlc.arg(page_idx)::int;

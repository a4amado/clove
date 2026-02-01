-- name: App_Select :one
SELECT * from "app"
where "id" = @app_id
limit 1;

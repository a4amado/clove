-- name: FindAppById :one
SELECT * from "app"
where "id" = @appId
limit 1;

-- name: ListAppRegions :one
SELECT "region"
FROM "app"
WHERE "id" = @id
LIMIT 1;
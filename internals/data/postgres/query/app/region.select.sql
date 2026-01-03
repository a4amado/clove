-- name: App_Region_Select :one
SELECT "region"
FROM "app"
WHERE "id" = @id
LIMIT 1;
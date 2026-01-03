-- name: App_Region_Update :exec
UPDATE "app"
SET "region" = @region
WHERE "id" = @id
RETURNING "region";
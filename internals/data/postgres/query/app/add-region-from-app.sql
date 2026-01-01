-- name: AddRegionToApp :one
UPDATE "app"
SET "region" = @region
WHERE "id" = @id
RETURNING "region";
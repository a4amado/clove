-- name: AddRegionToApp :one
UPDATE "app"
SET "region" = ARRAY_APPEND("region", @region::region)
WHERE "id" = @id
RETURNING "region";
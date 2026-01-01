-- name: RemoveRegionFromApp :one
UPDATE "app"
SET "region" = ARRAY_REMOVE("region", @region::region)
WHERE "id" = @id
RETURNING "region";
-- name: UpdateAppRegions :exec
UPDATE "app"
SET "region" = @regions
WHERE "id" = @app_id
RETURNING *;
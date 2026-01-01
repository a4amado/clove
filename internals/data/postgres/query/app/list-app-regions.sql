<<<<<<< HEAD
-- name: RemoveRegionFromApp :one
SELECT "region"
WHERE "id" = @id
from "app"
=======
-- name: ListAppRegions :one
SELECT "region"
FROM "app"
WHERE "id" = @id

>>>>>>> fresh-start
LIMIT 1;
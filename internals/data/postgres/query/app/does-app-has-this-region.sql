-- name: DoesAppHaveRegion :one
SELECT EXISTS(
    SELECT 1 
    FROM "app"
    WHERE "id" = @id 
      AND @region = ANY("region")
) AS has_region;
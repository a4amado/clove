-- name: App_Delete :exec
DELETE FROM "app"
where "id" = @id;
-- name: DeleteApp :exec

DELETE FROM "app"
where "id" = @id;
-- name: User_Insert :one
insert into "user"("hash") values(@hash) RETURNING *;
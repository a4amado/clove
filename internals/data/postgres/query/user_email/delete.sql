-- name: UserEmail_Delete :exec
DELETE FROM "user_email"
WHERE "id" = @id;

-- name: User_Select :one
SELECT * FROM "user"
WHERE "id" = @user_id
LIMIT 1;

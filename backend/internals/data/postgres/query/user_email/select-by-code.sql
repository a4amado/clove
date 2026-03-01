-- name: UserEmail_SelectByCode :one
SELECT *
FROM "user_email"
WHERE "code" = @code
LIMIT 1;

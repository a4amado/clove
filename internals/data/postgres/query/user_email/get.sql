-- name: UserEmail_Select :one
SELECT
    *
FROM "user_email"
WHERE "id" = @user_id;


-- name: UserEmail_SelectByEmail :one
SELECT
    *
FROM "user_email"
WHERE "email" = @email;

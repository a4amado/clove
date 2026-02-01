-- name: UserEmails_List_ByUserID :many
SELECT
    *
FROM "user_email"
WHERE "user_id" = @user_id
ORDER BY created_at DESC;

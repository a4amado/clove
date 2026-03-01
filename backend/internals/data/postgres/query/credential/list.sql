-- name: Credential_List_ByUserID :many
SELECT *
FROM "credential"
WHERE "user_id" = @user_id
ORDER BY "created_at" DESC;

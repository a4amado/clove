-- name: App_List_ByUserID :many
SELECT * FROM "app"
WHERE "user_id" = @user_id
ORDER BY "app_slug" ASC;

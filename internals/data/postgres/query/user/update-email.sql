 
-- name: UpdateUserEmail :exec
UPDATE "user"
SET "email" = @email
Where "id" = @user_Id;
 
-- name: User_Password_Update :exec
UPDATE "user"
SET "hash" = @hash
Where "id" =  @user_id;
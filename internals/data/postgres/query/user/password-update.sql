 
-- name: UpdateUserPassword :exec
UPDATE "user"
SET "hash" = @hash
Where "id" =  @user_id;
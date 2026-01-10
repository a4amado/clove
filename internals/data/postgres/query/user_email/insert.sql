-- name: UserEmail_Insert :one
INSERT INTO "user_email"
("email", "user_id", "code")
VALUES(@email, @user_id, @code)
RETURNING *;
-- name: UserEmails_Verify :exec
UPDATE "user_email"
SET
    "verified_at" = NOW()
WHERE "id" = @id;

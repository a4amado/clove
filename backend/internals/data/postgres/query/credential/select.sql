-- name: Credential_SelectByToken :one
SELECT * FROM "credential"
WHERE "token" = @token
LIMIT 1;

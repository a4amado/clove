-- name: Credential_Delete :exec
DELETE FROM "credential"
WHERE "id" = @id;

-- name: Credential_DeleteExpired :exec
DELETE FROM "credential"
WHERE "expires_at" < NOW();

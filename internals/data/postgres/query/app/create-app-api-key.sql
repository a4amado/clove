-- CREATE TABLE "app_api_key" (
--     "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     "appId" UUID,
--     "createdAt" TIMESTAMP DEFAULT NOW(),
--     "updatedAt" TIMESTAMP DEFAULT NOW(),
--     CONSTRAINT "api_key_app_fk" FOREIGN KEY ("appId") REFERENCES "app"("id")
-- );

-- name: CreateAppApiKey :one
INSERT INTO "app_api_key"
("appId", "key", "name")
VALUES
(@app_id, @key, @name)
RETURNING *;
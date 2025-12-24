-- "id" UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
-- "app_slug" VARCHAR(255) NOT NULL UNIQUE,
-- "region" region[] NOT NULL,
-- "app_type" app_type NOT NULL 

-- name: InsertApp :one
INSERT INTO "app"
("appSlug", "region", "appType", "userId", "allowedOrigins")
values
(@appSlug, @region, @appType, @userId, @allowedOrigins)
RETURNING *;
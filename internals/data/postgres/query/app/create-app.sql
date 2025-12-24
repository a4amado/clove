 
-- name: InsertApp :one
INSERT INTO "app"
("appSlug", "region", "appType", "userId", "allowedOrigins")
values
(@appSlug, @region, @appType, @userId, @allowedOrigins)
RETURNING *;
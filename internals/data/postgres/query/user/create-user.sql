-- name: InsertUser :one
insert into "user"
("email", "hash")
values(
    @email, @hash
)
RETURNING *;
-- name: CreateUser :one
insert into users (
  username, email, password_hash, password_salt
)
select $1, $2, $3, $4
where not exists (
  select 1 from users
  where username = $1 or email = $2
)
returning *;
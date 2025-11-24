-- name: CreateUser :one
insert into users (
  username, email, password_hash
)
select $1, $2, $3
where not exists (
  select 1 from users
  where username = $1 or email = $2
)
returning *;
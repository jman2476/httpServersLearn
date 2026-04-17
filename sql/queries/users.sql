-- name: CreateUser :one
insert into users (id, created_at, updated_at, email, hashed_password)
values (
    gen_random_uuid(),
    now(),
    now(),
    $1, $2
) returning id, created_at, updated_at, email;

-- name: ClearUsers :exec
delete from users;
-- name: CreateRefreshToken :one
insert into refresh_tokens(token, created_at, updated_at, user_id, expires_at)
values (
    $1,
    now(),
    now(),
    $2,
    date_add(now(), '60 days')
) returning token, created_at, expires_at;


-- name: GetRefreshToken :one
select * from refresh_tokens
where token = $1;

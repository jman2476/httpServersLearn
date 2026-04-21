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

-- name: GetUserfromRefreshToken :one
select users.* from users
join refresh_tokens on users.id = refresh_tokens.user_id
where refresh_tokens.token = $1
and revoked_at is null
and expires_at > now();

-- name: RevokeToken :exec
update refresh_tokens
set updated_at = now(), revoked_at = now()
where token = $1;

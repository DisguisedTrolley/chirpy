-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
	$1,
	NOW(),
	NOW(),
	$2,
	$3,
	NULL
);


-- name: GetRefreshToken :one
select *
from refresh_tokens
where token = $1
;


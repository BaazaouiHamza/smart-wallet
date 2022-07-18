-- name: GetUserPolicy :one
SELECT * FROM user_policies WHERE id = $1 LIMIT 1;

-- name: DeleteUserPolicy :exec
DELETE FROM user_policies WHERE id = $1;

-- name: ListUserPolicies :many
SELECT * FROM user_policies WHERE nym_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

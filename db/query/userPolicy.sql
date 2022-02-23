-- name: GetUserPolicy :one
SELECT * FROM ONLY user_policy up JOIN transaction_trigger_policy ttp ON up.id=ttp.id
WHERE up.id = $1 LIMIT 1;

-- name: DeleteUserPolicy :exec
DELETE FROM user_policy WHERE id = $1;

-- name: ListUserPolicies :many
SELECT * FROM user_policy WHERE nym_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;
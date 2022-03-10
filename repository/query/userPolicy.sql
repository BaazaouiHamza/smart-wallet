-- name: GetUserPolicy :one
SELECT * FROM ONLY user_policies up JOIN transaction_trigger_policies ttp ON up.id=ttp.id
WHERE up.id = $1 LIMIT 1;

-- name: DeleteUserPolicy :exec
DELETE FROM user_policies WHERE id = $1;

-- name: ListUserPolicies :many
SELECT * FROM user_policies WHERE nym_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

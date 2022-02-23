-- name: CreateTransactionTriggerPolicy :one
INSERT INTO transaction_trigger_policy (
  name,
  description,
  nym_id,
  targeted_balance,
  amount
) VALUES (
  $1,$2,$3,$4,$5
) RETURNING *;

-- name: UpdateTransactionTriggerPolicy :one
UPDATE transaction_trigger_policy 
SET name = $2,
description=$3,
nym_id=$4,
targeted_balance=$5,
amount=$6
WHERE id = $1
RETURNING *;

-- name: GetTransactionTriggerPolicy :one
SELECT * FROM transaction_trigger_policy
WHERE id = $1 LIMIT 1;

-- name: DeleteTransactionTriggerPolicy :exec
DELETE FROM transaction_trigger_policy WHERE id = $1;

-- name: ListTransactionTriggerPolicies :many
SELECT * FROM transaction_trigger_policy WHERE nym_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;
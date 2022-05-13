-- name: CreateTTP :one
INSERT INTO transaction_trigger_policies (
  name,
  description,
  nym_id,
  targeted_balance,
  amount,
  recipient
) VALUES (
  $1,$2,$3,$4,$5,$6
) RETURNING *;

-- name: UpdateTTP :one
UPDATE transaction_trigger_policies 
SET name = $2,
description=$3,
nym_id=$4,
targeted_balance=$5,
amount=$6,
recipient=$7
WHERE id = $1
RETURNING *;

-- name: GetTTP :one
SELECT * FROM transaction_trigger_policies
WHERE id = $1 LIMIT 1;

-- name: DeleteTTP :exec
DELETE FROM transaction_trigger_policies WHERE id = $1;

-- name: ListTTP :many
SELECT *, count(*) OVER() AS full_count FROM transaction_trigger_policies WHERE nym_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

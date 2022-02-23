-- name: CreateRoutineTransactionPolicy :one
INSERT INTO routine_transaction_policy (
  name,
  description,
  nym_id,
  schedule_start_date, 
  schedule_end_date,
  frequency,
  amount
) VALUES (
  $1,$2,$3,$4,$5,$6,$7
) RETURNING *;


-- name: UpdateRoutineTransactionPolicy :one
UPDATE routine_transaction_policy 
SET name = $2,
description=$3,
nym_id=$4,
schedule_start_date=$5,
schedule_end_date=$6,
frequency=$7,
amount=$8
WHERE id = $1
RETURNING *;

-- name: GetRoutineTransactionPolicy :one
SELECT * FROM routine_transaction_policy
WHERE id = $1 LIMIT 1;

-- name: DeleteRoutineTransactionPolicy :exec
DELETE FROM routine_transaction_policy WHERE id = $1;

-- name: ListRoutineTransactionPolicies :many
SELECT * FROM routine_transaction_policy WHERE nym_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;
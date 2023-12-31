-- name: CreateRTP :one
INSERT INTO routine_transaction_policies (name, description, nym_id, schedule_start_date,
  schedule_end_date, frequency, amount, recipient)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
  *;

-- name: UpdateRTP :one
UPDATE
  routine_transaction_policies
SET
  name = $2,
  description = $3,
  nym_id = $4,
  schedule_start_date = $5,
  schedule_end_date = $6,
  frequency = $7,
  amount = $8,
  recipient = $9
WHERE
  id = $1
RETURNING
  *;

-- name: GetRTP :one
SELECT
  *
FROM
  routine_transaction_policies
WHERE
  nym_id = $1
  AND id = $2;

-- name: DeleteRTP :exec
DELETE FROM routine_transaction_policies
WHERE id = $1;

-- name: ListRTP :many
SELECT
  *,
  COUNT(*) OVER () AS full_count
FROM
  routine_transaction_policies
WHERE
  nym_id = $1
ORDER BY
  id OFFSET $3
LIMIT $2;

-- name: GetRTPs :many
SELECT
  *
FROM
  routine_transaction_policies
WHERE
  nym_id = $1
ORDER BY
  id;

  -- name: GetALlRoutinePolicies :many
  SELECT * from routine_transaction_policies ORDER BY id;



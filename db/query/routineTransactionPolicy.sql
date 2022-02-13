-- name: CreateRoutineTransactionPolicy :one
INSERT INTO routine_transaction_policy (
  name,
  description,
  sender,
	receiver,
	created_at,
  schedule_start_date, 
  schedule_end_date,
  frequency,
  amount
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9
) RETURNING *;
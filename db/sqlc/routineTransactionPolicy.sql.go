// Code generated by sqlc. DO NOT EDIT.
// source: routineTransactionPolicy.sql

package db

import (
	"context"
	"encoding/json"
	"time"
)

const createRTP = `-- name: CreateRTP :one
INSERT INTO routine_transaction_policy (
  name,
  description,
  nym_id,
  schedule_start_date, 
  schedule_end_date,
  frequency,
  amount,
  recipient
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8
) RETURNING id, name, description, nym_id, recipient, created_at, schedule_start_date, schedule_end_date, frequency, amount
`

type CreateRTPParams struct {
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	NymID             string          `json:"nym_id"`
	ScheduleStartDate time.Time       `json:"schedule_start_date"`
	ScheduleEndDate   time.Time       `json:"schedule_end_date"`
	Frequency         string          `json:"frequency"`
	Amount            json.RawMessage `json:"amount"`
	Recipient         string          `json:"recipient"`
}

func (q *Queries) CreateRTP(ctx context.Context, arg CreateRTPParams) (RoutineTransactionPolicy, error) {
	row := q.db.QueryRowContext(ctx, createRTP,
		arg.Name,
		arg.Description,
		arg.NymID,
		arg.ScheduleStartDate,
		arg.ScheduleEndDate,
		arg.Frequency,
		arg.Amount,
		arg.Recipient,
	)
	var i RoutineTransactionPolicy
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.NymID,
		&i.Recipient,
		&i.CreatedAt,
		&i.ScheduleStartDate,
		&i.ScheduleEndDate,
		&i.Frequency,
		&i.Amount,
	)
	return i, err
}

const deleteRTP = `-- name: DeleteRTP :exec
DELETE FROM routine_transaction_policy WHERE id = $1
`

func (q *Queries) DeleteRTP(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteRTP, id)
	return err
}

const getRTP = `-- name: GetRTP :one
SELECT id, name, description, nym_id, recipient, created_at, schedule_start_date, schedule_end_date, frequency, amount FROM routine_transaction_policy
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetRTP(ctx context.Context, id int64) (RoutineTransactionPolicy, error) {
	row := q.db.QueryRowContext(ctx, getRTP, id)
	var i RoutineTransactionPolicy
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.NymID,
		&i.Recipient,
		&i.CreatedAt,
		&i.ScheduleStartDate,
		&i.ScheduleEndDate,
		&i.Frequency,
		&i.Amount,
	)
	return i, err
}

const listRTP = `-- name: ListRTP :many
SELECT id, name, description, nym_id, recipient, created_at, schedule_start_date, schedule_end_date, frequency, amount FROM routine_transaction_policy WHERE nym_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListRTPParams struct {
	NymID  string `json:"nym_id"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) ListRTP(ctx context.Context, arg ListRTPParams) ([]RoutineTransactionPolicy, error) {
	rows, err := q.db.QueryContext(ctx, listRTP, arg.NymID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []RoutineTransactionPolicy{}
	for rows.Next() {
		var i RoutineTransactionPolicy
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.NymID,
			&i.Recipient,
			&i.CreatedAt,
			&i.ScheduleStartDate,
			&i.ScheduleEndDate,
			&i.Frequency,
			&i.Amount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateRTP = `-- name: UpdateRTP :one
UPDATE routine_transaction_policy 
SET name = $2,
description=$3,
nym_id=$4,
schedule_start_date=$5,
schedule_end_date=$6,
frequency=$7,
amount=$8,
recipient=$9
WHERE id = $1
RETURNING id, name, description, nym_id, recipient, created_at, schedule_start_date, schedule_end_date, frequency, amount
`

type UpdateRTPParams struct {
	ID                int64           `json:"id"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	NymID             string          `json:"nym_id"`
	ScheduleStartDate time.Time       `json:"schedule_start_date"`
	ScheduleEndDate   time.Time       `json:"schedule_end_date"`
	Frequency         string          `json:"frequency"`
	Amount            json.RawMessage `json:"amount"`
	Recipient         string          `json:"recipient"`
}

func (q *Queries) UpdateRTP(ctx context.Context, arg UpdateRTPParams) (RoutineTransactionPolicy, error) {
	row := q.db.QueryRowContext(ctx, updateRTP,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.NymID,
		arg.ScheduleStartDate,
		arg.ScheduleEndDate,
		arg.Frequency,
		arg.Amount,
		arg.Recipient,
	)
	var i RoutineTransactionPolicy
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.NymID,
		&i.Recipient,
		&i.CreatedAt,
		&i.ScheduleStartDate,
		&i.ScheduleEndDate,
		&i.Frequency,
		&i.Amount,
	)
	return i, err
}

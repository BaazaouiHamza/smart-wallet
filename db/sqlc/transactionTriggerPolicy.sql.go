// Code generated by sqlc. DO NOT EDIT.
// source: transactionTriggerPolicy.sql

package db

import (
	"context"
	"encoding/json"
)

const createTTP = `-- name: CreateTTP :one
INSERT INTO transaction_trigger_policy (
  name,
  description,
  nym_id,
  targeted_balance,
  amount,
  recipient
) VALUES (
  $1,$2,$3,$4,$5,$6
) RETURNING id, name, description, nym_id, recipient, created_at, targeted_balance, amount
`

type CreateTTPParams struct {
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	NymID           string          `json:"nym_id"`
	TargetedBalance json.RawMessage `json:"targeted_balance"`
	Amount          json.RawMessage `json:"amount"`
	Recipient       string          `json:"recipient"`
}

func (q *Queries) CreateTTP(ctx context.Context, arg CreateTTPParams) (TransactionTriggerPolicy, error) {
	row := q.db.QueryRowContext(ctx, createTTP,
		arg.Name,
		arg.Description,
		arg.NymID,
		arg.TargetedBalance,
		arg.Amount,
		arg.Recipient,
	)
	var i TransactionTriggerPolicy
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.NymID,
		&i.Recipient,
		&i.CreatedAt,
		&i.TargetedBalance,
		&i.Amount,
	)
	return i, err
}

const deleteTTP = `-- name: DeleteTTP :exec
DELETE FROM transaction_trigger_policy WHERE id = $1
`

func (q *Queries) DeleteTTP(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTTP, id)
	return err
}

const getTTP = `-- name: GetTTP :one
SELECT id, name, description, nym_id, recipient, created_at, targeted_balance, amount FROM transaction_trigger_policy
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTTP(ctx context.Context, id int64) (TransactionTriggerPolicy, error) {
	row := q.db.QueryRowContext(ctx, getTTP, id)
	var i TransactionTriggerPolicy
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.NymID,
		&i.Recipient,
		&i.CreatedAt,
		&i.TargetedBalance,
		&i.Amount,
	)
	return i, err
}

const listTTP = `-- name: ListTTP :many
SELECT id, name, description, nym_id, recipient, created_at, targeted_balance, amount FROM transaction_trigger_policy WHERE nym_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListTTPParams struct {
	NymID  string `json:"nym_id"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) ListTTP(ctx context.Context, arg ListTTPParams) ([]TransactionTriggerPolicy, error) {
	rows, err := q.db.QueryContext(ctx, listTTP, arg.NymID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []TransactionTriggerPolicy{}
	for rows.Next() {
		var i TransactionTriggerPolicy
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.NymID,
			&i.Recipient,
			&i.CreatedAt,
			&i.TargetedBalance,
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

const updateTTP = `-- name: UpdateTTP :one
UPDATE transaction_trigger_policy 
SET name = $2,
description=$3,
nym_id=$4,
targeted_balance=$5,
amount=$6,
recipient=$7
WHERE id = $1
RETURNING id, name, description, nym_id, recipient, created_at, targeted_balance, amount
`

type UpdateTTPParams struct {
	ID              int64           `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	NymID           string          `json:"nym_id"`
	TargetedBalance json.RawMessage `json:"targeted_balance"`
	Amount          json.RawMessage `json:"amount"`
	Recipient       string          `json:"recipient"`
}

func (q *Queries) UpdateTTP(ctx context.Context, arg UpdateTTPParams) (TransactionTriggerPolicy, error) {
	row := q.db.QueryRowContext(ctx, updateTTP,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.NymID,
		arg.TargetedBalance,
		arg.Amount,
		arg.Recipient,
	)
	var i TransactionTriggerPolicy
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.NymID,
		&i.Recipient,
		&i.CreatedAt,
		&i.TargetedBalance,
		&i.Amount,
	)
	return i, err
}

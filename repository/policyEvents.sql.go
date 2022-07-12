// Code generated by sqlc. DO NOT EDIT.
// source: policyEvents.sql

package repository

import (
	"context"

	"git.digitus.me/prosperus/protocol/identity"
)

const getPolicyEvent = `-- name: GetPolicyEvent :one
select nym_id, transfer_sequence, policy_id from policy_events WHERE nym_id=$1 AND transfer_sequence=$2
`

type GetPolicyEventParams struct {
	NymID            identity.PublicKey `json:"nym_id"`
	TransferSequence int64              `json:"transfer_sequence"`
}

func (q *Queries) GetPolicyEvent(ctx context.Context, arg GetPolicyEventParams) (PolicyEvent, error) {
	row := q.db.QueryRowContext(ctx, getPolicyEvent, arg.NymID, arg.TransferSequence)
	var i PolicyEvent
	err := row.Scan(&i.NymID, &i.TransferSequence, &i.PolicyID)
	return i, err
}

const insertPolicyEvent = `-- name: InsertPolicyEvent :one
INSERT INTO policy_events (nym_id,transfer_sequence,policy_id) VALUES($1,$2,$3) RETURNING nym_id, transfer_sequence, policy_id
`

type InsertPolicyEventParams struct {
	NymID            identity.PublicKey `json:"nym_id"`
	TransferSequence int64              `json:"transfer_sequence"`
	PolicyID         int64              `json:"policy_id"`
}

func (q *Queries) InsertPolicyEvent(ctx context.Context, arg InsertPolicyEventParams) (PolicyEvent, error) {
	row := q.db.QueryRowContext(ctx, insertPolicyEvent, arg.NymID, arg.TransferSequence, arg.PolicyID)
	var i PolicyEvent
	err := row.Scan(&i.NymID, &i.TransferSequence, &i.PolicyID)
	return i, err
}
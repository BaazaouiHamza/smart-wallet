// Code generated by sqlc. DO NOT EDIT.

package repository

import (
	"time"

	"git.digitus.me/prosperus/protocol/identity"
	"git.digitus.me/prosperus/protocol/types"
)

type PolicyEvent struct {
	NymID            identity.PublicKey `json:"nym_id"`
	TransferSequence int64              `json:"transfer_sequence"`
	PolicyID         int64              `json:"policy_id"`
}

type RoutineTransactionPolicy struct {
	ID                int64              `json:"id"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	NymID             identity.PublicKey `json:"nym_id"`
	Recipient         identity.PublicKey `json:"recipient"`
	CreatedAt         time.Time          `json:"created_at"`
	ScheduleStartDate time.Time          `json:"schedule_start_date"`
	ScheduleEndDate   time.Time          `json:"schedule_end_date"`
	Frequency         string             `json:"frequency"`
	Amount            types.Balance      `json:"amount"`
}

type TransactionTriggerPolicy struct {
	ID              int64              `json:"id"`
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	NymID           identity.PublicKey `json:"nym_id"`
	Recipient       identity.PublicKey `json:"recipient"`
	CreatedAt       time.Time          `json:"created_at"`
	TargetedBalance types.Balance      `json:"targeted_balance"`
	Amount          types.Balance      `json:"amount"`
}

type UserPolicy struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	NymID       identity.PublicKey `json:"nym_id"`
	Recipient   string             `json:"recipient"`
	CreatedAt   time.Time          `json:"created_at"`
}

package service

import (
	"context"

	"git.digitus.me/pfe/smart-wallet/types"
	"git.digitus.me/prosperus/protocol/identity"
)

type SmartWallet interface {
	CreateRoutineTransactionPolicy(context.Context, types.RoutineTransactionPolicy) error
	UpdateRoutineTransactionPolicy(context.Context, types.RoutineTransactionPolicy) error
	GetRoutineTransactionPolicy(context.Context, int) (*types.RoutineTransactionPolicy, error)
	ListRoutineTransactionPolicies(
		ctx context.Context, nym identity.PublicKey, page, itemsPerPage int,
	) ([]types.RoutineTransactionPolicy, int, error)

	CreateTransactionTriggerPolicy(context.Context, types.TransactionTriggerPolicy) error
	UpdateTransactionTriggerPolicy(context.Context, types.TransactionTriggerPolicy) error
	GetTransactionTriggerPolicy(context.Context, identity.PublicKey, int) (*types.TransactionTriggerPolicy, error)
	ListTransactionTriggerPolicies(
		ctx context.Context, nym identity.PublicKey, page, itemsPerPage int,
	) ([]types.TransactionTriggerPolicy, int, error)

	CreateUserPolicy(context.Context, types.UserPolicy) error
	UpdateUserPolicy(context.Context, types.UserPolicy) error

	DeletePolicy(context.Context, identity.PublicKey, int) error

	IsNotFoundError(error) bool
	IsUserError(error) bool
}

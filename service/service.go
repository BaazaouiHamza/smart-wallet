package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"git.digitus.me/pfe/smart-wallet/repository"
	"git.digitus.me/pfe/smart-wallet/types"
	"git.digitus.me/prosperus/protocol/identity"
)

type SmartWallet interface {
	CreateTransactionTriggerPolicy(context.Context, types.TransactionTriggerPolicy) error
	UpdateTransactionTriggerPolicy(context.Context, types.TransactionTriggerPolicy) error
	GetTransactionTriggerPolicy(
		context.Context, identity.PublicKey, int,
	) (*types.TransactionTriggerPolicy, error)
	ListTransactionTriggerPolicies(
		ctx context.Context, nym identity.PublicKey, page, itemsPerPage int,
	) ([]types.TransactionTriggerPolicy, error)

	DeletePolicy(context.Context, identity.PublicKey, int) error

	IsNotFoundError(error) bool
	IsUserError(error) bool
}

var _ SmartWallet = (*SmartWalletStd)(nil)

type SmartWalletStd struct{ DB *sql.DB }

func (r *SmartWalletStd) CreateTransactionTriggerPolicy(
	ctx context.Context, ttp types.TransactionTriggerPolicy,
) error {
	balance, amount, err := marshalBoth(ttp.TargetedBalance, ttp.Amount)
	if err != nil {
		return fmt.Errorf("could not marshal balances: %w", err)
	}

	if _, err := repository.New(r.DB).CreateTTP(ctx, repository.CreateTTPParams{
		NymID:           ttp.NymID.String(),
		Name:            ttp.Name,
		Recipient:       ttp.Recipient.String(),
		Description:     ttp.Description,
		TargetedBalance: balance,
		Amount:          amount,
	}); err != nil {
		return fmt.Errorf("could not insert TTP: %w", err)
	}

	return nil
}

func (r *SmartWalletStd) UpdateTransactionTriggerPolicy(
	ctx context.Context, ttp types.TransactionTriggerPolicy,
) error {
	balance, amount, err := marshalBoth(ttp.TargetedBalance, ttp.Amount)
	if err != nil {
		return fmt.Errorf("could not unmarshal balances: %w", err)
	}

	arg := repository.UpdateTTPParams{
		ID:              int64(ttp.ID),
		NymID:           ttp.NymID.String(),
		Recipient:       ttp.Recipient.String(),
		Name:            ttp.Name,
		Description:     ttp.Description,
		TargetedBalance: balance,
		Amount:          amount,
	}

	if _, err := repository.New(r.DB).UpdateTTP(ctx, arg); err != nil {
		return fmt.Errorf("could not update TTP '%d': %w", ttp.ID, err)
	}

	return nil
}

func (r *SmartWalletStd) GetTransactionTriggerPolicy(
	ctx context.Context, pk identity.PublicKey, id int,
) (*types.TransactionTriggerPolicy, error) {
	ttp, err := repository.New(r.DB).GetTTP(ctx, int64(id))
	if err != nil {
		return nil, err
	}

	balance, amount, err := unmarshalBoth(ttp.TargetedBalance, ttp.Amount)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal balances: %w", err)
	}

	return &types.TransactionTriggerPolicy{
		Name:            ttp.Name,
		Description:     ttp.Description,
		NymID:           pk,
		TargetedBalance: balance,
		Amount:          amount,
	}, nil
}

func (r *SmartWalletStd) ListTransactionTriggerPolicies(
	ctx context.Context, nym identity.PublicKey, page, itemsPerPage int,
) ([]types.TransactionTriggerPolicy, error) {
	ttps, err := repository.New(r.DB).ListTTP(ctx, repository.ListTTPParams{
		NymID:  nym.String(),
		Limit:  int32(itemsPerPage),
		Offset: (int32(page) - 1) * int32(itemsPerPage),
	})
	if err != nil {
		return nil, err
	}

	ts := make([]types.TransactionTriggerPolicy, len(ttps))

	for _, ttp := range ttps {
		balance, amount, err := unmarshalBoth(ttp.TargetedBalance, ttp.Amount)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal balances: %w", err)
		}

		ts = append(ts, types.TransactionTriggerPolicy{
			Name:            ttp.Name,
			Description:     ttp.Description,
			NymID:           nym,
			TargetedBalance: balance,
			Amount:          amount,
		})
	}

	return ts, nil

}

func (r *SmartWalletStd) DeletePolicy(ctx context.Context, pk identity.PublicKey, id int) error {
	if err := repository.New(r.DB).DeleteUserPolicy(ctx, int64(id)); err != nil {
		return fmt.Errorf("could not delete policy %d: %w", id, err)
	}

	return nil
}

func (r *SmartWalletStd) IsNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

type userError struct{ err error }

func (uErr userError) Error() string { return uErr.err.Error() }
func (uErr userError) Unwrap() error { return uErr.err }

func (r *SmartWalletStd) IsUserError(err error) bool {
	var uErr userError
	return errors.As(err, &uErr)
}

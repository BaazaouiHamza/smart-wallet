package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"git.digitus.me/pfe/smart-wallet/repository"
	"git.digitus.me/pfe/smart-wallet/types"
	"git.digitus.me/prosperus/protocol/identity"
	ptclTypes "git.digitus.me/prosperus/protocol/types"
)

type SmartWallet interface {
	CreateRoutineTransactionPolicy(context.Context, types.RoutineTransactionPolicy) error
	GetRoutineTransactionPolicy(context.Context, identity.PublicKey, int) (*types.RoutineTransactionPolicy, error)
	UpdateRoutineTransactionPolicy(context.Context, types.RoutineTransactionPolicy) error
	ListRoutineTransactionPolicies(
		ctx context.Context, nym identity.PublicKey, page, itemsPerPage int,
	) ([]types.RoutineTransactionPolicy, int, error)

	CreateTransactionTriggerPolicy(context.Context, types.TransactionTriggerPolicy) error
	UpdateTransactionTriggerPolicy(context.Context, types.TransactionTriggerPolicy) error
	GetTransactionTriggerPolicy(
		context.Context, identity.PublicKey, int,
	) (*types.TransactionTriggerPolicy, error)
	ListTransactionTriggerPolicies(
		ctx context.Context, nym identity.PublicKey, page, itemsPerPage int,
	) ([]types.TransactionTriggerPolicy, int, error)

	DeletePolicy(context.Context, identity.PublicKey, int) error

	DeleteUserPolicy(ctx context.Context, id int) error

	IsNotFoundError(error) bool
	IsUserError(error) bool
}

var _ SmartWallet = (*SmartWalletStd)(nil)

type SmartWalletStd struct{ DB *sql.DB }

func NewSmartWallet(db *sql.DB) SmartWallet {
	return &SmartWalletStd{DB: db}
}

func (r *SmartWalletStd) GetRoutineTransactionPolicy(ctx context.Context, pk identity.PublicKey, id int) (*types.RoutineTransactionPolicy, error) {
	rtp, err := repository.New(r.DB).GetRTP(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	var balance map[ptclTypes.UnitID]int64
	err = json.Unmarshal(rtp.Amount, &balance)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal balance %w", err)
	}
	recipient, err := identity.PublicKeyFromString(rtp.Recipient)
	if err != nil {
		return nil, fmt.Errorf("could not get public key from recipient %w", err)
	}
	return &types.RoutineTransactionPolicy{
		Name:              rtp.Name,
		Description:       rtp.Description,
		ScheduleStartDate: rtp.ScheduleStartDate,
		ScheduleEndDate:   rtp.ScheduleEndDate,
		Amount:            balance,
		Frequency:         rtp.Frequency,
		NymID:             pk,
		Recipient:         *recipient,
	}, nil

}

func (r *SmartWalletStd) CreateRoutineTransactionPolicy(ctx context.Context, rtp types.RoutineTransactionPolicy) error {
	println("fucking thing wont work")
	amount, err := rtp.Amount.MarshalJSON()
	if err != nil {
		return fmt.Errorf("could not marshal amount: %w", err)
	}
	if _, err := repository.New(r.DB).CreateRTP(ctx, repository.CreateRTPParams{
		Name:              rtp.Name,
		Description:       rtp.Description,
		ScheduleStartDate: rtp.ScheduleStartDate,
		ScheduleEndDate:   rtp.ScheduleEndDate,
		Amount:            amount,
		Frequency:         rtp.Frequency,
		NymID:             rtp.NymID.String(),
		Recipient:         rtp.Recipient.String(),
	}); err != nil {
		return fmt.Errorf("could not insert RTP: %w", err)
	}
	return nil
}

func (r *SmartWalletStd) CreateTransactionTriggerPolicy(
	ctx context.Context, ttp types.TransactionTriggerPolicy,
) error {
	println("here")
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
func (r *SmartWalletStd) UpdateRoutineTransactionPolicy(c context.Context, rtp types.RoutineTransactionPolicy) error {
	amount, err := rtp.Amount.MarshalJSON()
	if err != nil {
		return fmt.Errorf("could not marshal amount: %w", err)
	}
	arg := repository.UpdateRTPParams{
		ID:                int64(rtp.ID),
		NymID:             rtp.NymID.String(),
		Recipient:         rtp.Recipient.String(),
		Name:              rtp.Name,
		Description:       rtp.Description,
		Frequency:         rtp.Frequency,
		ScheduleStartDate: rtp.ScheduleStartDate,
		ScheduleEndDate:   rtp.ScheduleEndDate,
		Amount:            amount,
	}
	if _, err := repository.New(r.DB).UpdateRTP(c, arg); err != nil {
		return fmt.Errorf("could not update RTP %w", err)
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

func (r *SmartWalletStd) ListRoutineTransactionPolicies(
	ctx context.Context, nym identity.PublicKey, page, itemsPerPage int,
) ([]types.RoutineTransactionPolicy, int, error) {
	rtps, err := repository.New(r.DB).ListRTP(ctx, repository.ListRTPParams{
		NymID:  nym.String(),
		Limit:  int32(itemsPerPage),
		Offset: (int32(page) - 1) * int32(itemsPerPage),
	})
	if err != nil {
		return nil, 0, nil
	}

	rts := make([]types.RoutineTransactionPolicy, 0, len(rtps))
	total := int(rtps[0].FullCount)
	for _, rtp := range rtps {
		amount, err := unmarshalAmount(rtp.Amount)
		if err != nil {
			return nil, 0, fmt.Errorf("could not unmarshal amount: %w", err)
		}
		recipient, err := identity.PublicKeyFromString(rtp.Recipient)
		if err != nil {
			return nil, 0, fmt.Errorf("could not get public key %w", err)
		}
		rts = append(rts, types.RoutineTransactionPolicy{
			ID:                int(rtp.ID),
			Name:              rtp.Name,
			Description:       rtp.Description,
			NymID:             nym,
			Recipient:         *recipient,
			Amount:            amount,
			Frequency:         rtp.Frequency,
			ScheduleStartDate: rtp.ScheduleStartDate,
			ScheduleEndDate:   rtp.ScheduleEndDate,
		})

	}
	return rts, total, nil
}

func (r *SmartWalletStd) ListTransactionTriggerPolicies(
	ctx context.Context, nym identity.PublicKey, page, itemsPerPage int,
) ([]types.TransactionTriggerPolicy, int, error) {
	ttps, err := repository.New(r.DB).ListTTP(ctx, repository.ListTTPParams{
		NymID:  nym.String(),
		Limit:  int32(itemsPerPage),
		Offset: (int32(page) - 1) * int32(itemsPerPage),
	})
	if err != nil {
		return nil, 0, err
	}
	ts := make([]types.TransactionTriggerPolicy, 0, len(ttps))
	total := int(ttps[0].FullCount)
	for _, ttp := range ttps {
		balance, amount, err := unmarshalBoth(ttp.TargetedBalance, ttp.Amount)
		if err != nil {
			return nil, 0, fmt.Errorf("could not unmarshal balances: %w", err)
		}
		recipient, err := identity.PublicKeyFromString(ttp.Recipient)
		if err != nil {
			return nil, 0, fmt.Errorf("could not get public key %w", err)
		}

		ts = append(ts, types.TransactionTriggerPolicy{
			ID:              int(ttp.ID),
			Name:            ttp.Name,
			Description:     ttp.Description,
			NymID:           nym,
			TargetedBalance: balance,
			Recipient:       *recipient,
			Amount:          amount,
		})
	}

	return ts, total, nil

}

func (r *SmartWalletStd) DeletePolicy(ctx context.Context, pk identity.PublicKey, id int) error {
	if err := repository.New(r.DB).DeleteUserPolicy(ctx, int64(id)); err != nil {
		return fmt.Errorf("could not delete policy %d: %w", id, err)
	}

	return nil
}

func (r *SmartWalletStd) DeleteUserPolicy(ctx context.Context, id int) error {
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

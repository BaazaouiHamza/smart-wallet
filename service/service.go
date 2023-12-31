package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"git.digitus.me/pfe/smart-wallet/internal"
	"git.digitus.me/pfe/smart-wallet/repository"
	"git.digitus.me/pfe/smart-wallet/types"
	"git.digitus.me/prosperus/protocol/identity"
	"github.com/nsqio/go-nsq"
	"go.opentelemetry.io/otel/trace"
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

const otelName = "git.digitus.me/pfe/smart-wallet/service"

var _ SmartWallet = (*SmartWalletStd)(nil)

type SmartWalletStd struct {
	DB        *sql.DB
	publisher *nsq.Producer
}

func NewSmartWallet(db *sql.DB, p *nsq.Producer) SmartWallet {
	return &SmartWalletStd{DB: db, publisher: p}
}

func (r *SmartWalletStd) GetRoutineTransactionPolicy(
	ctx context.Context, pk identity.PublicKey, id int,
) (*types.RoutineTransactionPolicy, error) {

	ctx, span := trace.SpanFromContext(ctx).
		TracerProvider().
		Tracer(otelName).
		Start(ctx, "RTP.Get")
	defer span.End()

	rtp, err := repository.New(r.DB).GetRTP(ctx, repository.GetRTPParams{
		NymID: pk,
		ID:    int64(id),
	})
	if err != nil {
		return nil, err
	}

	return &types.RoutineTransactionPolicy{
		Name:              rtp.Name,
		Description:       rtp.Description,
		ScheduleStartDate: rtp.ScheduleStartDate,
		ScheduleEndDate:   rtp.ScheduleEndDate,
		Amount:            rtp.Amount,
		Frequency:         rtp.Frequency,
		NymID:             pk,
		Recipient:         rtp.Recipient,
	}, nil

}

func (r *SmartWalletStd) CreateRoutineTransactionPolicy(
	ctx context.Context, rtp types.RoutineTransactionPolicy,
) error {
	ctx, span := trace.SpanFromContext(ctx).
		TracerProvider().
		Tracer(otelName).
		Start(ctx, "RTP.Create")
	defer span.End()

	Rtp, err := repository.New(r.DB).CreateRTP(ctx, repository.CreateRTPParams{
		Name:              rtp.Name,
		Description:       rtp.Description,
		ScheduleStartDate: rtp.ScheduleStartDate,
		ScheduleEndDate:   rtp.ScheduleEndDate,
		Amount:            rtp.Amount,
		Frequency:         rtp.Frequency,
		NymID:             rtp.NymID,
		Recipient:         rtp.Recipient,
	})
	if err != nil {
		return fmt.Errorf("could not insert RTP: %w", err)
	}
	data, err := json.Marshal(types.RoutineTransactionPolicy{
		ID:                int(Rtp.ID),
		Description:       rtp.Description,
		Name:              rtp.Name,
		ScheduleStartDate: rtp.ScheduleStartDate,
		ScheduleEndDate:   rtp.ScheduleEndDate,
		Frequency:         rtp.Frequency,
		Amount:            rtp.Amount,
		NymID:             rtp.NymID,
		Recipient:         rtp.Recipient,
		RequestType:       "NEW",
	})
	if err != nil {
		return fmt.Errorf("could not marshal data %w", err)
	}

	if err := r.publisher.Publish(internal.RoutineTransactionPolicyTopic, data); err != nil {
		return err
	}

	return nil
}

func (r *SmartWalletStd) CreateTransactionTriggerPolicy(
	ctx context.Context, ttp types.TransactionTriggerPolicy,
) error {
	println("here")
	if _, err := repository.New(r.DB).CreateTTP(ctx, repository.CreateTTPParams{
		NymID:           ttp.NymID,
		Name:            ttp.Name,
		Recipient:       ttp.Recipient,
		Description:     ttp.Description,
		TargetedBalance: ttp.TargetedBalance,
		Amount:          ttp.Amount,
	}); err != nil {
		return fmt.Errorf("could not insert TTP: %w", err)
	}

	return nil
}
func (r *SmartWalletStd) UpdateRoutineTransactionPolicy(
	ctx context.Context, rtp types.RoutineTransactionPolicy,
) error {
	arg := repository.UpdateRTPParams{
		ID:                int64(rtp.ID),
		NymID:             rtp.NymID,
		Recipient:         rtp.Recipient,
		Name:              rtp.Name,
		Description:       rtp.Description,
		Frequency:         rtp.Frequency,
		ScheduleStartDate: rtp.ScheduleStartDate,
		ScheduleEndDate:   rtp.ScheduleEndDate,
		Amount:            rtp.Amount,
	}
	if _, err := repository.New(r.DB).UpdateRTP(ctx, arg); err != nil {
		return fmt.Errorf("could not update RTP %w", err)
	}
	data, err := json.Marshal(types.RoutineTransactionPolicy{
		ID:                rtp.ID,
		NymID:             rtp.NymID,
		Recipient:         rtp.Recipient,
		Amount:            rtp.Amount,
		ScheduleStartDate: rtp.ScheduleStartDate,
		ScheduleEndDate:   rtp.ScheduleEndDate,
		Frequency:         rtp.Frequency,
		RequestType:       "NEW",
	})
	if err != nil {
		return fmt.Errorf("could not marshal data %w", err)
	}

	if err := r.publisher.Publish(internal.RoutineTransactionPolicyTopic, data); err != nil {
		return err
	}

	return nil
}

func (r *SmartWalletStd) UpdateTransactionTriggerPolicy(
	ctx context.Context, ttp types.TransactionTriggerPolicy,
) error {
	arg := repository.UpdateTTPParams{
		ID:              int64(ttp.ID),
		NymID:           ttp.NymID,
		Recipient:       ttp.Recipient,
		Name:            ttp.Name,
		Description:     ttp.Description,
		TargetedBalance: ttp.TargetedBalance,
		Amount:          ttp.Amount,
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

	return &types.TransactionTriggerPolicy{
		Name:            ttp.Name,
		Description:     ttp.Description,
		NymID:           pk,
		TargetedBalance: ttp.TargetedBalance,
		Amount:          ttp.Amount,
		Recipient:       ttp.Recipient,
	}, nil
}

func (r *SmartWalletStd) ListRoutineTransactionPolicies(
	ctx context.Context, nym identity.PublicKey, page, itemsPerPage int,
) ([]types.RoutineTransactionPolicy, int, error) {
	ctx, span := trace.SpanFromContext(ctx).
		TracerProvider().
		Tracer(otelName).
		Start(ctx, "RTP.Create")
	defer span.End()

	rtps, err := repository.New(r.DB).ListRTP(ctx, repository.ListRTPParams{
		NymID:  nym,
		Limit:  int32(itemsPerPage),
		Offset: int32(page) * int32(itemsPerPage),
	})
	if err != nil {
		return nil, 0, nil
	}

	var (
		total int
		rts   = make([]types.RoutineTransactionPolicy, 0, len(rtps))
	)

	for _, rtp := range rtps {
		total = int(rtp.FullCount)

		rts = append(rts, types.RoutineTransactionPolicy{
			ID:                int(rtp.ID),
			Name:              rtp.Name,
			Description:       rtp.Description,
			NymID:             nym,
			Recipient:         rtp.Recipient,
			Amount:            rtp.Amount,
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
	ctx, span := trace.SpanFromContext(ctx).
		TracerProvider().
		Tracer(otelName).
		Start(ctx, "RTP.Create")
	defer span.End()

	ttps, err := repository.New(r.DB).ListTTP(ctx, repository.ListTTPParams{
		NymID:  nym,
		Limit:  int32(itemsPerPage),
		Offset: int32(page) * int32(itemsPerPage),
	})
	if err != nil {
		return nil, 0, err
	}

	ts := make([]types.TransactionTriggerPolicy, 0, len(ttps))
	var total int
	for _, ttp := range ttps {
		total = int(ttp.FullCount)

		ts = append(ts, types.TransactionTriggerPolicy{
			ID:              int(ttp.ID),
			Name:            ttp.Name,
			Description:     ttp.Description,
			NymID:           nym,
			TargetedBalance: ttp.TargetedBalance,
			Recipient:       ttp.Recipient,
			Amount:          ttp.Amount,
		})
	}

	return ts, total, nil

}

func (r *SmartWalletStd) DeletePolicy(ctx context.Context, pk identity.PublicKey, id int) error {
	if err := repository.New(r.DB).DeleteUserPolicy(ctx, int64(id)); err != nil {
		return fmt.Errorf("could not delete policy %d: %w", id, err)
	}
	data, err := json.Marshal(types.RoutineTransactionPolicy{
		ID:          id,
		RequestType: "DELETE",
	})
	if err != nil {
		return fmt.Errorf("could not marshal data %w", err)
	}

	if err = r.publisher.Publish(internal.RoutineTransactionPolicyTopic, data); err != nil {
		return err
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

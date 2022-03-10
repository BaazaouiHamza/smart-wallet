package service

import (
	"context"
	"database/sql"
	"encoding/json"

	db "git.digitus.me/pfe/smart-wallet/db/sqlc"
	"git.digitus.me/pfe/smart-wallet/types"
	"git.digitus.me/prosperus/protocol/identity"
	ptclTypes "git.digitus.me/prosperus/protocol/types"
)

type SmartWallet interface {
	// CreateRoutineTransactionPolicy(context.Context, types.RoutineTransactionPolicy, identity.PublicKey) error
	// UpdateRoutineTransactionPolicy(context.Context, types.RoutineTransactionPolicy) error
	// GetRoutineTransactionPolicy(context.Context, int) (*types.RoutineTransactionPolicy, error)
	// ListRoutineTransactionPolicies(
	// 	ctx context.Context, nym identity.PublicKey, page, itemsPerPage int,
	// ) ([]types.RoutineTransactionPolicy, int, error)

	CreateTransactionTriggerPolicy(context.Context, types.TransactionTriggerPolicy) error
	UpdateTransactionTriggerPolicy(context.Context, types.TransactionTriggerPolicy) error
	GetTransactionTriggerPolicy(context.Context, identity.PublicKey, int) (*types.TransactionTriggerPolicy, error)
	ListTransactionTriggerPolicies(
		ctx context.Context, nym identity.PublicKey, page, itemsPerPage int,
	) ([]types.TransactionTriggerPolicy, error)

	DeletePolicy(context.Context, identity.PublicKey, int) error

	IsNotFoundError(error) bool
	IsUserError(error) bool
}

// func NewSmartWallet(db *sql.DB) SmartWallet {
// 	return &repoSvc{
// 		Queries: New(db),
// 		db:      db,
// 	}
// }

type repoSvc struct {
	*db.Queries
	db *sql.DB
}

func (r *repoSvc) CreateTransactionTriggerPolicy(c context.Context, ttp types.TransactionTriggerPolicy) error {
	balance, err := ttp.TargetedBalance.MarshalJSON()
	if err != nil {
		return err
	}
	amount, err := ttp.Amount.MarshalJSON()
	if err != nil {
		return err
	}
	arg := db.CreateTTPParams{
		NymID:           ttp.NymID.String(),
		Name:            ttp.Name,
		Recipient:       ttp.Recipient.String(),
		Description:     ttp.Description,
		TargetedBalance: balance,
		Amount:          amount,
	}
	_, err = r.CreateTTP(c, arg)
	if err != nil {
		return err
	}
	return nil
}

func (r *repoSvc) UpdateTransactionTriggerPolicy(c context.Context, ttp types.TransactionTriggerPolicy) error {
	balance, err := ttp.TargetedBalance.MarshalJSON()
	if err != nil {
		return err
	}
	amount, err := ttp.Amount.MarshalJSON()
	if err != nil {
		return err
	}
	arg := db.UpdateTTPParams{
		ID:              int64(ttp.ID),
		NymID:           ttp.NymID.String(),
		Recipient:       ttp.Recipient.String(),
		Name:            ttp.Name,
		Description:     ttp.Description,
		TargetedBalance: balance,
		Amount:          amount,
	}
	_, err = r.UpdateTTP(c, arg)
	if err != nil {
		return err
	}
	return nil
}

func (r *repoSvc) GetTransactionTriggerPolicy(c context.Context, pk identity.PublicKey, id int) (*types.TransactionTriggerPolicy, error) {
	ttp, err := r.GetTTP(c, int64(id))
	if err != nil {
		return nil, err
	}
	var balance map[ptclTypes.UnitID]int64
	err = json.Unmarshal(ttp.TargetedBalance, &balance)
	if err != nil {
		return nil, err
	}
	var amount map[ptclTypes.UnitID]int64
	err = json.Unmarshal(ttp.TargetedBalance, &amount)
	if err != nil {
		return nil, err
	}
	Ttp := types.TransactionTriggerPolicy{
		Name:            ttp.Name,
		Description:     ttp.Description,
		NymID:           pk,
		TargetedBalance: balance,
		Amount:          amount,
	}
	return &Ttp, nil

}

func (r *repoSvc) ListTransactionTriggerPolicies(
	c context.Context, nym identity.PublicKey, page, itemsPerPage int,
) ([]types.TransactionTriggerPolicy, error) {
	arg := db.ListTTPParams{
		NymID:  nym.String(),
		Limit:  int32(itemsPerPage),
		Offset: (int32(page) - 1) * int32(itemsPerPage),
	}
	ttps, err := r.ListTTP(c, arg)
	if err != nil {
		return nil, err
	}
	var amount []map[ptclTypes.UnitID]int64
	var balance []map[ptclTypes.UnitID]int64
	var Ttps []types.TransactionTriggerPolicy
	for i, ttp := range ttps {

		err = json.Unmarshal(ttp.TargetedBalance, &balance[i])
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(ttp.TargetedBalance, &amount[i])
		if err != nil {
			return nil, err
		}
		Ttps[i] = types.TransactionTriggerPolicy{
			Name:            ttp.Name,
			Description:     ttp.Description,
			NymID:           nym,
			TargetedBalance: balance[i],
			Amount:          amount[i],
		}
	}

	return Ttps, nil

}
func (r *repoSvc) DeletePolicy(c context.Context, pk identity.PublicKey, id int) error {
	err := r.Queries.DeleteUserPolicy(c, int64(id))
	if err != nil {
		return err
	}
	return nil
}

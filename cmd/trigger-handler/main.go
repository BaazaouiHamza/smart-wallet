package main

import (
	"context"
	"database/sql"
	"errors"

	"git.digitus.me/library/prosper-kit/wallet"
	walletrepository "git.digitus.me/library/prosper-kit/wallet/repository"
	"git.digitus.me/pfe/smart-wallet/types"
	"git.digitus.me/prosperus/protocol/identity"
	ptclTypes "git.digitus.me/prosperus/protocol/types"
)

type TriggerHandler interface {
	HandleNotarization(context.Context, *ptclTypes.Notarization) error
	HandleTrigger(context.Context, types.TriggerMessage)
}

type Cloudwallet interface {
	SendTransaction(
		context.Context, map[identity.PublicKey]ptclTypes.Balance,
	) (*ptclTypes.RawTransfer, error)
	GetUserState(context.Context, identity.PublicKey) (*ptclTypes.UserState, error)
}

type ITriggerHandler struct {
	db          *sql.DB
	cloudwallet Cloudwallet
	dcn         wallet.DCN
}

func (th *ITriggerHandler) HandleNotarization(
	ctx context.Context, n *ptclTypes.Notarization,
) error {
	if !n.TransactionResult.Success() {
		// TODO: get policy_event by nymID, TransferSequence pairs
		// SELECT * FROM policy_events WHERE (nym_id, transfer_sequence) = ANY (@something)
		// If any exists republish in trigger topic
		return nil
	} else {
		var matchingPolicies []types.TransactionTriggerPolicy
		// TODO: get matching policies by peers

		nyms := make([]string, 0, len(matchingPolicies))
		for _, p := range matchingPolicies {
			nyms = append(nyms, p.NymID.String())
		}

		// TODO: use errgroup to run all these requests in parallel
		usLookup := map[identity.PublicKey]ptclTypes.Balance{}
		for _, m := range matchingPolicies {
			// BUG: we might be racing the cloudwallet here, it would be wise to check if
			// the UserState contains the current notarization or not.
			us, err := th.cloudwallet.GetUserState(ctx, m.NymID)
			if err != nil {
				return err
			}

      var balancePre,balancePost ptclTypes.Balance

      // for trigger to trigger balancePre < m.TargetedBalance
      // and balancePost > m.TargetedBalance

			if ts := n.Peers[us.UserNymID].TransferSequence; ts != nil {
				if us.TransferSequences.Contains(*ts){
          // has notarization
        } else {
          // does not have notarization
          // calculate balance yourself
        }
			}


			usLookup[m.NymID] = us.Balance
		}

		for range matchingPolicies {
			// TODO: compare p.TargetedBalance to usLookup[p.NymID],
			// if inferior then publish message or handle directly
		}

	}

	return nil
}

func (th *ITriggerHandler) HandleTrigger(
	ctx context.Context, m types.TriggerMessage,
) error {
	dcnNymID, err := th.dcn.NymID(ctx)
	if err != nil {
		return err
	}

	t := ptclTypes.RawTransfer{
		DCNNymID: dcnNymID,
		// TODO: place policy in context or something
		Context: "some context",
		Peers:   ptclTypes.Peers{},
	}

	for pk, a := range m.Amounts {
		t.Peers[pk] = ptclTypes.Peer{TransferQuantity: a}
	}

	// BUG: fees not being set

	res, err := th.cloudwallet.SendTransaction(ctx, m.Amounts)
	if err != nil {
		return err
	}

	if res.TransactionResult != nil && !res.TransactionResult.Success() {
		return errors.New("transaction failed")
	}

	// Insert into policy_events

	return nil
}

func main() {
	// This is where we read events and launch transactions according to policies
}

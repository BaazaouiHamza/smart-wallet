package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	prospercontext "git.digitus.me/library/prosper-kit/context"
	"git.digitus.me/library/prosper-kit/wallet"
	"git.digitus.me/pfe/smart-wallet/internal"
	"git.digitus.me/pfe/smart-wallet/repository"
	helpers "git.digitus.me/pfe/smart-wallet/service"
	"git.digitus.me/pfe/smart-wallet/types"
	"git.digitus.me/prosperus/protocol/identity"
	ptclTypes "git.digitus.me/prosperus/protocol/types"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Cloudwallet interface {
	SendTransaction(
		context.Context, map[identity.PublicKey]ptclTypes.Balance,
	) (*ptclTypes.RawTransfer, error)
	GetUserState(context.Context, identity.PublicKey) (*ptclTypes.UserState, error)
}

type TriggerHandler interface {
	HandleNotarization(context.Context, *ptclTypes.Notarization) error
	HandleTrigger(context.Context, types.TriggerMessage) error
}

type TtpHandler struct {
	db          *sql.DB
	cloudwallet Cloudwallet
	dcn         wallet.DCN
	publisher   *nsq.Producer
}

var _ TriggerHandler = (*TtpHandler)(nil)

func NewTriggerHandler(db *sql.DB, p *nsq.Producer) TriggerHandler {
	return &TtpHandler{
		db:        db,
		publisher: p,
	}
}

func (th *TtpHandler) HandleNotarization(
	ctx context.Context, n *ptclTypes.Notarization,
) error {
	logger := prospercontext.GetLogger(ctx)
	if !n.TransactionResult.Success() {
		// TODO: get policy_event by nymID, TransferSequence pairs
		// SELECT * FROM policy_events WHERE (nym_id, transfer_sequence) = ANY (@something)
		// If any exists republish in trigger topic
		for peerPk, peer := range n.Peers {
			pe, err := repository.New(th.db).GetPolicyEvent(ctx, repository.GetPolicyEventParams{
				NymID:            peerPk,
				TransferSequence: *peer.TransferSequence,
			})
			if err != nil {
				return err
			}
			ttp, err := repository.New(th.db).GetTTP(ctx, pe.PolicyID)
			if err != nil {
				return err
			}
			data, err := json.Marshal(types.TriggerMessage{
				PolicyID: int(ttp.ID),
				Amounts: map[identity.PublicKey]ptclTypes.Balance{
					ttp.NymID:     helpers.NegativeAmount(ttp.Amount),
					ttp.Recipient: ttp.Amount,
				},
			})
			if err != nil {
				return fmt.Errorf("could not marshal %w", err)
			}
			if err := th.publisher.Publish(internal.TransactionsTopic, data); err != nil {
				logger.Error("could not publish message", zap.Any("Trigger Message", data), zap.Error(err))
			}
		}
		return nil
	} else {
		var matchingPolicies []types.TransactionTriggerPolicy
		// TODO: get matching policies by peers
		for peerPk := range n.Peers {
			ttps, err := repository.New(th.db).ListMatchingPolicies(ctx, peerPk)
			if err != nil {
				return err
			}
			for _, ttp := range ttps {
				matchingPolicies = append(matchingPolicies, types.TransactionTriggerPolicy{
					Name:            ttp.Name,
					Description:     ttp.Description,
					NymID:           ttp.NymID,
					Recipient:       ttp.Recipient,
					TargetedBalance: ttp.TargetedBalance,
					Amount:          ttp.Amount,
				})
			}
		}

		nyms := make([]identity.PublicKey, 0, len(matchingPolicies))
		for _, p := range matchingPolicies {
			nyms = append(nyms, p.NymID)
		}
		// TODO: use errgroup to run all these requests in parallel
		usLookup := map[identity.PublicKey]ptclTypes.Balance{}
		g, ctx := errgroup.WithContext(ctx)
		for _, m := range matchingPolicies {
			currentMatchingPolicy := m
			g.Go(func() error {
				// BUG: we might be racing the cloudwallet here, it would be wise to check if
				// the UserState contains the current notarization or not.
				us, err := th.cloudwallet.GetUserState(ctx, currentMatchingPolicy.NymID)
				if err != nil {
					return err
				}
				// for trigger to trigger balancePre < m.TargetedBalance
				// and balancePost > m.TargetedBalance
				var balancePre, balancePost ptclTypes.Balance
				if ts := n.Peers[us.UserNymID].TransferSequence; ts != nil {
					if us.TransferSequences.Contains(*ts) {
						balancePost = us.Balance
						err := us.Balance.Add(helpers.NegativeAmount(n.Peers[us.UserNymID].TransferQuantity))
						if err != nil {
							return err
						}
						balancePre = us.Balance
					} else {
						balancePre = us.Balance
						err := us.Balance.Add(helpers.NegativeAmount(n.Peers[us.UserNymID].TransferQuantity))
						if err != nil {
							return err
						}
						balancePost = us.Balance
					}
				}
				for unitID, targetedBalance := range currentMatchingPolicy.TargetedBalance {
					if balancePost[unitID] > targetedBalance && balancePre[unitID] < targetedBalance {
						data, err := json.Marshal(types.TriggerMessage{
							PolicyID: currentMatchingPolicy.ID,
							Amounts: map[identity.PublicKey]ptclTypes.Balance{
								currentMatchingPolicy.NymID:     helpers.NegativeAmount(currentMatchingPolicy.Amount),
								currentMatchingPolicy.Recipient: currentMatchingPolicy.Amount,
							},
						})
						if err != nil {
							return fmt.Errorf("could not marshal data %w", err)
						}
						if err := th.publisher.Publish(internal.TransactionsTopic, data); err != nil {
							logger.Error("could not publish message", zap.Any("triggerPolicy", currentMatchingPolicy), zap.Error(err))
						}
					}

				}

				usLookup[currentMatchingPolicy.NymID] = us.Balance

				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return err
		}

		for _, p := range matchingPolicies {
			//Here We Check matching policies by peers if any exist that satisfy TargetedBalance
			// TODO: compare p.TargetedBalance to usLookup[p.NymID],
			// if inferior then publish message or handle directly
			for unitID, amount := range p.TargetedBalance {
				if usLookup[p.NymID][unitID] > amount {
					data, err := json.Marshal(types.TriggerMessage{
						PolicyID: p.ID,
						Amounts: map[identity.PublicKey]ptclTypes.Balance{
							p.NymID:     helpers.NegativeAmount(p.Amount),
							p.Recipient: p.Amount,
						},
					})
					if err != nil {
						return fmt.Errorf("could not marshal data %w", err)
					}
					if err := th.publisher.Publish(internal.TransactionsTopic, data); err != nil {
						logger.Error("could not publish message", zap.Any("triggerPolicy", p), zap.Error(err))
					}
				}
			}

		}

	}

	return nil
}

func (th *TtpHandler) HandleTrigger(
	ctx context.Context, m types.TriggerMessage,
) error {
	dcnNymID, err := th.dcn.NymID(ctx)
	if err != nil {
		return err
	}
	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("could not marshal data %w", err)
	}
	t := ptclTypes.RawTransfer{
		DCNNymID: dcnNymID,
		// TODO: place policy in context or something
		Context: string(data),
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

	for peerPk := range res.Peers {
		_, err = repository.New(th.db).InsertPolicyEvent(ctx, repository.InsertPolicyEventParams{
			PolicyID:         int64(m.PolicyID),
			TransferSequence: *res.Peers[peerPk].TransferSequence,
			NymID:            peerPk,
		})
		if err != nil {
			return errors.New("could not insert policy event")
		}
	}

	return nil
}

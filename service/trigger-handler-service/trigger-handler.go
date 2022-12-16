package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	prospercontext "git.digitus.me/library/prosper-kit/context"
	"git.digitus.me/pfe/smart-wallet/internal"
	"git.digitus.me/pfe/smart-wallet/repository"
	helpers "git.digitus.me/pfe/smart-wallet/service"
	"git.digitus.me/pfe/smart-wallet/types"
	"git.digitus.me/prosperus/protocol/identity"
	ptclTypes "git.digitus.me/prosperus/protocol/types"
	"github.com/google/uuid"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type TriggerHandler interface {
	HandleNotarization(context.Context, *ptclTypes.Notarization) error
	HandleTrigger(context.Context, types.TriggerMessage) error
}

type TtpHandler struct {
	db          *sql.DB
	cloudwallet internal.CloudwalletClient
	publisher   *nsq.Producer
}

var _ TriggerHandler = (*TtpHandler)(nil)

func NewTriggerHandler(
	db *sql.DB, p *nsq.Producer, cloudwallet *internal.CloudwalletClient,
) TriggerHandler {
	return &TtpHandler{
		db:          db,
		publisher:   p,
		cloudwallet: *cloudwallet,
	}
}

func (th *TtpHandler) retryTransactionIfMatchesPolicy(
	ctx context.Context, n *ptclTypes.Notarization,
) error {
	logger := prospercontext.GetLogger(ctx)

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
				ttp.Recipient: ttp.Amount,
			},
		})
		if err != nil {
			return err
		}

		if err := th.publisher.Publish(internal.TransactionsTopic, data); err != nil {
			logger.Error("could not publish message", zap.Any("Trigger Message", data), zap.Error(err))
			return err
		}
		logger.Debug("Retry Matching Policy and Publish", zap.Any("Trigger Message", data))
	}

	return nil
}

func (th *TtpHandler) HandleNotarization(
	ctx context.Context, n *ptclTypes.Notarization,
) error {
	logger := prospercontext.GetLogger(ctx)
	if !n.TransactionResult.Success() {
		err := th.retryTransactionIfMatchesPolicy(ctx, n)
		if err != nil {
			return err
		}
	}

	var (
		matchingPolicies []types.TransactionTriggerPolicy
		nyms             = make([]string, 0, len(n.Peers))
	)

	for peerPK := range n.Peers {
		nyms = append(nyms, peerPK.String())
	}

	ttps, err := repository.New(th.db).ListMatchingPoliciesBatch(ctx, nyms)
	if err != nil {
		return err
	}
	for _, ttp := range ttps {
		matchingPolicies = append(matchingPolicies, types.TransactionTriggerPolicy{
			ID:              int(ttp.ID),
			Name:            ttp.Name,
			Description:     ttp.Description,
			NymID:           ttp.NymID,
			Recipient:       ttp.Recipient,
			TargetedBalance: ttp.TargetedBalance,
			Amount:          ttp.Amount,
		})
	}
	pks := make([]identity.PublicKey, 0, len(matchingPolicies))
	for _, p := range matchingPolicies {
		pks = append(pks, p.NymID)
	}
	// TODO: use errgroup to run all these requests in parallel
	usLookup := map[identity.PublicKey]ptclTypes.Balance{}
	var usLookupLock sync.Mutex

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
							currentMatchingPolicy.Recipient: currentMatchingPolicy.Amount,
						},
					})
					if err != nil {
						return fmt.Errorf("could not marshal data %w", err)
					}
					if err := th.publisher.Publish(internal.TransactionsTopic, data); err != nil {
						logger.Error(
							"could not publish message",
							zap.Any("triggerPolicy", currentMatchingPolicy),
							zap.Error(err),
						)
					}
					logger.Debug("Published Trigger Policy From Matching Policy", zap.Any("Trigger Message", data))
				}

			}

			usLookupLock.Lock()
			defer usLookupLock.Unlock()
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
						p.Recipient: p.Amount,
					},
				})
				if err != nil {
					return fmt.Errorf("could not marshal data %w", err)
				}
				if err := th.publisher.Publish(internal.TransactionsTopic, data); err != nil {
					logger.Error("could not publish message", zap.Any("triggerPolicy", p), zap.Error(err))
				}
				logger.Debug("Published Trigger Policy From UsLookUp", zap.Any("Trigger Message", data))
			}
		}

	}

	return nil
}

func (th *TtpHandler) HandleTrigger(
	ctx context.Context, m types.TriggerMessage,
) error {
	policy, err := repository.New(th.db).GetUserPolicy(ctx, int64(m.PolicyID))
	if err != nil {
		return err
	}

	// BUG: fees not being set
	res, err := th.cloudwallet.SendAmounts(
		ctx,
		policy.NymID,
		m.Amounts,
		fmt.Sprintf("policy-%d-%s", m.PolicyID, uuid.New()),
	)
	if err != nil {
		return err
	}
	if _, err = repository.New(th.db).InsertPolicyEvent(ctx, repository.InsertPolicyEventParams{
		PolicyID:         int64(m.PolicyID),
		TransferSequence: *res.Peers[policy.NymID].TransferSequence,
		NymID:            policy.NymID,
	}); err != nil {
		return fmt.Errorf("could not insert policy event: %w", err)
	}

	return nil
}

package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	prospercontext "git.digitus.me/library/prosper-kit/context"
	"git.digitus.me/pfe/smart-wallet/repository"
	"git.digitus.me/pfe/smart-wallet/types"
	"git.digitus.me/prosperus/protocol/identity"
	ptclTypes "git.digitus.me/prosperus/protocol/types"
	"git.digitus.me/prosperus/publisher"
	"github.com/nsqio/go-nsq"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type RTPTicker interface {
	HandleNewPolicy(context.Context, types.RoutineTransactionPolicy) error
	HandleDeletePolicy(context.Context, int) error
	StartAllRoutinePolicies(context.Context) error
}

type InMemory struct {
	routinePolicy map[int]cron.EntryID
	mutex         sync.RWMutex
	db            *sql.DB
	scheduler     *cron.Cron
	publisher     *nsq.Producer
}

var _ (RTPTicker) = (*InMemory)(nil)

func NewRTPTicker(db *sql.DB, s *cron.Cron, p *nsq.Producer) RTPTicker {
	return &InMemory{
		db:            db,
		routinePolicy: map[int]cron.EntryID{},
		mutex:         sync.RWMutex{},
		scheduler:     s,
		publisher:     p,
	}
}

func (im *InMemory) runCronJobs(
	ctx context.Context, rtp types.RoutineTransactionPolicy,
) error {
	logger := prospercontext.GetLogger(ctx)

	// TODO: use actual frequency from rtp
	frequency := "every 5s"

	negativeAmount := ptclTypes.Balance{}

	for u, a := range rtp.Amount {
		negativeAmount[u] = -a
	}

	data, err := json.Marshal(types.TriggerMessage{
		PolicyID: rtp.ID,
		Amounts: map[identity.PublicKey]ptclTypes.Balance{
			rtp.NymID:     negativeAmount,
			rtp.Recipient: rtp.Amount,
		},
	})
	if err != nil {
		return fmt.Errorf("could not marshal data %w", err)
	}

	entryId, err := im.scheduler.AddFunc("@"+frequency, func() {
		now := time.Now()
		if now.After(rtp.ScheduleStartDate) && now.Before(rtp.ScheduleEndDate) {
			if err := im.publisher.Publish("transactions", data); err != nil {
				logger.Error("could not publish message", zap.Any("policy", rtp), zap.Error(err))
			}
		}
	})
	if err != nil {
		return err
	}

	func() {
		im.mutex.Lock()
		defer im.mutex.Unlock()
		im.routinePolicy[rtp.ID] = entryId
		logger.Debug("", zap.Any("RoutinePolicyMap", im.routinePolicy))
	}()

	logger.Debug("", zap.Any("RoutinePolicy", rtp))

	return nil
}

func HandleRTPMessages(rtpTicker RTPTicker) func(ctx context.Context, m *nsq.Message) error {
	return func(ctx context.Context, msg *nsq.Message) error {
		logger := prospercontext.GetLogger(ctx)
		var data types.RoutineTransactionPolicy
		if err := publisher.Decode(msg.Body, data); err != nil {
			return err
		}

		switch data.RequestType {
		case "NEW":
			return rtpTicker.HandleNewPolicy(ctx, data)
		case "DELETE":
			return rtpTicker.HandleDeletePolicy(ctx, data.ID)
		default:
			// log a warning about unknown type
			logger.Warn("unknown request type", zap.Any("message", data))
			return nil
		}
	}
}

func (im *InMemory) HandleNewPolicy(ctx context.Context, rtp types.RoutineTransactionPolicy) error {
	if err := im.HandleDeletePolicy(ctx, rtp.ID); err != nil {
		return err
	}

	im.runCronJobs(ctx, rtp)

	return nil
}

func (im *InMemory) HandleDeletePolicy(ctx context.Context, id int) error {
	im.mutex.Lock()
	defer im.mutex.Unlock()
	if entrydId, ok := im.routinePolicy[id]; ok {
		im.scheduler.Remove(entrydId)
		delete(im.routinePolicy, id)
	}

	return nil
}

func (im *InMemory) StartAllRoutinePolicies(ctx context.Context) error {
	rtps, err := repository.New(im.db).GetALlRoutinePolicies(ctx)
	if err != nil {
		return err
	}

	for _, rtp := range rtps {
		im.runCronJobs(ctx, types.RoutineTransactionPolicy{
			ID:                int(rtp.ID),
			NymID:             rtp.NymID,
			Recipient:         rtp.Recipient,
			Amount:            rtp.Amount,
			Frequency:         rtp.Frequency,
			ScheduleStartDate: rtp.ScheduleStartDate,
			ScheduleEndDate:   rtp.ScheduleEndDate,
			RequestType:       "NEW",
		})
	}

	return nil

}

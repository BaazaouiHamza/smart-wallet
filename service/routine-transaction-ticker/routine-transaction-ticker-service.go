package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	prospercontext "git.digitus.me/library/prosper-kit/context"
	"git.digitus.me/pfe/smart-wallet/repository"
	"git.digitus.me/pfe/smart-wallet/types"
	"git.digitus.me/prosperus/publisher"
	"github.com/nsqio/go-nsq"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type RTPTicker interface {
	HandleNewPolicy(context.Context, types.RoutineTransactionPolicy) error
	HandleDeletePolicy(context.Context, int) error
	StartAllRoutinePolicies(context.Context) error
	HandleRTPMessages(rtpTicker RTPTicker) func(context.Context, *nsq.Message) error
}

type InMemory struct {
	RoutinePolicy map[int]cron.EntryID
	Mu            sync.Mutex
	DB            *sql.DB
	S             *cron.Cron
	P             *nsq.Producer
}

var _ (RTPTicker) = (*InMemory)(nil)

func NewRTPTicker(db *sql.DB, s *cron.Cron, p *nsq.Producer) RTPTicker {
	return &InMemory{DB: db, RoutinePolicy: make(map[int]cron.EntryID), Mu: sync.Mutex{}, S: s, P: p}
}

func (im *InMemory) runCronJobs(ctx context.Context, s *cron.Cron, rtp types.RoutineTransactionPolicy) error {
	logger := prospercontext.GetLogger(ctx)
	var frequency = "every 5s"
	data, err := json.Marshal(rtp)
	if err != nil {
		return fmt.Errorf("could not marshal data %w", err)
	}
	entryId, err := s.AddFunc("@"+frequency, func() {
		if time.Now().After(rtp.ScheduleStartDate) && time.Now().Before(rtp.ScheduleEndDate) {
			err = im.P.Publish("transactions", data)
			if err != nil {
				logger.Error("could not Publish")
			}
		}
	})
	if err != nil {
		return err
	}
	func() {
		im.Mu.Lock()
		defer im.Mu.Unlock()
		im.RoutinePolicy[rtp.ID] = entryId
	}()
	// never use fmt package for logging use logger from context with prospercontext.GetLogger
	logger.Debug("", zap.String("RoutinePolicyMap", fmt.Sprint(im.RoutinePolicy)))
	// never use fmt package for logging use logger from context with prospercontext.GetLogger
	logger.Debug("", zap.String("RoutinePolicy", fmt.Sprint(rtp)))
	return nil
}

func (im *InMemory) HandleRTPMessages(rtpTicker RTPTicker) func(ctx context.Context, m *nsq.Message) error {
	return func(ctx context.Context, msg *nsq.Message) error {
		logger := prospercontext.GetLogger(ctx)
		var data types.RoutineTransactionPolicy
		err := publisher.Decode(msg.Body, data)
		if err != nil {
			return err
		}
		switch data.RequestType {
		case "NEW":
			return rtpTicker.HandleNewPolicy(ctx, data)
		case "DELETE":
			return rtpTicker.HandleDeletePolicy(ctx, data.ID)
		default:
			// log a warning about unknown type
			logger.Warn("unknown request type")
			return nil
		}
	}
}

func (im *InMemory) HandleNewPolicy(ctx context.Context, rtp types.RoutineTransactionPolicy) error {
	if entryId, ok := im.RoutinePolicy[rtp.ID]; ok {
		im.S.Remove(entryId)
		delete(im.RoutinePolicy, rtp.ID)
		fmt.Println(im.RoutinePolicy)
		im.runCronJobs(ctx, im.S, rtp)
	} else {
		im.runCronJobs(ctx, im.S, rtp)
	}
	return nil
}

func (im *InMemory) HandleDeletePolicy(ctx context.Context, id int) error {
	if entrydId, ok := im.RoutinePolicy[id]; ok {
		im.S.Remove(entrydId)
		delete(im.RoutinePolicy, id)
	} else {
		return errors.New("could not find entryId for the given id")
	}

	return nil
}
func (im *InMemory) StartAllRoutinePolicies(ctx context.Context) error {
	rtps, err := repository.New(im.DB).GetALlRoutinePolicies(ctx)
	if err != nil {
		return err
	}
	rts := make([]types.RoutineTransactionPolicy, 0, len(rtps))

	for i, rtp := range rtps {
		rts = append(rts, types.RoutineTransactionPolicy{
			ID:                int(rtp.ID),
			NymID:             rtp.NymID,
			Recipient:         rtp.Recipient,
			Amount:            rtp.Amount,
			Frequency:         rtp.Frequency,
			ScheduleStartDate: rtp.ScheduleStartDate,
			ScheduleEndDate:   rtp.ScheduleEndDate,
			RequestType:       "NEW",
		})
		im.runCronJobs(ctx, im.S, rts[i])
	}

	return nil

}

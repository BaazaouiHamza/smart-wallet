package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"git.digitus.me/library/prosper-kit/application"
	"git.digitus.me/pfe/smart-wallet/types"
	"git.digitus.me/prosperus/publisher"
	"github.com/nsqio/go-nsq"
	"github.com/robfig/cron/v3"
)

// This should be in its own package
type inMemory struct {
	routinePolicy map[int]cron.EntryID
	mu            sync.Mutex
}

type RTPTicker interface {
	HandleNewPolicy(context.Context, types.RoutineTransactionPolicy) error
	HandleDeletePolicy(context.Context, int) error
}

var _ (RTPTicker) = (*inMemory)(nil)

func HandleRTPMessages(rtpTicker RTPTicker) func(context.Context, *nsq.Message) error {
	return func(ctx context.Context, msg *nsq.Message) error {
    var type string
    // here we deserialize and hande to the approriate method
    switch type {
    case "NEW":
      return rtpTrtpTicker.HandleNewPolicy(ctx, types.RoutineTransactionPolicy{})
    case "DELETE":
      return rtpTrtpTicker.HandleDeletePolicy(ctx, 0)
    default:
      // log a warning about unknown type
      return nil
    }
	}
}

var inMemoryRoutinePolicy = &inMemory{routinePolicy: make(map[int]cron.EntryID)}

func runCronJobs(s *cron.Cron, rtp types.RoutineTransactionPolicy) {
	var frequency = "every 5s"

	entryId, err := s.AddFunc("@"+frequency, func() {
		if time.Now().After(rtp.ScheduleStartDate) && time.Now().Before(rtp.ScheduleEndDate) {
			// never use fmt package for logging use logger from context with prospercontext.GetLogger
			fmt.Println("Every hour on the half hour")
		}
	})
	if err != nil {
		// never use log package for logging use logger from context with prospercontext.GetLogger
		log.Println(err)
	}

	entry := s.Entry(entryId)
	// never use fmt package for logging use logger from context with prospercontext.GetLogger
	fmt.Println(entry)
	func() {
		inMemoryRoutinePolicy.mu.Lock()
		defer inMemoryRoutinePolicy.mu.Unlock()
		inMemoryRoutinePolicy.routinePolicy[rtp.ID] = entryId
	}()
	// never use fmt package for logging use logger from context with prospercontext.GetLogger
	fmt.Println(inMemoryRoutinePolicy.routinePolicy)
	// never use fmt package for logging use logger from context with prospercontext.GetLogger
	fmt.Println(rtp.RequestType)
	s.Start()
}

func run(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	s := cron.New()

	wg.Add(1)
	s.Start()

	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	defer func() {
		defer wg.Done()
		<-s.Stop().Done()
	}()

	config := publisher.NewNSQConfig()
	c, err := publisher.NewConsumer(
		context.Background(), "routine-transaction-policies", "ticker", config,
	)
	if err != nil {
		return err
	}

  var rtpTicker RTPTicker 
  {
    // Here we initialize the ticker
  }

	// handler func should be a method on inMemory
	c.AddHandler(HandleRTPMessages(rtp))

	// Should be read from config
	err = c.ConnectToNSQLookupd("127.0.0.1:4150")
	if err != nil {
		return err
	}

	return nil
}

func main() {
	application.WaitForInterrupt(run)
}

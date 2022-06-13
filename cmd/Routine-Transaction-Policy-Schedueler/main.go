package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"git.digitus.me/pfe/smart-wallet/types"
	"git.digitus.me/prosperus/publisher"
	"github.com/nsqio/go-nsq"
	"github.com/robfig/cron/v3"
)

type inMemory struct {
	routinePolicy map[int]cron.EntryID
}

var inMemoryRoutinePolicy = &inMemory{routinePolicy: make(map[int]cron.EntryID)}
var mu sync.Mutex

func runCronJobs(s *cron.Cron, rtp types.RoutineTransactionPolicy) {
	var frequency = "every 5s"
	entryId, err := s.AddFunc("@"+frequency, func() { fmt.Println("Every hour on the half hour") })
	if err != nil {
		log.Println(err)
	}
	mu.Lock()
	inMemoryRoutinePolicy.routinePolicy[rtp.ID] = entryId
	mu.Unlock()
	fmt.Println(inMemoryRoutinePolicy.routinePolicy)
	fmt.Println(rtp.RequestType)
	s.Start()
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	s := cron.New()
	var data types.RoutineTransactionPolicy
	config := publisher.NewNSQConfig()
	c, err := publisher.NewConsumer(context.Background(), "Add-Routine-Transaction-Policy", "Rtp", config)
	if err != nil {
		return
	}
	c.AddHandler(func(ctx context.Context, m *nsq.Message) error {
		err := json.Unmarshal(m.Body, &data)
		if err != nil {
			return fmt.Errorf("error %w", err)
		}
		switch data.RequestType {
		case "POST":
			runCronJobs(s, data)
		case "DELETE":
			if entryId, ok := inMemoryRoutinePolicy.routinePolicy[data.ID]; ok {
				s.Remove(entryId)
				delete(inMemoryRoutinePolicy.routinePolicy, data.ID)
				fmt.Println(inMemoryRoutinePolicy.routinePolicy)
			}
		case "PUT":
			if entryId, ok := inMemoryRoutinePolicy.routinePolicy[data.ID]; ok {
				s.Remove(entryId)
				delete(inMemoryRoutinePolicy.routinePolicy, data.ID)
				fmt.Println(inMemoryRoutinePolicy.routinePolicy)
				runCronJobs(s, data)
			}
		}
		log.Println(data)
		return nil
	})

	err = c.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		log.Panic("Could not connect")
	}
	wg.Wait()
	fmt.Scanln()
}

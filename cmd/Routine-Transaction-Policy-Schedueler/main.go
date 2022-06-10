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
	"github.com/robfig/cron"
)

func runCronJobs() {
	s := cron.New()

	var frequency = "every 5s"
	s.AddFunc("@"+frequency, func() { fmt.Println("Every hour on the half hour") })
	s.Start()
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

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
		log.Println(data)
		return nil
	})

	err = c.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		log.Panic("Could not connect")
	}
	wg.Wait()

	// runCronJobs()
	// fmt.Scanln()
}

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"time"

	"git.digitus.me/library/prosper-kit/application"
	"git.digitus.me/library/prosper-kit/config"
	prospercontext "git.digitus.me/library/prosper-kit/context"
	"git.digitus.me/pfe/smart-wallet/internal"
	service "git.digitus.me/pfe/smart-wallet/service/routine-transaction-ticker"
	"git.digitus.me/prosperus/publisher"
	"github.com/nsqio/go-nsq"
	"github.com/robfig/cron/v3"
	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/pgxv4"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type smartWalletConfig struct {
	DB struct {
		Host       string `json:"host"`
		Name       string `json:"name"`
		Pass       string `json:"pass"`
		Port       int    `json:"port"`
		User       string `json:"user"`
		SSLEnabled bool   `json:"sslEnabled"`
	} `json:"db"`
	Address          string `json:"address"`
	Port             int    `json:"port"`
	ConsulAddress    string `json:"consulAddress"`
	NsqLookupAddress string `json:"nsqLookUpAddress"`
}

func withCron(ctx context.Context, f func(*cron.Cron) error) error {
	endCh := make(chan struct{})
	c := cron.New()

	c.Start()

	go func() {
		<-ctx.Done()
		<-c.Stop().Done()
		close(endCh)
	}()

	defer func() { <-endCh }()

	if err := f(c); err != nil {
		return err
	}

	return nil
}

func run(ctx context.Context) (err error) {
	var logger *zap.Logger
	if logger, err = zap.NewDevelopment(); err != nil {
		return
	}

	ctx = prospercontext.WithLogger(ctx, logger)

	cfgPath := flag.String("config", "config file path", "")
	flag.Parse()

	_ = logger

	var cfg *smartWalletConfig
	{
		cfg, err = config.GetConfigFromVaultAndFile[smartWalletConfig](
			ctx, "smart-wallet", *cfgPath,
		)
		if err != nil {
			return err
		}
	}

	var db *sql.DB
	{
		sslMode := "disable"
		if cfg.DB.SSLEnabled {
			sslMode = "required"
		}

		connectionString := fmt.Sprintf(
			"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
			cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, sslMode,
		)

		db, err = apmsql.Open(
			"pgx",
			fmt.Sprintf("%s search_path=smart-wallet,wallet", connectionString),
		)
		if err != nil {
			return fmt.Errorf("could not open db connection: %w", err)
		}

		db.SetConnMaxIdleTime(time.Minute * 30)
		db.SetConnMaxLifetime(time.Hour)
		db.SetMaxIdleConns(25)
		db.SetMaxOpenConns(25)

		defer func() {
			err = multierr.Combine(err, db.Close())
		}()
	}

	config := publisher.NewNSQConfig()
	c, err := publisher.NewConsumer(
		context.Background(), internal.RoutineTransactionPolicyTopic, "Rtp", config,
	)
	if err != nil {
		return err
	}

	defer func() {
		c.Stop()
		<-c.StopChan
	}()

	p, err := nsq.NewProducer(cfg.NsqLookupAddress, config)
	if err != nil {
		return err
	}

	return withCron(ctx, func(s *cron.Cron) error {
		rtpTicker := service.NewRTPTicker(db, s, p)
		if err := rtpTicker.StartAllRoutinePolicies(ctx); err != nil {
			return err
		}

		// handler func should be a method on inMemory
		c.AddHandler(service.HandleRTPMessages(rtpTicker))

		// Should be read from config
		err = c.ConnectToNSQLookupd(cfg.NsqLookupAddress)
		if err != nil {
			return err
		}

		return nil
	})

}

func main() {
	application.WaitForInterrupt(run)
}

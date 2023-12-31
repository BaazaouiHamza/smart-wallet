package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.digitus.me/library/prosper-kit/application"
	"git.digitus.me/library/prosper-kit/authentication"
	"git.digitus.me/library/prosper-kit/config"
	prospercontext "git.digitus.me/library/prosper-kit/context"
	"git.digitus.me/library/prosper-kit/discovery"
	"git.digitus.me/pfe/smart-wallet/internal"
	service "git.digitus.me/pfe/smart-wallet/service/trigger-handler-service"
	"git.digitus.me/pfe/smart-wallet/types"

	ptclTypes "git.digitus.me/prosperus/protocol/types"
	"git.digitus.me/prosperus/publisher"
	"github.com/nsqio/go-nsq"
	"go.elastic.co/apm/module/apmhttp"
	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/pgxv4"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

func run(ctx context.Context) (err error) {
	var logger *zap.Logger
	if logger, err = zap.NewDevelopment(); err != nil {
		return
	}
	ctx = prospercontext.WithLogger(ctx, logger)

	cfgPath := flag.String("config", "config file path", "")
	flag.Parse()

	_ = logger

	var cfg *types.SmartWalletConfig
	{
		cfg, err = config.GetConfigFromVaultAndFile[types.SmartWalletConfig](
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
	p, err := nsq.NewProducer(cfg.NsqLookupAddress, config)
	if err != nil {
		return err
	}

	var client *http.Client
	{

		var opts []discovery.ConsulDiscoveryOption
		if cfg.ConsulAddress != "" {
			opts = append(opts, discovery.ConsulDiscoveryWithAddress(cfg.ConsulAddress))
		}

		var discoverer discovery.Discoverer
		discoverer, err = discovery.NewConsulDiscovery(opts...)
		if err != nil {
			return
		}

		discoverer, err = discovery.DiscoverWithPolling(ctx, discoverer)
		if err != nil {
			return
		}

		tc := &authentication.ClientCredentialsTokenCreator{
			Client: discovery.ClientWithDiscovery(apmhttp.WrapClient(&http.Client{
				Timeout: time.Second * 15,
			}), discoverer),
			ClientCredentials: authentication.ClientCredentials{
				ID:     cfg.ClientCredentials.ID,
				Secret: cfg.ClientCredentials.Secret,
			},
		}

		client = discovery.ClientWithDiscovery(
			authentication.ClientWithAuthentication(
				apmhttp.WrapClient(&http.Client{
					Timeout: time.Second * 15,
				}),
				tc,
			),
			discoverer,
		)
	}

	triggerHandler := service.NewTriggerHandler(db, p, &internal.CloudwalletClient{
		Client: client,
	})

	{
		triggerConsumer, err := publisher.NewConsumer(ctx, internal.TransactionsTopic, "smart-wallet-transaction", config)
		if err != nil {
			return err
		}
		defer func() {
			triggerConsumer.Stop()
			<-triggerConsumer.StopChan
		}()

		triggerConsumer.AddHandler(func(ctx context.Context, m *nsq.Message) error {
			var data types.TriggerMessage
			if err := publisher.Decode(m.Body, &data); err != nil {
				return err
			}

			return triggerHandler.HandleTrigger(ctx, data)
		})

		if err := triggerConsumer.ConnectToNSQD(cfg.NsqLookupAddress); err != nil {
			return err
		}
	}

	{
		notarizationConsumer, err := publisher.NewConsumer(
			ctx, "notarizations", "smart-wallet-trigger", config,
		)
		if err != nil {
			return err
		}
		defer func() {
			notarizationConsumer.Stop()
			<-notarizationConsumer.StopChan
		}()

		notarizationConsumer.AddHandler(func(ctx context.Context, m *nsq.Message) error {
			var n ptclTypes.Notarization
			if err := publisher.Decode(m.Body, &n); err != nil {
				return err
			}

			return triggerHandler.HandleNotarization(ctx, &n)
		})

		if err := notarizationConsumer.ConnectToNSQLookupd(cfg.NsqLookupAddressProsperus); err != nil {
			return err
		}
	}
	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)
	<-killSig
	return nil

}

func main() {
	application.WaitForInterrupt(run)
}

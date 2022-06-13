package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"time"

	"git.digitus.me/library/prosper-kit/application"
	"git.digitus.me/library/prosper-kit/authentication"
	"git.digitus.me/library/prosper-kit/config"
	prospercontext "git.digitus.me/library/prosper-kit/context"
	"git.digitus.me/library/prosper-kit/discovery"
	"git.digitus.me/library/prosper-kit/server"
	walletrepository "git.digitus.me/library/prosper-kit/wallet/repository"
	"git.digitus.me/pfe/smart-wallet/api"
	"git.digitus.me/pfe/smart-wallet/repository"
	"git.digitus.me/pfe/smart-wallet/service"
	"git.digitus.me/prosperus/publisher"
	"github.com/nsqio/go-nsq"
	"go.elastic.co/apm/module/apmhttp"
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
	Address       string `json:"address"`
	Port          int    `json:"port"`
	ConsulAddress string `json:"consulAddress"`
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
		cfg, err = config.GetConfigFromVaultAndFile[smartWalletConfig](ctx, "smart-wallet", *cfgPath)
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

		if _, err := db.Exec(`CREATE SCHEMA IF NOT EXISTS "smart-wallet";`); err != nil {
			return err
		}

		if _, err := db.Exec(`CREATE SCHEMA IF NOT EXISTS "wallet";`); err != nil {
			return err
		}

		err = repository.MigrateUp(fmt.Sprintf("%s search_path=smart-wallet", connectionString))
		if err != nil {
			return err
		}

		err = walletrepository.MigrateUp(fmt.Sprintf("%s search_path=wallet", connectionString))
		if err != nil {
			return err
		}
	}

	// TODO: initialize service

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

	client := discovery.ClientWithDiscovery(apmhttp.WrapClient(&http.Client{
		Timeout: time.Second * 15,
	}), discoverer)

	var jwsGetter *authentication.UAAJWSGetter
	jwsGetter, err = authentication.NewUAAJWSGetter(
		ctx, authentication.UAAJWSGetterWithClient(client),
	)
	if err != nil {
		return
	}
	config := publisher.NewNSQConfig()
	p, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		return
	}
	var svc = service.NewSmartWallet(db, p)
	engine := server.ReasonableEngine()
	api.NewServer(svc, engine, jwsGetter)

	consulHost := cfg.ConsulAddress
	if consulHost == "" {
		consulHost = "http://testing.prosperus.tech:8500"
	}

	return discovery.WithRegistration(
		ctx,
		discovery.Registration{
			ConsulHost:     consulHost,
			ServiceName:    "smart-wallet",
			ServiceAddress: cfg.Address,
			ServicePort:    cfg.Port,
		},
		&http.Client{Timeout: 15 * time.Second},
		func() error {
			return server.WithServer(
				ctx,
				server.SetAddress(cfg.Address, cfg.Port),
				server.SetHandler(engine),
			)
		},
	)
}

func main() {
	application.WaitForInterrupt(run)
}

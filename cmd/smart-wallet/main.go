package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"git.digitus.me/library/prosper-kit/application"
	"git.digitus.me/library/prosper-kit/authentication"
	"git.digitus.me/library/prosper-kit/config"
	"git.digitus.me/library/prosper-kit/discovery"
	"git.digitus.me/library/prosper-kit/server"
	"git.digitus.me/pfe/smart-wallet/api"
	"git.digitus.me/pfe/smart-wallet/repository"
	"git.digitus.me/pfe/smart-wallet/service"
	"go.elastic.co/apm/module/apmhttp"
	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/pq"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type smartWalletConfig struct {
	DB struct {
		Host string `json:"host"`
		Name string `json:"name"`
		Pass string `json:"pass"`
		Port int    `json:"port"`
		User string `json:"user"`
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

	_ = logger

	var cfg smartWalletConfig
	if err = config.GetConfigFromVault(ctx, "smart-wallet", &cfg); err != nil {
		if fErr := config.GetConfigFromFile(&cfg, "../../config.json"); fErr != nil {
			return multierr.Combine(
				fmt.Errorf("could not read config from vault: %w", err),
				fmt.Errorf("could not read config from file: %w", fErr),
			)
		}
	}

	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name,
	)

	if err = repository.MigrateUp(connectionString); err != nil {
		logger.Error(
			"no migration",
			zap.Error(fmt.Errorf("migration failed: %w", err)),
		)
	}

	var db *sql.DB
	db, err = apmsql.Open("postgres", connectionString)
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
	var svc = service.NewSmartWallet(db)
	engine := server.ReasonableEngine(logger)
	api.NewServer(svc, engine, jwsGetter)

	consulHost := cfg.ConsulAddress
	if consulHost == "" {
		consulHost = "http://localhost:8500"
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
		logger,
		func() error {
			return server.WithServer(
				ctx,
				logger,
				server.SetAddress(cfg.Address, cfg.Port),
				server.SetHandler(engine),
			)
		},
	)
}

func main() {
	application.WaitForInterrupt(run)
}

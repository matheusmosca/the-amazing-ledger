package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/instrumentators"
	"github.com/stone-co/the-amazing-ledger/app/domain/usecases"
	"github.com/stone-co/the-amazing-ledger/app/gateways/db/postgres"
	"github.com/stone-co/the-amazing-ledger/app/gateways/rpc"
	"github.com/stone-co/the-amazing-ledger/app/instrumentation/newrelic"
)

func main() {
	logger := log.With().
		Str("module", "main").
		Str("build_time", BuildTime).
		Str("build_commit", BuildGitCommit).
		Logger()

	logger.Info().Msg("starting ledger process...")

	cfg, err := app.LoadConfig()
	if err != nil {
		logger.Panic().Err(err).Msg("failed to load app configurations")
	}

	nr, err := newrelic.NewApp(cfg.NewRelic.AppName, cfg.NewRelic.LicenseKey, logrus.NewEntry(logrus.New()))
	if err != nil {
		logger.Panic().Err(err).Msg("failed to start new relic")
	}

	ledgerInstrumentator := instrumentators.NewLedgerInstrumentator(nr)

	conn, err := postgres.ConnectPool(cfg.Postgres.DSN(), zerolog.New(os.Stderr))
	if err != nil {
		logger.Panic().Err(err).Msg("failed to connect to database")
	}
	logger.Info().Msg("connected to postgres pool")
	defer conn.Close()

	logger.Info().Msg("running migrations")
	if err = postgres.RunMigrations(cfg.Postgres.URL()); err != nil {
		logger.Panic().Err(err).Msg("failed to run database migrations")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.RPCServer.Host, cfg.RPCServer.Port))
	if err != nil {
		logger.Panic().Err(err).Msg("failed to listen")
	}

	ledgerRepository := postgres.NewLedgerRepository(conn, ledgerInstrumentator)
	ledgerUseCase := usecases.NewLedgerUseCase(ledgerRepository, ledgerInstrumentator)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	rpcServer, gwServer, err := rpc.NewServer(ctx, ledgerUseCase, nr, cfg, BuildGitCommit, BuildTime)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to create servers")
	}
	logger.Info().Msg("created rpc and gateway servers")

	go func() {
		logger.Info().Msg("rpcServer listening")
		if err = rpcServer.Serve(lis); err != nil {
			logger.Panic().Err(err).Msg("failed to serve rpc server")
		}
	}()

	go func() {
		<-ctx.Done()
		logger.Info().Msg("context canceled, initiating graceful stop")

		ctx, cancel = context.WithTimeout(context.Background(), cfg.HttpServer.ShutdownTimeout)
		defer cancel()

		rpcServer.GracefulStop()
		logger.Info().Msg("rpcServer stopped")

		if err = gwServer.Shutdown(ctx); err != nil {
			_ = gwServer.Close()
			logger.Error().Err(err).Msg("failed to stop gateway server gracefully")
		}
		logger.Info().Msg("gateway stopped")
	}()

	go handleInterrupt(cancel)

	logger.Info().Msg("gatewayServer up")
	err = gwServer.ListenAndServe()
	if err != nil {
		logger.Panic().Err(err).Msg("failed to listen and serve gateway server")
	}
}

func handleInterrupt(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	sig := <-signals
	log.Info().Str("signal", sig.String()).Msg("captured signal - server shutdown")
	signal.Stop(signals)
	cancel()
}

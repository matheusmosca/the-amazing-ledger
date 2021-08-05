package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	logger := log.With().
		Str("module", "main").
		Str("build_time", BuildTime).
		Str("build_commit", BuildGitCommit).
		Logger()

	logger.Info().Msg("starting ledger process...")

	flag.Parse()
	if *cpuprofile != "" {
		logger.Printf("profilling cpu")
		f, err := os.Create(*cpuprofile)
		if err != nil {
			logger.Panic().Err(err).Msg("could not create CPU profile: ")
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			logger.Panic().Err(err).Msg("could not start CPU profile: ")
		}
		defer pprof.StopCPUProfile()
	}

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
	defer conn.Close()

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

	go func() {
		if err = rpcServer.Serve(lis); err != nil {
			logger.Panic().Err(err).Msg("failed to serve rpc server")
		}
	}()

	go func() {
		<-ctx.Done()

		ctx, cancel = context.WithTimeout(context.Background(), cfg.HttpServer.ShutdownTimeout)
		defer cancel()

		rpcServer.GracefulStop()

		if err = gwServer.Shutdown(ctx); err != nil {
			_ = gwServer.Close()
			logger.Error().Err(err).Msg("failed to stop gateway server gracefully")
		}
	}()

	go handleInterrupt(cancel)

	err = gwServer.ListenAndServe()
	if err != nil {
		logger.Panic().Err(err).Msg("failed to listen and serve gateway server")
	}

	if *memprofile != "" {

		logger.Printf("profilling memory")
		f, err := os.Create(*memprofile)
		if err != nil {
			logger.Panic().Err(err).Msg("could not create memory profile: ")
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			logger.Panic().Err(err).Msg("could not write memory profile: ")
		}
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

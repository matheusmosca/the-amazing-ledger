package testutils

import (
	"context"
	"fmt"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/instrumentators"
	"github.com/stone-co/the-amazing-ledger/app/domain/usecases"
	"github.com/stone-co/the-amazing-ledger/app/gateways/db/postgres"
	"github.com/stone-co/the-amazing-ledger/app/gateways/rpc"
	"github.com/stone-co/the-amazing-ledger/app/tests/testenv"
)

func StartServer(ctx context.Context, db *pgxpool.Pool, cfg *app.Config, startGatewayServer bool) {
	zerolog.SetGlobalLevel(zerolog.FatalLevel)

	nr, err := newrelic.NewApplication(newrelic.ConfigEnabled(false))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create newrelic application")
	}

	ledgerInstrumentator := instrumentators.NewLedgerInstrumentator(nr)
	ledgerRepository := postgres.NewLedgerRepository(db, ledgerInstrumentator)
	ledgerUsecase := usecases.NewLedgerUseCase(ledgerRepository, ledgerInstrumentator)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.RPCServer.Host, cfg.RPCServer.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	buildCommit := "undefined"
	buildTime := "undefined"

	rpcServer, gwServer, err := rpc.NewServer(ctx, ledgerUsecase, nr, cfg, buildCommit, buildTime)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create servers")
	}

	go func() {
		if err = rpcServer.Serve(listener); err != nil {
			log.Fatal().Err(err).Msg("failed to start rpc server")
		}
	}()

	if startGatewayServer {
		go func() {
			_ = gwServer.ListenAndServe()
		}()
	}

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), cfg.HttpServer.ShutdownTimeout)
		defer cancel()

		rpcServer.GracefulStop()

		if startGatewayServer {
			if err = gwServer.Shutdown(ctx); err != nil {
				_ = gwServer.Close()
				log.Fatal().Err(err).Msg("failed to stop the server gracefully")
			}
		}
	}()

	testenv.LedgerRepository = ledgerRepository
	testenv.GatewayServer = fmt.Sprintf("http://%s", gwServer.Addr)
}

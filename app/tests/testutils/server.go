package testutils

import (
	"context"
	"fmt"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/usecases"
	"github.com/stone-co/the-amazing-ledger/app/gateways/db/postgres"
	"github.com/stone-co/the-amazing-ledger/app/gateways/rpc"
	"github.com/stone-co/the-amazing-ledger/app/tests/testenv"
)

func StartServer(ctx context.Context, db *pgxpool.Pool, cfg *app.Config, startGatewayServer bool) {
	log := logrus.New()

	nr, err := newrelic.NewApplication(newrelic.ConfigEnabled(false))
	if err != nil {
		log.Fatalf("failed to create newrelic application: %v", err)
	}

	ledgerRepository := postgres.NewLedgerRepository(db, log)
	ledgerUsecase := usecases.NewLedgerUseCase(log, ledgerRepository)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.RPCServer.Host, cfg.RPCServer.Port))
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}

	buildCommit := "undefined"
	buildTime := "undefined"

	rpcServer, gwServer, err := rpc.NewServer(ctx, ledgerUsecase, nr, cfg, log, buildCommit, buildTime)
	if err != nil {
		log.Panicf("could not create server: %v", err)
	}

	go func() {
		if err = rpcServer.Serve(listener); err != nil {
			log.Panicf("could not start rpc server: %v", err)
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
				log.WithError(err).Fatal("could not stop server gracefully")
			}
		}
	}()

	testenv.LedgerRepository = ledgerRepository
	testenv.GatewayServer = fmt.Sprintf("http://%s", gwServer.Addr)
}

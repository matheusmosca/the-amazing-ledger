package testutils

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/usecases"
	"github.com/stone-co/the-amazing-ledger/app/gateways/db/postgres"
	"github.com/stone-co/the-amazing-ledger/app/gateways/rpc"
	"github.com/stone-co/the-amazing-ledger/app/tests/testenv"
)

func StartServer(db *pgxpool.Pool, listener *bufconn.Listener, startGatewayServer bool) (*grpc.Server, *http.Server) {
	log := logrus.New()

	nr, err := newrelic.NewApplication(newrelic.ConfigEnabled(false))
	if err != nil {
		log.Fatalf("failed to create newrelic application: %v", err)
	}

	ledgerRepository := postgres.NewLedgerRepository(db, log)
	ledgerUsecase := usecases.NewLedgerUseCase(log, ledgerRepository)

	cfg := &app.Config{
		RPCServer: app.RPCServerConfig{
			Host:            "0.0.0.0",
			Port:            3000,
			ShutdownTimeout: 5 * time.Second,
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    10 * time.Second,
		},
		HttpServer: app.HttpServerConfig{
			Host:            "0.0.0.0",
			Port:            3001,
			ShutdownTimeout: 5 * time.Second,
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    10 * time.Second,
		},
	}

	ctx := context.Background()
	commit := "undefined"
	time := "undefined"

	rpcServer, gwServer, err := rpc.NewServer(ctx, ledgerUsecase, nr, cfg, log, commit, time)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := rpcServer.Serve(listener); err != nil {
			panic(err)
		}
	}()

	if startGatewayServer {
		go func() {
			if err := gwServer.ListenAndServe(); err != nil {
				panic(err)
			}
		}()
	}

	testenv.LedgerRepository = ledgerRepository

	return rpcServer, gwServer
}

func GetBufDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
}

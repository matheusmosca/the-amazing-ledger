package rpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/tests"
	"github.com/stone-co/the-amazing-ledger/app/tests/testenv"
	"github.com/stone-co/the-amazing-ledger/app/tests/testutils"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func TestMain(m *testing.M) {
	pgDocker := tests.SetupTest("../../../gateways/db/postgres/migrations")

	ctx, cancel := context.WithCancel(context.Background())
	cfg := &app.Config{
		RPCServer: app.RPCServerConfig{
			Host:            "0.0.0.0",
			Port:            6000,
			ShutdownTimeout: 5 * time.Second,
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    10 * time.Second,
		},
		HttpServer: app.HttpServerConfig{
			Host:            "0.0.0.0",
			Port:            6001,
			ShutdownTimeout: 5 * time.Second,
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    10 * time.Second,
		},
	}

	testutils.StartServer(ctx, pgDocker.DB, cfg, false)

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", cfg.RPCServer.Host, cfg.RPCServer.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not create client connection: %v", err)
	}

	testenv.RPCClient = proto.NewLedgerServiceClient(conn)
	testenv.DB = pgDocker.DB

	_, err = pgDocker.DB.Exec(context.Background(), `insert into event (id, name) values (1, 'default');`)
	if err != nil {
		log.Fatalf("could not insert default event values: %v", err)
	}

	exitCode := m.Run()

	cancel()
	tests.RemoveContainer(pgDocker)
	os.Exit(exitCode)
}

package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/tests"
	"github.com/stone-co/the-amazing-ledger/app/tests/testenv"
	"github.com/stone-co/the-amazing-ledger/app/tests/testutils"
)

func TestMain(m *testing.M) {
	pgDocker := tests.SetupTest("../../../gateways/db/postgres/migrations")

	testenv.DB = pgDocker.DB
	ctx, cancel := context.WithCancel(context.Background())
	cfg := &app.Config{
		RPCServer: app.RPCServerConfig{
			Host:            "0.0.0.0",
			Port:            5000,
			ShutdownTimeout: 5 * time.Second,
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    10 * time.Second,
		},
		HttpServer: app.HttpServerConfig{
			Host:            "0.0.0.0",
			Port:            5001,
			ShutdownTimeout: 5 * time.Second,
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    10 * time.Second,
		},
	}

	testutils.StartServer(ctx, pgDocker.DB, cfg, true)

	_, err := pgDocker.DB.Exec(context.Background(), `insert into event (id, name) values (1, 'default');`)
	if err != nil {
		log.Fatalf("could not insert default event values: %v", err)
	}

	ch := make(chan int, 1)

	go func() {
		if err := waitForGateway(context.Background(), cfg.HttpServer.Port); err != nil {
			log.Printf("waitForGateway failed with %v", err)
		}

		ch <- m.Run()
	}()

	exitCode := <-ch
	cancel()
	tests.RemoveContainer(pgDocker)

	os.Exit(exitCode)
}

func waitForGateway(ctx context.Context, port int) error {
	ch := time.After(10 * time.Second)

	for {
		r, err := http.Get(fmt.Sprintf("http://localhost:%d/health", port))
		if err != nil {
			defer r.Body.Close()
		}

		if err == nil && r.StatusCode == http.StatusOK {
			return nil
		}

		log.Printf("Waiting for localhost:%d to get ready", port)

		select {
		case <-ctx.Done():
			return err
		case <-ch:
			return err
		case <-time.After(10 * time.Millisecond):
		}
	}
}

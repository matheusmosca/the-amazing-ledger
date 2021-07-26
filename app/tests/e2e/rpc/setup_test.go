package rpc

import (
	"context"
	"log"
	"os"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/stone-co/the-amazing-ledger/app/tests"
	"github.com/stone-co/the-amazing-ledger/app/tests/testenv"
	"github.com/stone-co/the-amazing-ledger/app/tests/testutils"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func TestMain(m *testing.M) {
	pgDocker := tests.SetupTest("../../../gateways/db/postgres/migrations")

	listener := bufconn.Listen(1024 * 1024)
	rpcServer, _ := testutils.StartServer(pgDocker.DB, listener, false)

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(testutils.GetBufDialer(listener)), grpc.WithInsecure())
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

	tests.RemoveContainer(pgDocker)
	rpcServer.Stop()

	os.Exit(exitCode)
}

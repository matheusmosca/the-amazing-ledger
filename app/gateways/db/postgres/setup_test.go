package postgres

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests"
)

var pgDocker *tests.PostgresDocker

func TestMain(m *testing.M) {
	pgDocker = tests.SetupTest("./migrations")

	_, err := pgDocker.DB.Exec(context.Background(), `insert into event (id, name) values (1, 'default'), (2, 'new');`)
	if err != nil {
		log.Fatalf("could not insert default event values: %v", err)
	}

	exitCode := m.Run()

	tests.RemoveContainer(pgDocker)

	os.Exit(exitCode)
}

func createEntry(t *testing.T, op vos.OperationType, account string, version vos.Version, amount int) entities.Entry {
	t.Helper()

	entry, err := entities.NewEntry(
		uuid.New(),
		op,
		account,
		version,
		amount,
		json.RawMessage(`{}`),
	)
	assert.NoError(t, err)

	return entry
}

func createTransaction(t *testing.T, ctx context.Context, r *LedgerRepository, entries ...entities.Entry) entities.Transaction {
	t.Helper()

	tx, err := entities.NewTransaction(
		uuid.New(),
		uint32(1),
		"abc",
		time.Now().Round(time.Microsecond),
		entries...,
	)
	assert.NoError(t, err)

	err = r.CreateTransaction(ctx, tx)
	assert.NoError(t, err)

	return tx
}

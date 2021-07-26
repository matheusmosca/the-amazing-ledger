package testseed

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/tests/testenv"
)

func CreateTransaction(t *testing.T, entries ...entities.Entry) entities.Transaction {
	t.Helper()

	tx, err := entities.NewTransaction(uuid.New(), uint32(1), "abc", time.Now(), entries...)
	assert.NoError(t, err)

	err = testenv.LedgerRepository.CreateTransaction(context.Background(), tx)
	assert.NoError(t, err)

	return tx
}

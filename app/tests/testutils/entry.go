package testutils

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func CreateEntry(t *testing.T, op vos.OperationType, account string, version vos.Version, amount int) entities.Entry {
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

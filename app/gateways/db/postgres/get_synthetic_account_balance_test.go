package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/instrumentators"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests"
)

func TestLedgerRepository_QueryAggregatedBalanceFailure(t *testing.T) {
	t.Run("should return an error if accounts do not exist", func(t *testing.T) {
		r := NewLedgerRepository(pgDocker.DB, &instrumentators.LedgerInstrumentator{})
		ctx := context.Background()

		query, err := vos.NewAccount("liability.agg.*")
		assert.NoError(t, err)

		_, err = r.GetSyntheticAccountBalance(ctx, query)
		assert.ErrorIs(t, err, app.ErrAccountNotFound)
	})
}

func TestLedgerRepository_QueryAggregatedBalanceSuccess(t *testing.T) {
	acc1, err := vos.NewAccount("liability.agg.agg1")
	assert.NoError(t, err)

	acc2, err := vos.NewAccount("liability.agg.agg2")
	assert.NoError(t, err)

	acc3, err := vos.NewAccount("liability.abc.agg3")
	assert.NoError(t, err)

	query, err := vos.NewAccount("liability.agg.*")
	assert.NoError(t, err)

	type wants struct {
		accountBalance int
		snapBalance    int
		snapErr        error
	}

	testCases := []struct {
		name     string
		repoSeed func(t *testing.T, ctx context.Context, r *LedgerRepository)
		wants    wants
	}{
		{
			name: "should query aggregated balance involving two accounts",
			repoSeed: func(t *testing.T, ctx context.Context, r *LedgerRepository) {
				e1 := createEntry(t, vos.DebitOperation, acc1.Value(), vos.NextAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, acc2.Value(), vos.IgnoreAccountVersion, 100)
				createTransaction(t, ctx, r, e1, e2)
			},
			wants: wants{
				accountBalance: 0,
				snapErr:        pgx.ErrNoRows,
			},
		},
		{
			name: "should query aggregated balance involving three accounts (first snapshot)",
			repoSeed: func(t *testing.T, ctx context.Context, r *LedgerRepository) {
				e1 := createEntry(t, vos.DebitOperation, acc1.Value(), vos.NextAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, acc2.Value(), vos.IgnoreAccountVersion, 100)
				createTransaction(t, ctx, r, e1, e2)

				e1 = createEntry(t, vos.CreditOperation, acc1.Value(), vos.NextAccountVersion, 100)
				e2 = createEntry(t, vos.CreditOperation, acc2.Value(), vos.NextAccountVersion, 100)
				e3 := createEntry(t, vos.DebitOperation, acc3.Value(), vos.IgnoreAccountVersion, 200)
				createTransaction(t, ctx, r, e1, e2, e3)
			},
			wants: wants{
				accountBalance: 200,
				snapBalance:    0,
			},
		},
		{
			name: "should query aggregated balance involving three accounts (second snapshot)",
			repoSeed: func(t *testing.T, ctx context.Context, r *LedgerRepository) {
				e1 := createEntry(t, vos.DebitOperation, acc1.Value(), vos.NextAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, acc2.Value(), vos.IgnoreAccountVersion, 100)
				createTransaction(t, ctx, r, e1, e2)

				e1 = createEntry(t, vos.CreditOperation, acc1.Value(), vos.NextAccountVersion, 100)
				e2 = createEntry(t, vos.CreditOperation, acc2.Value(), vos.NextAccountVersion, 100)
				e3 := createEntry(t, vos.DebitOperation, acc3.Value(), vos.IgnoreAccountVersion, 200)
				createTransaction(t, ctx, r, e1, e2, e3)

				e1 = createEntry(t, vos.DebitOperation, acc1.Value(), vos.IgnoreAccountVersion, 100)
				e3 = createEntry(t, vos.CreditOperation, acc3.Value(), vos.NextAccountVersion, 100)
				createTransaction(t, ctx, r, e1, e3)

				_, err = r.GetSyntheticAccountBalance(ctx, query)
				assert.NoError(t, err)
			},
			wants: wants{
				accountBalance: 100,
				snapBalance:    200,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			r := NewLedgerRepository(pgDocker.DB, &instrumentators.LedgerInstrumentator{})
			tt.repoSeed(t, ctx, r)

			defer tests.TruncateTables(ctx, pgDocker.DB, "entry", "account_version", "account_balance")

			balance, err := r.GetSyntheticAccountBalance(ctx, query)
			assert.NoError(t, err)
			assert.Equal(t, tt.wants.accountBalance, balance.Balance)

			if tt.wants.snapErr != nil {
				_, err = fetchQuerySnapshot(ctx, pgDocker.DB, query)
				assert.ErrorIs(t, err, pgx.ErrNoRows)
			} else {
				snap, err := fetchQuerySnapshot(ctx, pgDocker.DB, query)
				assert.NoError(t, err)
				assert.Equal(t, tt.wants.snapBalance, snap.balance)
			}
		})
	}
}

type querySnapshot struct {
	balance int
	date    time.Time
}

func fetchQuerySnapshot(ctx context.Context, db *pgxpool.Pool, query vos.Account) (querySnapshot, error) {
	const cmd = "select balance, tx_date from account_balance where account = $1;"

	var snap querySnapshot

	err := db.QueryRow(ctx, cmd, query.Value()).Scan(&snap.balance, &snap.date)
	if err != nil {
		return querySnapshot{}, fmt.Errorf("failed to fetch query snapshot: %w", err)
	}

	return snap, nil
}

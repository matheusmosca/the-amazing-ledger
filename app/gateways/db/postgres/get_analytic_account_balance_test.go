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
	"github.com/stone-co/the-amazing-ledger/app/tests/testdata"
)

func TestLedgerRepository_GetAccountBalanceSuccess(t *testing.T) {
	acc1, err := vos.NewAnalyticAccount(testdata.GenerateAccountPath())
	assert.NoError(t, err)

	acc2, err := vos.NewAnalyticAccount(testdata.GenerateAccountPath())
	assert.NoError(t, err)

	type accountValues struct {
		acc1Balance int
		acc2balance int
	}

	type wants struct {
		total    accountValues
		snapshot accountValues
		snapErr  error
	}

	testCases := []struct {
		name     string
		repoSeed func(t *testing.T, ctx context.Context, r *LedgerRepository)
		wants    wants
	}{
		{
			name: "should get account balance successfully when is the first request",
			repoSeed: func(t *testing.T, ctx context.Context, r *LedgerRepository) {
				e1 := createEntry(t, vos.DebitOperation, acc1.Value(), vos.NextAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, acc2.Value(), vos.NextAccountVersion, 100)

				createTransaction(t, ctx, r, e1, e2)
			},
			wants: wants{
				total: accountValues{
					acc1Balance: -100,
					acc2balance: 100,
				},
				snapErr: pgx.ErrNoRows,
			},
		},
		{
			name: "should get account balance successfully when is the second request",
			repoSeed: func(t *testing.T, ctx context.Context, r *LedgerRepository) {
				e1 := createEntry(t, vos.DebitOperation, acc1.Value(), vos.NextAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, acc2.Value(), vos.NextAccountVersion, 100)
				createTransaction(t, ctx, r, e1, e2)

				e1 = createEntry(t, vos.DebitOperation, acc1.Value(), vos.NextAccountVersion, 100)
				e2 = createEntry(t, vos.CreditOperation, acc2.Value(), vos.NextAccountVersion, 100)
				createTransaction(t, ctx, r, e1, e2)
			},
			wants: wants{
				total: accountValues{
					acc1Balance: -200,
					acc2balance: 200,
				},
				snapshot: accountValues{
					acc1Balance: -100,
					acc2balance: 100,
				},
			},
		},
		{
			name: "should get account balance successfully when is the third request",
			repoSeed: func(t *testing.T, ctx context.Context, r *LedgerRepository) {
				e1 := createEntry(t, vos.DebitOperation, acc1.Value(), vos.NextAccountVersion, 100)
				e2 := createEntry(t, vos.CreditOperation, acc2.Value(), vos.NextAccountVersion, 100)
				createTransaction(t, ctx, r, e1, e2)

				e1 = createEntry(t, vos.DebitOperation, acc1.Value(), vos.NextAccountVersion, 100)
				e2 = createEntry(t, vos.CreditOperation, acc2.Value(), vos.NextAccountVersion, 100)
				createTransaction(t, ctx, r, e1, e2)

				e1 = createEntry(t, vos.CreditOperation, acc1.Value(), vos.NextAccountVersion, 100)
				e2 = createEntry(t, vos.DebitOperation, acc2.Value(), vos.NextAccountVersion, 100)
				createTransaction(t, ctx, r, e1, e2)
			},
			wants: wants{
				total: accountValues{
					acc1Balance: -100,
					acc2balance: 100,
				},
				snapshot: accountValues{
					acc1Balance: -200,
					acc2balance: 200,
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			r := NewLedgerRepository(pgDocker.DB, &instrumentators.LedgerInstrumentator{})

			tt.repoSeed(t, ctx, r)

			defer tests.TruncateTables(ctx, pgDocker.DB, "entry", "account_version", "account_balance")

			balance, err := r.GetAnalyticAccountBalance(ctx, acc1)
			assert.NoError(t, err)
			assert.Equal(t, tt.wants.total.acc1Balance, balance.Balance)

			balance, err = r.GetAnalyticAccountBalance(ctx, acc2)
			assert.NoError(t, err)
			assert.Equal(t, tt.wants.total.acc2balance, balance.Balance)

			if tt.wants.snapErr != nil {
				_, err = fetchSnapshot(ctx, pgDocker.DB, acc1)
				assert.ErrorIs(t, err, tt.wants.snapErr)

				_, err = fetchSnapshot(ctx, pgDocker.DB, acc2)
				assert.ErrorIs(t, err, tt.wants.snapErr)
			} else {
				snap, err := fetchSnapshot(ctx, pgDocker.DB, acc1)
				assert.NoError(t, err)
				assert.Equal(t, tt.wants.snapshot.acc1Balance, snap.balance)

				snap, err = fetchSnapshot(ctx, pgDocker.DB, acc2)
				assert.NoError(t, err)
				assert.Equal(t, tt.wants.snapshot.acc2balance, snap.balance)
			}
		})
	}
}

func TestLedgerRepository_GetAccountBalanceFailure(t *testing.T) {
	t.Run("should return an error if account does not exist", func(t *testing.T) {
		r := NewLedgerRepository(pgDocker.DB, &instrumentators.LedgerInstrumentator{})

		acc, err := vos.NewAnalyticAccount(testdata.GenerateAccountPath())
		assert.NoError(t, err)

		_, err = r.GetAnalyticAccountBalance(context.Background(), acc)
		assert.ErrorIs(t, app.ErrAccountNotFound, err)
	})
}

type snapshot struct {
	balance int
	date    time.Time
}

func fetchSnapshot(ctx context.Context, db *pgxpool.Pool, account vos.Account) (snapshot, error) {
	const query = "select balance, tx_date from account_balance where account = $1;"

	var snap snapshot

	err := db.QueryRow(ctx, query, account.Value()).Scan(
		&snap.balance,
		&snap.date,
	)
	if err != nil {
		return snapshot{}, fmt.Errorf("failed to fetch snapshot: %w", err)
	}

	return snap, nil
}

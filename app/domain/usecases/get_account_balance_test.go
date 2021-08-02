package usecases

import (
	"context"
	"testing"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/instrumentators"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests/mocks"
	"github.com/stone-co/the-amazing-ledger/app/tests/testdata"
)

func TestLedgerUseCase_GetAccountBalance_Analytic(t *testing.T) {
	t.Run("should return account balance successfully", func(t *testing.T) {
		accountPath, err := vos.NewAnalyticAccount(testdata.GenerateAccountPath())
		assert.NoError(t, err)

		accountBalance := vos.NewAnalyticAccountBalance(accountPath, vos.Version(1), 150)
		mockedRepository := &mocks.RepositoryMock{
			GetAnalyticAccountBalanceFunc: func(ctx context.Context, account vos.Account) (vos.AccountBalance, error) {
				return accountBalance, nil
			},
		}
		usecase := NewLedgerUseCase(mockedRepository, instrumentators.NewLedgerInstrumentator(logrus.New(), &newrelic.Application{}))

		got, err := usecase.GetAccountBalance(context.Background(), accountPath)
		assert.NoError(t, err)

		assert.Equal(t, accountBalance.Account, got.Account)
		assert.Equal(t, accountBalance.CurrentVersion, got.CurrentVersion)
		assert.Equal(t, accountBalance.Balance, got.Balance)
	})

	t.Run("should return an error if account does not exist", func(t *testing.T) {
		accountPath, err := vos.NewAnalyticAccount(testdata.GenerateAccountPath())
		assert.NoError(t, err)

		mockedRepository := &mocks.RepositoryMock{
			GetAnalyticAccountBalanceFunc: func(ctx context.Context, account vos.Account) (vos.AccountBalance, error) {
				return vos.AccountBalance{}, app.ErrAccountNotFound
			},
		}
		usecase := NewLedgerUseCase(mockedRepository, instrumentators.NewLedgerInstrumentator(logrus.New(), &newrelic.Application{}))

		got, err := usecase.GetAccountBalance(context.Background(), accountPath)
		assert.Empty(t, got)
		assert.ErrorIs(t, err, app.ErrAccountNotFound)
	})
}

func TestLedgerUseCase_GetAccountBalance_Synthetic(t *testing.T) {
	t.Run("should return aggregated balance successfully", func(t *testing.T) {
		account, err := vos.NewAccount("liability.stone.clients.*")
		assert.NoError(t, err)

		queryBalance := vos.NewSyntheticAccountBalance(account, 20)
		mockedRepository := &mocks.RepositoryMock{
			GetSyntheticAccountBalanceFunc: func(ctx context.Context, account vos.Account) (vos.AccountBalance, error) {
				return queryBalance, nil
			},
		}

		nr, _ := newrelic.NewApplication()
		usecase := NewLedgerUseCase(mockedRepository, instrumentators.NewLedgerInstrumentator(logrus.New(), nr))

		got, err := usecase.GetAccountBalance(context.Background(), account)
		assert.NoError(t, err)
		assert.Equal(t, queryBalance.Balance, got.Balance)
	})

	t.Run("should return an error if account does not exist", func(t *testing.T) {
		query, err := vos.NewAccount("liability.stone.clients.*")
		assert.NoError(t, err)

		mockedRepository := &mocks.RepositoryMock{
			GetSyntheticAccountBalanceFunc: func(ctx context.Context, account vos.Account) (vos.AccountBalance, error) {
				return vos.AccountBalance{}, app.ErrAccountNotFound
			},
		}

		nr, _ := newrelic.NewApplication()
		usecase := NewLedgerUseCase(mockedRepository, instrumentators.NewLedgerInstrumentator(logrus.New(), nr))

		got, err := usecase.GetAccountBalance(context.Background(), query)
		assert.Empty(t, got)
		assert.ErrorIs(t, err, app.ErrAccountNotFound)
	})
}

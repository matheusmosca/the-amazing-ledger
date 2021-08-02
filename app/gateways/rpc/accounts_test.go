package rpc

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests/mocks"
	"github.com/stone-co/the-amazing-ledger/app/tests/testdata"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func TestAPI_GetAccountBalance_Analytical_Success(t *testing.T) {
	t.Run("should get account balance successfully", func(t *testing.T) {
		account, err := vos.NewAccount(testdata.GenerateAccountPath())
		assert.NoError(t, err)

		accountBalance := vos.NewAnalyticalAccountBalance(account, vos.Version(1), 200, 100)
		mockedUsecase := &mocks.UseCaseMock{
			GetAccountBalanceFunc: func(ctx context.Context, accountPath vos.Account) (vos.AccountBalance, error) {
				accountBalance.Account = accountPath

				return accountBalance, nil
			},
		}
		api := NewAPI(logrus.New(), mockedUsecase)

		request := &proto.GetAccountBalanceRequest{
			Account: account.Value(),
		}

		got, err := api.GetAccountBalance(context.Background(), request)
		assert.NoError(t, err)

		assert.Equal(t, &proto.GetAccountBalanceResponse{
			Account:        request.Account,
			CurrentVersion: accountBalance.CurrentVersion.AsInt64(),
			Balance:        int64(accountBalance.TotalCredit - accountBalance.TotalDebit),
		}, got)
	})
}

func TestAPI_GetAccountBalance_Synthetic_Success(t *testing.T) {
	t.Run("should get aggregated balance successfully", func(t *testing.T) {
		account, err := vos.NewAccount("liability.stone.clients.*")
		assert.NoError(t, err)

		balance := vos.NewSyntheticAccountBalance(account, 100)
		mockedUsecase := &mocks.UseCaseMock{
			GetAccountBalanceFunc: func(ctx context.Context, account vos.Account) (vos.AccountBalance, error) {
				return balance, nil
			},
		}
		api := NewAPI(logrus.New(), mockedUsecase)

		request := &proto.GetAccountBalanceRequest{
			Account: "liability.stone.clients.*",
		}

		got, err := api.GetAccountBalance(context.Background(), request)
		assert.NoError(t, err)

		assert.Equal(t, &proto.GetAccountBalanceResponse{
			Account:        account.Value(),
			CurrentVersion: -1,
			Balance:        100,
		}, got)
	})
}

func TestAPI_GetAccountBalance_InvalidRequest(t *testing.T) {
	testCases := []struct {
		name            string
		useCaseSetup    *mocks.UseCaseMock
		request         *proto.GetAccountBalanceRequest
		expectedCode    codes.Code
		expectedMessage string
	}{
		{
			name:         "should return an error if account name is invalid",
			useCaseSetup: &mocks.UseCaseMock{},
			request: &proto.GetAccountBalanceRequest{
				Account: "liability.clients.abc-123.*",
			},
			expectedCode:    codes.InvalidArgument,
			expectedMessage: app.ErrInvalidAccountComponentCharacters.Error(),
		},
		{
			name: "should return an error if account does not exist",
			useCaseSetup: &mocks.UseCaseMock{
				GetAccountBalanceFunc: func(ctx context.Context, account vos.Account) (vos.AccountBalance, error) {
					return vos.AccountBalance{}, app.ErrAccountNotFound
				},
			},
			request: &proto.GetAccountBalanceRequest{
				Account: testdata.GenerateAccountPath(),
			},
			expectedCode:    codes.NotFound,
			expectedMessage: "account not found",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			api := NewAPI(logrus.New(), tt.useCaseSetup)

			_, err := api.GetAccountBalance(context.Background(), tt.request)
			respStatus, ok := status.FromError(err)

			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, respStatus.Code())
			assert.Equal(t, tt.expectedMessage, respStatus.Message())
		})
	}
}

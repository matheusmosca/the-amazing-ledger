package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests"
	"github.com/stone-co/the-amazing-ledger/app/tests/testdata"
	"github.com/stone-co/the-amazing-ledger/app/tests/testenv"
	"github.com/stone-co/the-amazing-ledger/app/tests/testseed"
	"github.com/stone-co/the-amazing-ledger/app/tests/testutils"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func TestE2E_RPC_CreateTransactionSuccess(t *testing.T) {
	t.Run("should create a transaction successfully", func(t *testing.T) {
		request := &proto.CreateTransactionRequest{
			Id: uuid.New().String(),
			Entries: []*proto.Entry{
				{
					Id:              uuid.New().String(),
					AccountId:       testdata.GenerateAccountPath(),
					ExpectedVersion: 3,
					Operation:       proto.Operation_OPERATION_DEBIT,
					Amount:          123,
				},
				{
					Id:              uuid.New().String(),
					AccountId:       testdata.GenerateAccountPath(),
					ExpectedVersion: 3,
					Operation:       proto.Operation_OPERATION_CREDIT,
					Amount:          123,
				},
			},
			Company:        "abc",
			Event:          1,
			CompetenceDate: timestamppb.Now(),
		}

		_, err := testenv.RPCClient.CreateTransaction(context.Background(), request)
		assert.NoError(t, err)
	})
}

func TestE2E_RPC_CreateTransactionFailure(t *testing.T) {
	testCases := []struct {
		name         string
		seedRepo     func(t *testing.T) entities.Transaction
		request      *proto.CreateTransactionRequest
		expectedCode codes.Code
		expectedMsg  string
	}{
		{
			name: "should return an error if id is invalid",
			request: &proto.CreateTransactionRequest{
				Id: "invalid UUID",
				Entries: []*proto.Entry{
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_DEBIT,
						Amount:          123,
					},
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_CREDIT,
						Amount:          123,
					},
				},
				Company:        "abc",
				Event:          1,
				CompetenceDate: timestamppb.Now(),
			},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "error parsing transaction id",
		},
		{
			name: "should return an error if entry id is invalid",
			request: &proto.CreateTransactionRequest{
				Id: uuid.New().String(),
				Entries: []*proto.Entry{
					{
						Id:              "invalid-entry-id",
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_DEBIT,
						Amount:          123,
					},
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_CREDIT,
						Amount:          123,
					},
				},
				Company:        "abc",
				Event:          1,
				CompetenceDate: timestamppb.Now(),
			},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "error parsing entry id",
		},
		{
			name: "should return an error if operation is invalid",
			request: &proto.CreateTransactionRequest{
				Id: uuid.New().String(),
				Entries: []*proto.Entry{
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_UNSPECIFIED,
						Amount:          123,
					},
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_CREDIT,
						Amount:          123,
					},
				},
				Company:        "abc",
				Event:          1,
				CompetenceDate: timestamppb.Now(),
			},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "invalid operation",
		},
		{
			name: "should return an error if amount is invalid",
			request: &proto.CreateTransactionRequest{
				Id: uuid.New().String(),
				Entries: []*proto.Entry{
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_DEBIT,
						Amount:          -123,
					},
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_CREDIT,
						Amount:          123,
					},
				},
				Company:        "abc",
				Event:          1,
				CompetenceDate: timestamppb.Now(),
			},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "invalid amount",
		},
		{
			name: "should return an error if number of entries is less than two",
			request: &proto.CreateTransactionRequest{
				Id: uuid.New().String(),
				Entries: []*proto.Entry{
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_CREDIT,
						Amount:          123,
					},
				},
				Company:        "abc",
				Event:          1,
				CompetenceDate: timestamppb.Now(),
			},
			expectedCode: codes.Aborted,
			expectedMsg:  "invalid entries number",
		},
		{
			name: "should return an error if account is invalid",
			request: &proto.CreateTransactionRequest{
				Id: uuid.New().String(),
				Entries: []*proto.Entry{
					{
						Id:              uuid.New().String(),
						AccountId:       "assets",
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_DEBIT,
						Amount:          123,
					},
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_CREDIT,
						Amount:          123,
					},
				},
				Company:        "abc",
				Event:          1,
				CompetenceDate: timestamppb.Now(),
			},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "invalid depth value",
		},
		{
			name: "should return if competence date is in the future",
			request: &proto.CreateTransactionRequest{
				Id: uuid.New().String(),
				Entries: []*proto.Entry{
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_DEBIT,
						Amount:          123,
					},
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_CREDIT,
						Amount:          123,
					},
				},
				Company:        "abc",
				Event:          1,
				CompetenceDate: timestamppb.New(time.Now().UTC().Add(1 * time.Minute)),
			},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "competence date set to the future",
		},
		{
			name: "should return an error when occurs idempotency key violation",
			seedRepo: func(t *testing.T) entities.Transaction {
				e1 := testutils.CreateEntry(t, vos.DebitOperation, testdata.GenerateAccountPath(), vos.NextAccountVersion, 100)
				e2 := testutils.CreateEntry(t, vos.CreditOperation, testdata.GenerateAccountPath(), vos.NextAccountVersion, 100)

				tx := testseed.CreateTransaction(t, e1, e2)
				return tx
			},
			request: &proto.CreateTransactionRequest{
				Id: uuid.New().String(),
				Entries: []*proto.Entry{
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_DEBIT,
						Amount:          123,
					},
					{
						Id:              uuid.New().String(),
						AccountId:       testdata.GenerateAccountPath(),
						ExpectedVersion: 3,
						Operation:       proto.Operation_OPERATION_CREDIT,
						Amount:          123,
					},
				},
				Company:        "abc",
				Event:          1,
				CompetenceDate: timestamppb.Now(),
			},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "failed to create transaction: idempotency key violation",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			request := tt.request

			if tt.seedRepo != nil {
				tx := tt.seedRepo(t)
				request.Entries[0].Id = tx.Entries[0].ID.String()

				defer tests.TruncateTables(context.Background(), testenv.DB, "entry", "account_version")
			}

			response, err := testenv.RPCClient.CreateTransaction(context.Background(), request)
			assert.Nil(t, response)

			status, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, status.Code())
			assert.Equal(t, tt.expectedMsg, status.Message())
		})
	}
}

package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests"
	"github.com/stone-co/the-amazing-ledger/app/tests/testdata"
	"github.com/stone-co/the-amazing-ledger/app/tests/testenv"
	"github.com/stone-co/the-amazing-ledger/app/tests/testseed"
	"github.com/stone-co/the-amazing-ledger/app/tests/testutils"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func TestE2E_RPC_ListAccountEntriesSuccess(t *testing.T) {
	e1 := testutils.CreateEntry(t, vos.DebitOperation, "liability.clients.available.acc1", vos.NextAccountVersion, 100)
	e2 := testutils.CreateEntry(t, vos.CreditOperation, "liability.clients.available.acc2", vos.NextAccountVersion, 100)
	e3 := testutils.CreateEntry(t, vos.DebitOperation, "liability.clients.available.acc1", vos.NextAccountVersion, 100)
	e4 := testutils.CreateEntry(t, vos.CreditOperation, "liability.clients.available.acc2", vos.NextAccountVersion, 100)

	testCases := []struct {
		name               string
		seedRepo           func(t *testing.T)
		requestSetup       func(t *testing.T) *proto.ListAccountEntriesRequest
		expectedNumEntries int
	}{
		{
			name: "should return a list of account entries",
			seedRepo: func(t *testing.T) {
				testseed.CreateTransaction(t, e1, e2)
			},
			requestSetup: func(t *testing.T) *proto.ListAccountEntriesRequest {
				return &proto.ListAccountEntriesRequest{
					AccountPath: e1.Account.Value(),
					StartDate:   timestamppb.New(time.Now().Add(-1 * time.Minute)),
					EndDate:     timestamppb.New(time.Now().Add(1 * time.Minute)),
				}
			},
			expectedNumEntries: 1,
		},
		{
			name: "should return an empty list of account entries if account path has no entries",
			seedRepo: func(t *testing.T) {
				testseed.CreateTransaction(t, e1, e2)
			},
			requestSetup: func(t *testing.T) *proto.ListAccountEntriesRequest {
				return &proto.ListAccountEntriesRequest{
					AccountPath: testdata.GenerateAccountPath(),
					StartDate:   timestamppb.New(time.Now().Add(-1 * time.Minute)),
					EndDate:     timestamppb.New(time.Now().Add(1 * time.Minute)),
				}
			},
			expectedNumEntries: 0,
		},
		{
			name: "should return an empty list of account entries if startDate/endDate is in the future",
			seedRepo: func(t *testing.T) {
				testseed.CreateTransaction(t, e1, e2)
			},
			requestSetup: func(t *testing.T) *proto.ListAccountEntriesRequest {
				return &proto.ListAccountEntriesRequest{
					AccountPath: e1.Account.Value(),
					StartDate:   timestamppb.New(time.Now().Add(1 * time.Minute)),
					EndDate:     timestamppb.New(time.Now().Add(2 * time.Minute)),
				}
			},
			expectedNumEntries: 0,
		},
		{
			name: "should return first page",
			seedRepo: func(t *testing.T) {
				testseed.CreateTransaction(t, e1, e2)
				testseed.CreateTransaction(t, e3, e4)
			},
			requestSetup: func(t *testing.T) *proto.ListAccountEntriesRequest {
				return &proto.ListAccountEntriesRequest{
					AccountPath: e1.Account.Value(),
					StartDate:   timestamppb.New(time.Now().Add(-1 * time.Minute)),
					EndDate:     timestamppb.New(time.Now().Add(1 * time.Minute)),
					Page: &proto.RequestPagination{
						PageSize: 1,
					},
				}
			},
			expectedNumEntries: 1,
		},
		{
			name: "should return second page",
			seedRepo: func(t *testing.T) {
				testseed.CreateTransaction(t, e1, e2)
				testseed.CreateTransaction(t, e3, e4)
			},
			requestSetup: func(t *testing.T) *proto.ListAccountEntriesRequest {
				request := &proto.ListAccountEntriesRequest{
					AccountPath: e1.Account.Value(),
					StartDate:   timestamppb.New(time.Now().Add(-1 * time.Minute)),
					EndDate:     timestamppb.New(time.Now().Add(1 * time.Minute)),
					Page: &proto.RequestPagination{
						PageSize: 1,
					},
				}

				accountEntries, err := testenv.RPCClient.ListAccountEntries(context.Background(), request)
				assert.NoError(t, err)

				return &proto.ListAccountEntriesRequest{
					AccountPath: e1.Account.Value(),
					StartDate:   timestamppb.New(time.Now().Add(-1 * time.Minute)),
					EndDate:     timestamppb.New(time.Now().Add(1 * time.Minute)),
					Page: &proto.RequestPagination{
						PageSize:  1,
						PageToken: accountEntries.NextPageToken,
					},
				}
			},
			expectedNumEntries: 1,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.seedRepo(t)

			defer tests.TruncateTables(context.Background(), testenv.DB, "entry", "account_version")

			request := tt.requestSetup(t)

			accountEntries, err := testenv.RPCClient.ListAccountEntries(context.Background(), request)
			assert.NoError(t, err)

			assert.Len(t, accountEntries.Entries, tt.expectedNumEntries)
		})
	}
}

func TestE2E_RPC_ListAccountEntriesFailure(t *testing.T) {
	testCases := []struct {
		name         string
		request      *proto.ListAccountEntriesRequest
		expectedCode codes.Code
		expectedMsg  string
	}{
		{
			name: "should return an error if account is invalid",
			request: &proto.ListAccountEntriesRequest{
				AccountPath: "liability..asset",
				StartDate:   timestamppb.Now(),
				EndDate:     timestamppb.Now(),
			},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "account component cannot be empty and must be less than 256 characters",
		},
		{
			name: "should return an error if start date is empty",
			request: &proto.ListAccountEntriesRequest{
				AccountPath: testdata.GenerateAccountPath(),
				EndDate:     timestamppb.Now(),
			},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "start_date must have a value",
		},
		{
			name: "should return an error if end date is empty",
			request: &proto.ListAccountEntriesRequest{
				AccountPath: testdata.GenerateAccountPath(),
				StartDate:   timestamppb.Now(),
			},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "end_date must have a value",
		},
		{
			name: "should return an error if page size <= 0",
			request: &proto.ListAccountEntriesRequest{
				AccountPath: testdata.GenerateAccountPath(),
				StartDate:   timestamppb.Now(),
				EndDate:     timestamppb.Now(),
				Page: &proto.RequestPagination{
					PageSize: 0,
				},
			},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "invalid page size",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			accountEntries, err := testenv.RPCClient.ListAccountEntries(context.Background(), tt.request)
			assert.Nil(t, accountEntries)

			status, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, status.Code())
			assert.Equal(t, tt.expectedMsg, status.Message())
		})
	}
}

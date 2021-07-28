package rpc

import (
	"context"
	"testing"
	"time"

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

func TestE2E_RPC_GetSyntheticReportSuccess(t *testing.T) {
	type wants struct {
		totalCredit int64
		totalDebit  int64
		numPaths    int
	}

	testCases := []struct {
		name         string
		repoSeed     func(t *testing.T) entities.Transaction
		requestSetup func(tx entities.Transaction) *proto.GetSyntheticReportRequest
		wants        wants
	}{
		{
			name: "should return an empty report if there are no accounts",
			repoSeed: func(t *testing.T) entities.Transaction {
				return entities.Transaction{}
			},
			requestSetup: func(tx entities.Transaction) *proto.GetSyntheticReportRequest {
				return &proto.GetSyntheticReportRequest{
					Account:   testdata.GenerateAccountPath(),
					StartDate: timestamppb.New(time.Now().Add(-1 * time.Minute)),
					EndDate:   timestamppb.New(time.Now().Add(1 * time.Minute)),
				}
			},
			wants: wants{
				totalCredit: 0,
				totalDebit:  0,
				numPaths:    0,
			},
		},
		{
			name: "should return synthetic report successfully - one path",
			repoSeed: func(t *testing.T) entities.Transaction {
				e1 := testutils.CreateEntry(t, vos.DebitOperation, testdata.GenerateAccountPath(), vos.NextAccountVersion, 100)
				e2 := testutils.CreateEntry(t, vos.CreditOperation, testdata.GenerateAccountPath(), vos.NextAccountVersion, 100)

				tx := testseed.CreateTransaction(t, e1, e2)

				return tx
			},
			requestSetup: func(tx entities.Transaction) *proto.GetSyntheticReportRequest {
				query := tx.Entries[0].Account.Value()

				if tx.Entries[0].Operation == vos.CreditOperation {
					query = tx.Entries[1].Account.Value()
				}

				return &proto.GetSyntheticReportRequest{
					Account:   query,
					StartDate: timestamppb.New(time.Now().Add(-1 * time.Minute)),
					EndDate:   timestamppb.New(time.Now().Add(1 * time.Minute)),
				}
			},
			wants: wants{
				totalCredit: 0,
				totalDebit:  100,
				numPaths:    1,
			},
		},
		{
			name: "should return synthetic report successfully - two paths",
			repoSeed: func(t *testing.T) entities.Transaction {
				e1 := testutils.CreateEntry(t, vos.DebitOperation, testdata.GenerateAccountPath(), vos.NextAccountVersion, 100)
				e2 := testutils.CreateEntry(t, vos.CreditOperation, testdata.GenerateAccountPath(), vos.NextAccountVersion, 100)

				tx := testseed.CreateTransaction(t, e1, e2)

				return tx
			},
			requestSetup: func(tx entities.Transaction) *proto.GetSyntheticReportRequest {
				return &proto.GetSyntheticReportRequest{
					Account:   "liability.clients.available.*",
					StartDate: timestamppb.New(time.Now().Add(-1 * time.Minute)),
					EndDate:   timestamppb.New(time.Now().Add(1 * time.Minute)),
				}
			},
			wants: wants{
				totalCredit: 100,
				totalDebit:  100,
				numPaths:    2,
			},
		},
		{
			name: "should return an empty report when there are no accounts for the given date range",
			repoSeed: func(t *testing.T) entities.Transaction {
				e1 := testutils.CreateEntry(t, vos.DebitOperation, testdata.GenerateAccountPath(), vos.NextAccountVersion, 100)
				e2 := testutils.CreateEntry(t, vos.CreditOperation, testdata.GenerateAccountPath(), vos.NextAccountVersion, 100)

				tx := testseed.CreateTransaction(t, e1, e2)

				return tx
			},
			requestSetup: func(tx entities.Transaction) *proto.GetSyntheticReportRequest {
				return &proto.GetSyntheticReportRequest{
					Account:   tx.Entries[0].Account.Value(),
					StartDate: timestamppb.New(time.Now().Add(1 * time.Minute)),
					EndDate:   timestamppb.New(time.Now().Add(2 * time.Minute)),
				}
			},
			wants: wants{
				totalCredit: 0,
				totalDebit:  0,
				numPaths:    0,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tx := tt.repoSeed(t)
			request := tt.requestSetup(tx)

			defer tests.TruncateTables(context.Background(), testenv.DB, "entry", "account_version")

			report, err := testenv.RPCClient.GetSyntheticReport(context.Background(), request)
			assert.NoError(t, err)

			assert.Equal(t, tt.wants.totalCredit, report.TotalCredit)
			assert.Equal(t, tt.wants.totalDebit, report.TotalDebit)
			assert.Len(t, report.Results, tt.wants.numPaths)
		})
	}
}

func TestE2E_RPC_GetSyntheticReportFailure(t *testing.T) {
	testCases := []struct {
		name         string
		request      *proto.GetSyntheticReportRequest
		expectedCode codes.Code
		expectedMsg  string
	}{
		{
			name: "should return an error if account is invalid",
			request: &proto.GetSyntheticReportRequest{
				Account:   "liability..",
				StartDate: timestamppb.New(time.Now()),
				EndDate:   timestamppb.New(time.Now()),
			},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "account component cannot be empty and must be less than 256 characters",
		},
		{
			name: "should return an error if start date is not provided",
			request: &proto.GetSyntheticReportRequest{
				Account: testdata.GenerateAccountPath(),
				EndDate: timestamppb.New(time.Now()),
			},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "start_date must have a value",
		},
		{
			name: "should return an error if end date is not provided",
			request: &proto.GetSyntheticReportRequest{
				Account:   testdata.GenerateAccountPath(),
				StartDate: timestamppb.New(time.Now()),
			},
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "end_date must have a value",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			report, err := testenv.RPCClient.GetSyntheticReport(context.Background(), tt.request)
			assert.Nil(t, report)

			status, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, status.Code())
			assert.Equal(t, tt.expectedMsg, status.Message())
		})
	}
}

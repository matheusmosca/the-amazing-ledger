package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests/mocks"
	"github.com/stone-co/the-amazing-ledger/app/tests/testdata"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func TestAPI_GetSyntheticReport(t *testing.T) {
	t.Run("should get synthetic report successfully", func(t *testing.T) {
		mockedUsecase := &mocks.UseCaseMock{
			GetSyntheticReportFunc: func(ctx context.Context, account vos.Account, level int, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {
				return &vos.SyntheticReport{}, nil
			},
		}
		api := NewAPI(logrus.New(), mockedUsecase)

		request := &proto.GetSyntheticReportRequest{
			Account:   testdata.GenerateAccount(),
			StartDate: timestamppb.Now(),
			EndDate:   timestamppb.Now(),
			Filters:   &proto.GetSyntheticReportFilters{Level: 4},
		}

		syntheticReport, err := api.GetSyntheticReport(context.Background(), request)
		assert.NoError(t, err)
		assert.NotNil(t, syntheticReport)
	})

	t.Run("should return an error if account query is invalid", func(t *testing.T) {
		mockedUsecase := &mocks.UseCaseMock{
			GetSyntheticReportFunc: func(ctx context.Context, account vos.Account, level int, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {
				return nil, app.ErrInvalidAccountComponentSize
			},
		}
		api := NewAPI(logrus.New(), mockedUsecase)

		request := &proto.GetSyntheticReportRequest{
			Account:   testdata.GenerateAccount(),
			StartDate: timestamppb.Now(),
			EndDate:   timestamppb.Now(),
			Filters:   &proto.GetSyntheticReportFilters{Level: 4},
		}

		_, err := api.GetSyntheticReport(context.Background(), request)
		respStatus, ok := status.FromError(err)

		assert.True(t, ok)
		assert.Equal(t, codes.Internal, respStatus.Code())
		assert.Equal(t, app.ErrInvalidAccountComponentSize.Error(), respStatus.Message())
	})

	t.Run("should not get synthetic report successfully, missing dates", func(t *testing.T) {
		mockedUsecase := &mocks.UseCaseMock{
			GetSyntheticReportFunc: func(ctx context.Context, account vos.Account, level int, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {
				return &vos.SyntheticReport{}, nil
			},
		}
		api := NewAPI(logrus.New(), mockedUsecase)

		request := &proto.GetSyntheticReportRequest{
			Account: testdata.GenerateAccount(),
			Filters: &proto.GetSyntheticReportFilters{},
		}

		syntheticReport, err := api.GetSyntheticReport(context.Background(), request)
		respStatus, ok := status.FromError(err)

		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, respStatus.Code())
		assert.Equal(t, "start_date must have a value", respStatus.Message())
		assert.Nil(t, syntheticReport)
	})

	t.Run("should get synthetic report successfully, zeroed level", func(t *testing.T) {
		mockedUsecase := &mocks.UseCaseMock{
			GetSyntheticReportFunc: func(ctx context.Context, account vos.Account, level int, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {
				return &vos.SyntheticReport{}, nil
			},
		}
		api := NewAPI(logrus.New(), mockedUsecase)

		request := &proto.GetSyntheticReportRequest{
			Account:   testdata.GenerateAccount(),
			StartDate: timestamppb.Now(),
			EndDate:   timestamppb.Now(),
			Filters:   &proto.GetSyntheticReportFilters{},
		}

		syntheticReport, err := api.GetSyntheticReport(context.Background(), request)
		assert.NoError(t, err)
		assert.NotNil(t, syntheticReport)
	})

	t.Run("should get synthetic report successfully, nil Filter", func(t *testing.T) {
		mockedUsecase := &mocks.UseCaseMock{
			GetSyntheticReportFunc: func(ctx context.Context, account vos.Account, level int, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {
				return &vos.SyntheticReport{}, nil
			},
		}
		api := NewAPI(logrus.New(), mockedUsecase)

		request := &proto.GetSyntheticReportRequest{
			Account:   testdata.GenerateAccount(),
			StartDate: timestamppb.Now(),
			EndDate:   timestamppb.Now(),
			Filters:   nil,
		}

		syntheticReport, err := api.GetSyntheticReport(context.Background(), request)
		assert.NoError(t, err)
		assert.NotNil(t, syntheticReport)
	})

}

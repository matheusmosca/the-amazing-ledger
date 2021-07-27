package rpc

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func (a *API) GetSyntheticReport(ctx context.Context, request *proto.GetSyntheticReportRequest) (*proto.GetSyntheticReportResponse, error) {
	log := a.log.WithFields(logrus.Fields{
		"handler": "GetSyntheticReport",
	})

	query, err := vos.NewAccount(request.Account)
	if err != nil {
		log.WithError(err).Error("Invalid account query")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	level := int(request.Filters.Level) // that's ok to convert int32 to int, since int can be int32 or int64 depending on the used system

	if request.StartDate == nil {
		return nil, status.Error(codes.InvalidArgument, "start_date must have a value")
	} else if !request.StartDate.IsValid() {
		return nil, status.Error(codes.InvalidArgument, "start_date must be valid")
	}

	if request.EndDate == nil {
		return nil, status.Error(codes.InvalidArgument, "end_date must have a value")
	} else if !request.EndDate.IsValid() {
		return nil, status.Error(codes.InvalidArgument, "end_date must be valid")
	}

	syntheticReport, err := a.UseCase.GetSyntheticReport(ctx, query, level, request.StartDate.AsTime(), request.EndDate.AsTime())
	if err != nil {
		log.WithError(err).Error("can't get synthetic report")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.GetSyntheticReportResponse{
		TotalCredit: syntheticReport.TotalCredit,
		TotalDebit:  syntheticReport.TotalDebit,
		Results:     toProto(syntheticReport.Results),
	}, nil
}

func toProto(paths []vos.AccountResult) []*proto.AccountResult {
	protoPaths := []*proto.AccountResult{}

	for _, element := range paths {
		protoPaths = append(protoPaths, &proto.AccountResult{
			Account: element.Account.Value(),
			Credit:  element.Credit,
			Debit:   element.Debit,
		})
	}

	return protoPaths
}

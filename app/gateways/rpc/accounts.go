package rpc

import (
	"context"
	"errors"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func (a *API) GetAccountBalance(ctx context.Context, request *proto.GetAccountBalanceRequest) (*proto.GetAccountBalanceResponse, error) {
	defer newrelic.FromContext(ctx).StartSegment("GetAccountBalance").End()

	log := a.log.WithFields(logrus.Fields{
		"handler": "GetAccountBalance",
	})

	accountName, err := vos.NewAccount(request.Account)
	if err != nil {
		log.WithError(err).Error("can't create account name")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	accountBalance, err := a.UseCase.GetAccountBalance(ctx, accountName)
	if err != nil {
		if errors.Is(err, app.ErrAccountNotFound) {
			log.WithError(err).Error("account name does not exist")
			return nil, status.Error(codes.NotFound, err.Error())
		}

		log.WithError(err).Error("can't get account")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.GetAccountBalanceResponse{
		Account:        accountBalance.Account.Value(),
		CurrentVersion: accountBalance.CurrentVersion.AsInt64(),
		TotalCredit:    int64(accountBalance.TotalCredit),
		TotalDebit:     int64(accountBalance.TotalDebit),
		Balance:        int64(accountBalance.Balance),
	}, nil
}

package rpc

import (
	"context"
	"errors"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func (a *API) GetAccountBalance(ctx context.Context, request *proto.GetAccountBalanceRequest) (*proto.GetAccountBalanceResponse, error) {
	defer newrelic.FromContext(ctx).StartSegment("GetAccountBalance").End()

	logger := zerolog.Ctx(ctx)
	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("handler", "GetAccountBalance")
	})

	accountName, err := vos.NewAccount(request.Account)
	if err != nil {
		logger.Error().Err(err).Msg("can't create account name")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	accountBalance, err := a.UseCase.GetAccountBalance(ctx, accountName)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get account balance")
		if errors.Is(err, app.ErrAccountNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &proto.GetAccountBalanceResponse{
		Account:        accountBalance.Account.Value(),
		CurrentVersion: accountBalance.CurrentVersion.AsInt64(),
		Balance:        int64(accountBalance.Balance),
	}, nil
}
